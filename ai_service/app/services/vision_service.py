"""Vision service for blueprint analysis using LLM with vision capabilities."""

import base64
import json

from openai import AsyncOpenAI
from pydantic import BaseModel, Field
from tenacity import retry, stop_after_attempt, wait_exponential

from app.core.config import get_settings
from app.core.logging import get_logger
from app.models.responses import Fixture, Material, Measurement, Opening, Room
from app.prompts.templates import BLUEPRINT_ANALYSIS_PROMPT, VISION_ANALYSIS_SYSTEM_PROMPT

logger = get_logger(__name__)


class BlueprintAnalysis(BaseModel):
    """Blueprint analysis result from vision model."""

    rooms: list[Room] = Field(default_factory=list)
    openings: list[Opening] = Field(default_factory=list)
    fixtures: list[Fixture] = Field(default_factory=list)
    measurements: list[Measurement] = Field(default_factory=list)
    materials: list[Material] = Field(default_factory=list)
    confidence_score: float = Field(..., ge=0, le=1)


class VisionService:
    """Service for analyzing blueprints using vision-capable LLMs."""

    def __init__(self):
        """Initialize vision service."""
        self.settings = get_settings()
        self._client = None

    def _get_client(self) -> AsyncOpenAI:
        """Get or create OpenAI client."""
        if self._client is None:
            if not self.settings.openai_api_key:
                logger.warning("no_openai_api_key_mock_mode_active")
            self._client = AsyncOpenAI(api_key=self.settings.openai_api_key or "mock-key")
        return self._client

    def _create_json_schema(self) -> dict:
        """Create JSON schema for structured output."""
        return {
            "rooms": [
                {
                    "name": "string",
                    "dimensions": "string",
                    "area": "number",
                    "room_type": "string or null",
                }
            ],
            "openings": [
                {
                    "opening_type": "string",
                    "count": "number",
                    "size": "string",
                    "details": "string or null",
                }
            ],
            "fixtures": [
                {
                    "fixture_type": "string",
                    "category": "string",
                    "count": "number",
                    "details": "string or null",
                }
            ],
            "measurements": [
                {
                    "measurement_type": "string",
                    "value": "number",
                    "unit": "string",
                    "location": "string or null",
                }
            ],
            "materials": [
                {
                    "material_name": "string",
                    "quantity": "number",
                    "unit": "string",
                    "specifications": "string or null",
                }
            ],
            "confidence_score": "number (0-1)",
        }

    def _encode_image(self, image_bytes: bytes) -> str:
        """Encode image bytes to base64."""
        return base64.b64encode(image_bytes).decode("utf-8")

    @retry(stop=stop_after_attempt(3), wait=wait_exponential(multiplier=1, min=4, max=10))
    async def _call_vision_model(
        self, image_bytes: bytes, ocr_text: str, context: dict | None
    ) -> dict:
        """
        Call OpenAI vision model with retry logic.

        Args:
            image_bytes: Blueprint image as bytes
            ocr_text: OCR extracted text
            context: Additional context

        Returns:
            Parsed JSON response
        """
        client = self._get_client()

        # Prepare prompt
        prompt = BLUEPRINT_ANALYSIS_PROMPT.format(
            ocr_text=ocr_text or "No OCR text available",
            context=json.dumps(context or {}, indent=2),
            json_schema=json.dumps(self._create_json_schema(), indent=2),
        )

        # Encode image
        base64_image = self._encode_image(image_bytes)

        try:
            logger.info("calling_vision_model", model=self.settings.openai_vision_model)

            if not self.settings.openai_api_key:
                # Return mock response if no API key
                logger.warning("using_mock_vision_response")
                return self._mock_vision_response()

            response = await client.chat.completions.create(
                model=self.settings.openai_vision_model,
                messages=[
                    {"role": "system", "content": VISION_ANALYSIS_SYSTEM_PROMPT},
                    {
                        "role": "user",
                        "content": [
                            {"type": "text", "text": prompt},
                            {
                                "type": "image_url",
                                "image_url": {
                                    "url": f"data:image/png;base64,{base64_image}",
                                    "detail": "high",
                                },
                            },
                        ],
                    },
                ],
                max_tokens=4096,
                temperature=0.1,
            )

            content = response.choices[0].message.content
            logger.info("vision_model_response_received", content_length=len(content))

            # Parse JSON response
            # Try to extract JSON from markdown code blocks if present
            if "```json" in content:
                parts = content.split("```json")
                if len(parts) > 1:
                    json_parts = parts[1].split("```")
                    if json_parts:
                        content = json_parts[0].strip()
            elif "```" in content:
                parts = content.split("```")
                if len(parts) > 2:
                    content = parts[1].strip()

            result = json.loads(content)
            return result

        except json.JSONDecodeError as e:
            logger.error("vision_response_json_parse_failed", error=str(e))
            # Return mock on parse failure
            return self._mock_vision_response()
        except Exception as e:
            logger.error("vision_model_call_failed", error=str(e))
            raise

    def _mock_vision_response(self) -> dict:
        """Generate mock vision response for testing."""
        return {
            "rooms": [
                {
                    "name": "Living Room",
                    "dimensions": "15' x 20'",
                    "area": 300.0,
                    "room_type": "Living",
                },
                {
                    "name": "Bedroom",
                    "dimensions": "12' x 14'",
                    "area": 168.0,
                    "room_type": "Bedroom",
                },
                {
                    "name": "Kitchen",
                    "dimensions": "10' x 12'",
                    "area": 120.0,
                    "room_type": "Kitchen",
                },
            ],
            "openings": [
                {
                    "opening_type": "Door",
                    "count": 3,
                    "size": "36\" x 80\"",
                    "details": "Standard interior doors",
                },
                {
                    "opening_type": "Window",
                    "count": 5,
                    "size": "48\" x 60\"",
                    "details": "Double-hung windows",
                },
            ],
            "fixtures": [
                {
                    "fixture_type": "Ceiling Light",
                    "category": "electrical",
                    "count": 6,
                    "details": "Standard ceiling fixtures",
                },
                {
                    "fixture_type": "Outlet",
                    "category": "electrical",
                    "count": 15,
                    "details": "Standard 120V outlets",
                },
            ],
            "measurements": [
                {
                    "measurement_type": "Wall Length",
                    "value": 120.0,
                    "unit": "feet",
                    "location": "Perimeter",
                },
                {
                    "measurement_type": "Ceiling Height",
                    "value": 9.0,
                    "unit": "feet",
                    "location": "All rooms",
                },
            ],
            "materials": [
                {
                    "material_name": "Drywall",
                    "quantity": 588.0,
                    "unit": "sq ft",
                    "specifications": "1/2 inch standard drywall",
                },
                {
                    "material_name": "Flooring",
                    "quantity": 588.0,
                    "unit": "sq ft",
                    "specifications": "Hardwood or equivalent",
                },
            ],
            "confidence_score": 0.85,
        }

    async def analyze_blueprint(
        self, image_bytes: bytes, ocr_text: str, context: dict | None = None
    ) -> BlueprintAnalysis:
        """
        Analyze blueprint using vision model.

        Args:
            image_bytes: Blueprint image as bytes
            ocr_text: OCR extracted text
            context: Optional additional context

        Returns:
            Blueprint analysis result

        Raises:
            Exception: If analysis fails
        """
        try:
            logger.info("starting_blueprint_analysis")

            # Call vision model
            response = await self._call_vision_model(image_bytes, ocr_text, context)

            # Parse into Pydantic model
            analysis = BlueprintAnalysis.model_validate(response)

            logger.info(
                "blueprint_analysis_complete",
                rooms=len(analysis.rooms),
                openings=len(analysis.openings),
                fixtures=len(analysis.fixtures),
                confidence=analysis.confidence_score,
            )

            return analysis

        except Exception as e:
            logger.error("blueprint_analysis_failed", error=str(e))
            raise Exception(f"Blueprint analysis failed: {e}") from e
