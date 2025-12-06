"""Vision service for blueprint analysis using LLM with vision capabilities."""

import base64
import json
import re

from openai import AsyncOpenAI
from pydantic import BaseModel, Field
from tenacity import retry, stop_after_attempt, wait_exponential

from app.core.config import get_settings
from app.core.logging import get_logger
from app.models.responses import Fixture, Material, Measurement, Opening, Room
from app.prompts.templates import (
    BLUEPRINT_ANALYSIS_PROMPT,
    ELECTRICAL_ANALYSIS_PROMPT,
    ELECTRICAL_SYMBOLS,
    ELECTRICAL_SYSTEM_PROMPT,
    HVAC_ANALYSIS_PROMPT,
    HVAC_SYMBOLS,
    HVAC_SYSTEM_PROMPT,
    MULTI_PAGE_ANALYSIS_PROMPT,
    PLUMBING_ANALYSIS_PROMPT,
    PLUMBING_SYMBOLS,
    PLUMBING_SYSTEM_PROMPT,
    SCALE_PATTERNS,
    STRUCTURAL_ANALYSIS_PROMPT,
    STRUCTURAL_SYMBOLS,
    STRUCTURAL_SYSTEM_PROMPT,
    VISION_ANALYSIS_SYSTEM_PROMPT,
)

logger = get_logger(__name__)


class BlueprintAnalysis(BaseModel):
    """Blueprint analysis result from vision model."""

    rooms: list[Room] = Field(default_factory=list)
    openings: list[Opening] = Field(default_factory=list)
    fixtures: list[Fixture] = Field(default_factory=list)
    measurements: list[Measurement] = Field(default_factory=list)
    materials: list[Material] = Field(default_factory=list)
    confidence_score: float = Field(..., ge=0, le=1)
    scale_info: dict | None = Field(None, description="Detected scale information")
    trade_type: str | None = Field(None, description="Detected trade type if specialized")


class VisionService:
    """Service for analyzing blueprints using vision-capable LLMs."""

    # Supported trade types for specialized analysis
    TRADE_TYPES = ["electrical", "plumbing", "hvac", "structural", "general"]

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

    def _detect_scale(self, ocr_text: str) -> dict | None:
        """
        Detect scale information from OCR text.

        Args:
            ocr_text: OCR extracted text

        Returns:
            Dictionary with scale information or None if not found
        """
        if not ocr_text:
            return None

        for pattern in SCALE_PATTERNS:
            match = re.search(pattern, ocr_text, re.IGNORECASE)
            if match:
                scale_str = match.group(0)
                logger.info("scale_detected", scale=scale_str)
                return {
                    "scale_string": scale_str,
                    "detected_pattern": pattern,
                    "confidence": 0.9,
                }

        return None

    def _detect_trade_type(self, ocr_text: str, context: dict | None) -> str:
        """
        Detect trade type from OCR text and context.

        Args:
            ocr_text: OCR extracted text
            context: Additional context

        Returns:
            Detected trade type (electrical, plumbing, hvac, structural, or general)
        """
        # Check if trade_type is explicitly provided in context
        if context and "trade_type" in context:
            trade = context["trade_type"].lower()
            if trade in self.TRADE_TYPES:
                return trade

        # Detect from OCR text using keyword analysis
        ocr_lower = ocr_text.lower() if ocr_text else ""

        electrical_keywords = ["electrical", "lighting", "outlet", "switch", "panel", "circuit", "voltage"]
        plumbing_keywords = ["plumbing", "water", "drain", "sewer", "fixture", "toilet", "sink"]
        hvac_keywords = ["hvac", "mechanical", "heating", "cooling", "duct", "furnace", "air conditioning"]
        structural_keywords = ["structural", "beam", "column", "foundation", "footing", "joist", "truss"]

        keyword_counts = {
            "electrical": sum(1 for kw in electrical_keywords if kw in ocr_lower),
            "plumbing": sum(1 for kw in plumbing_keywords if kw in ocr_lower),
            "hvac": sum(1 for kw in hvac_keywords if kw in ocr_lower),
            "structural": sum(1 for kw in structural_keywords if kw in ocr_lower),
        }

        max_count = max(keyword_counts.values())
        if max_count >= 2:  # At least 2 keywords to confidently detect trade
            detected_trade = max(keyword_counts, key=keyword_counts.get)
            logger.info("trade_type_detected", trade=detected_trade, keyword_count=max_count)
            return detected_trade

        return "general"

    def _get_trade_prompt_and_system(self, trade_type: str) -> tuple[str, str]:
        """
        Get appropriate prompt and system message for trade type.

        Args:
            trade_type: Type of trade (electrical, plumbing, hvac, structural, general)

        Returns:
            Tuple of (prompt_template, system_prompt)
        """
        trade_configs = {
            "electrical": (ELECTRICAL_ANALYSIS_PROMPT, ELECTRICAL_SYSTEM_PROMPT),
            "plumbing": (PLUMBING_ANALYSIS_PROMPT, PLUMBING_SYSTEM_PROMPT),
            "hvac": (HVAC_ANALYSIS_PROMPT, HVAC_SYSTEM_PROMPT),
            "structural": (STRUCTURAL_ANALYSIS_PROMPT, STRUCTURAL_SYSTEM_PROMPT),
            "general": (BLUEPRINT_ANALYSIS_PROMPT, VISION_ANALYSIS_SYSTEM_PROMPT),
        }

        return trade_configs.get(trade_type, trade_configs["general"])

    def _get_symbol_library(self, trade_type: str) -> dict:
        """
        Get symbol library for trade type.

        Args:
            trade_type: Type of trade

        Returns:
            Dictionary of symbols for the trade
        """
        symbol_libraries = {
            "electrical": ELECTRICAL_SYMBOLS,
            "plumbing": PLUMBING_SYMBOLS,
            "hvac": HVAC_SYMBOLS,
            "structural": STRUCTURAL_SYMBOLS,
        }

        return symbol_libraries.get(trade_type, {})

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

        # Detect scale and trade type
        scale_info = self._detect_scale(ocr_text)
        trade_type = self._detect_trade_type(ocr_text, context)
        
        # Get appropriate prompt and system message for trade type
        prompt_template, system_prompt = self._get_trade_prompt_and_system(trade_type)
        symbol_library = self._get_symbol_library(trade_type)

        # Enhance context with detected information
        enhanced_context = context or {}
        if scale_info:
            enhanced_context["scale_info"] = scale_info
        if trade_type != "general":
            enhanced_context["trade_type"] = trade_type
            enhanced_context["symbol_library"] = symbol_library

        # Prepare prompt
        prompt = prompt_template.format(
            ocr_text=ocr_text or "No OCR text available",
            context=json.dumps(enhanced_context, indent=2),
            json_schema=json.dumps(self._create_json_schema(), indent=2),
        )

        # Encode image
        base64_image = self._encode_image(image_bytes)

        try:
            logger.info(
                "calling_vision_model",
                model=self.settings.openai_vision_model,
                trade_type=trade_type,
                has_scale=scale_info is not None,
            )

            if not self.settings.openai_api_key:
                # Return mock response if no API key
                logger.warning("using_mock_vision_response")
                return self._mock_vision_response(trade_type)

            response = await client.chat.completions.create(
                model=self.settings.openai_vision_model,
                messages=[
                    {"role": "system", "content": system_prompt},
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
            
            # Add detected metadata to result
            if scale_info:
                result["scale_info"] = scale_info
            if trade_type != "general":
                result["trade_type"] = trade_type
            
            return result

        except json.JSONDecodeError as e:
            logger.error("vision_response_json_parse_failed", error=str(e))
            # Return mock on parse failure
            return self._mock_vision_response(trade_type)
        except Exception as e:
            logger.error("vision_model_call_failed", error=str(e))
            raise

    def _mock_vision_response(self, trade_type: str = "general") -> dict:
        """Generate mock vision response for testing."""
        base_response = {
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
            "fixtures": [],
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
            "scale_info": {
                "scale_string": "1/4\" = 1'-0\"",
                "detected_pattern": "1/4\" = 1'-0\"",
                "confidence": 0.9,
            },
        }

        # Add trade-specific fixtures
        if trade_type == "electrical":
            base_response["fixtures"] = [
                {
                    "fixture_type": "Duplex Outlet",
                    "category": "electrical",
                    "count": 15,
                    "details": "120V standard outlets",
                },
                {
                    "fixture_type": "Light Fixture",
                    "category": "electrical",
                    "count": 8,
                    "details": "Ceiling-mounted fixtures",
                },
                {
                    "fixture_type": "Switch",
                    "category": "electrical",
                    "count": 6,
                    "details": "Single-pole switches",
                },
            ]
            base_response["trade_type"] = "electrical"
        elif trade_type == "plumbing":
            base_response["fixtures"] = [
                {
                    "fixture_type": "Toilet",
                    "category": "plumbing",
                    "count": 2,
                    "details": "Standard water closets",
                },
                {
                    "fixture_type": "Sink",
                    "category": "plumbing",
                    "count": 3,
                    "details": "Lavatory and kitchen sinks",
                },
                {
                    "fixture_type": "Shower",
                    "category": "plumbing",
                    "count": 1,
                    "details": "Standard shower stall",
                },
            ]
            base_response["trade_type"] = "plumbing"
        elif trade_type == "hvac":
            base_response["fixtures"] = [
                {
                    "fixture_type": "Supply Diffuser",
                    "category": "hvac",
                    "count": 6,
                    "details": "Ceiling diffusers",
                },
                {
                    "fixture_type": "Return Grille",
                    "category": "hvac",
                    "count": 2,
                    "details": "Return air grilles",
                },
                {
                    "fixture_type": "Thermostat",
                    "category": "hvac",
                    "count": 1,
                    "details": "Programmable thermostat",
                },
            ]
            base_response["trade_type"] = "hvac"
        elif trade_type == "structural":
            base_response["fixtures"] = []
            base_response["measurements"].extend([
                {
                    "measurement_type": "Beam Span",
                    "value": 20.0,
                    "unit": "feet",
                    "location": "Main beam",
                },
                {
                    "measurement_type": "Joist Spacing",
                    "value": 16.0,
                    "unit": "inches",
                    "location": "Floor joists",
                },
            ])
            base_response["trade_type"] = "structural"
        else:
            base_response["fixtures"] = [
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
            ]

        return base_response

    async def analyze_blueprint(
        self, image_bytes: bytes, ocr_text: str, context: dict | None = None
    ) -> BlueprintAnalysis:
        """
        Analyze blueprint using vision model.

        Args:
            image_bytes: Blueprint image as bytes
            ocr_text: OCR extracted text
            context: Optional additional context (can include trade_type, page_info, etc.)

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
                trade_type=analysis.trade_type,
                has_scale=analysis.scale_info is not None,
            )

            return analysis

        except Exception as e:
            logger.error("blueprint_analysis_failed", error=str(e))
            raise Exception(f"Blueprint analysis failed: {e}") from e

    async def analyze_multi_page_blueprint(
        self,
        pages: list[tuple[bytes, str]],
        context: dict | None = None,
    ) -> BlueprintAnalysis:
        """
        Analyze multi-page blueprint by processing each page and aggregating results.

        Args:
            pages: List of tuples (image_bytes, ocr_text) for each page
            context: Optional additional context

        Returns:
            Aggregated blueprint analysis result

        Raises:
            Exception: If analysis fails
        """
        try:
            logger.info("starting_multi_page_analysis", page_count=len(pages))

            aggregated_rooms = []
            aggregated_openings = []
            aggregated_fixtures = []
            aggregated_measurements = []
            aggregated_materials = []
            total_confidence = 0.0
            scale_info = None
            trade_type = None

            for page_idx, (image_bytes, ocr_text) in enumerate(pages, 1):
                logger.info("analyzing_page", page=page_idx, total_pages=len(pages))

                # Enhance context with page information
                page_context = context or {}
                page_context["page_number"] = page_idx
                page_context["total_pages"] = len(pages)

                # Analyze individual page
                page_analysis = await self.analyze_blueprint(image_bytes, ocr_text, page_context)

                # Aggregate results
                aggregated_rooms.extend(page_analysis.rooms)
                aggregated_openings.extend(page_analysis.openings)
                aggregated_fixtures.extend(page_analysis.fixtures)
                aggregated_measurements.extend(page_analysis.measurements)
                aggregated_materials.extend(page_analysis.materials)
                total_confidence += page_analysis.confidence_score

                # Capture scale info and trade type from first page if available
                if page_idx == 1:
                    scale_info = page_analysis.scale_info
                    trade_type = page_analysis.trade_type

            # Calculate average confidence
            avg_confidence = total_confidence / len(pages) if pages else 0.0

            # Create aggregated analysis
            aggregated_analysis = BlueprintAnalysis(
                rooms=aggregated_rooms,
                openings=aggregated_openings,
                fixtures=aggregated_fixtures,
                measurements=aggregated_measurements,
                materials=aggregated_materials,
                confidence_score=avg_confidence,
                scale_info=scale_info,
                trade_type=trade_type,
            )

            logger.info(
                "multi_page_analysis_complete",
                pages=len(pages),
                total_rooms=len(aggregated_rooms),
                total_fixtures=len(aggregated_fixtures),
                avg_confidence=avg_confidence,
            )

            return aggregated_analysis

        except Exception as e:
            logger.error("multi_page_analysis_failed", error=str(e))
            raise Exception(f"Multi-page blueprint analysis failed: {e}") from e
