"""OCR service for text extraction from blueprints."""

import io
from typing import Any

import boto3
from pdf2image import convert_from_bytes
from PIL import Image
from pydantic import BaseModel, Field
from tenacity import retry, stop_after_attempt, wait_exponential

from app.core.config import get_settings
from app.core.logging import get_logger

logger = get_logger(__name__)


class TextBlock(BaseModel):
    """Text block with position and confidence."""

    text: str = Field(..., description="Extracted text")
    confidence: float = Field(..., description="Confidence score")
    bounding_box: dict[str, float] = Field(..., description="Bounding box coordinates")
    block_type: str = Field(default="LINE", description="Type of block")


class Table(BaseModel):
    """Table structure extracted from document."""

    rows: int = Field(..., description="Number of rows")
    columns: int = Field(..., description="Number of columns")
    cells: list[list[str]] = Field(..., description="Cell contents")


class FormField(BaseModel):
    """Form field extracted from document."""

    key: str = Field(..., description="Field key/label")
    value: str = Field(..., description="Field value")
    confidence: float = Field(..., description="Confidence score")


class OCRResult(BaseModel):
    """Result of OCR processing."""

    raw_text: str = Field(..., description="Complete extracted text")
    blocks: list[TextBlock] = Field(default_factory=list, description="Text blocks")
    tables: list[Table] | None = Field(None, description="Extracted tables")
    forms: list[FormField] | None = Field(None, description="Extracted form fields")
    page_count: int = Field(..., description="Number of pages processed")


class OCRService:
    """Service for OCR text extraction using AWS Textract."""

    def __init__(self):
        """Initialize OCR service."""
        self.settings = get_settings()
        self._textract_client = None

    def _get_textract_client(self) -> Any:
        """Get or create Textract client."""
        if self._textract_client is None and self.settings.aws_access_key_id:
            self._textract_client = boto3.client(
                "textract",
                aws_access_key_id=self.settings.aws_access_key_id,
                aws_secret_access_key=self.settings.aws_secret_access_key,
                region_name=self.settings.aws_region,
            )
        return self._textract_client

    def _pdf_to_images(self, pdf_bytes: bytes) -> list[Image.Image]:
        """
        Convert PDF to list of images.

        Args:
            pdf_bytes: PDF file content

        Returns:
            List of PIL Images
        """
        try:
            logger.info("converting_pdf_to_images")
            images = convert_from_bytes(pdf_bytes, dpi=300)
            logger.info("pdf_converted", page_count=len(images))
            return images
        except Exception as e:
            logger.error("pdf_conversion_failed", error=str(e))
            raise Exception(f"Failed to convert PDF to images: {e}") from e

    def _image_to_bytes(self, image: Image.Image) -> bytes:
        """Convert PIL Image to bytes."""
        img_byte_arr = io.BytesIO()
        image.save(img_byte_arr, format="PNG")
        return img_byte_arr.getvalue()

    @retry(stop=stop_after_attempt(3), wait=wait_exponential(multiplier=1, min=4, max=10))
    async def _call_textract(self, image_bytes: bytes) -> dict:
        """
        Call AWS Textract API with retry logic.

        Args:
            image_bytes: Image content as bytes

        Returns:
            Textract response
        """
        client = self._get_textract_client()
        if client is None:
            # Fallback to mock if no credentials
            logger.warning("no_textract_credentials_using_mock")
            return self._mock_textract_response(image_bytes)

        try:
            logger.info("calling_textract_api")
            response = client.detect_document_text(Document={"Bytes": image_bytes})
            logger.info("textract_api_success", block_count=len(response.get("Blocks", [])))
            return response
        except Exception as e:
            logger.error("textract_api_failed", error=str(e))
            raise

    def _mock_textract_response(self, image_bytes: bytes) -> dict:
        """Generate mock Textract response for testing."""
        return {
            "Blocks": [
                {
                    "BlockType": "LINE",
                    "Text": "Blueprint - Floor Plan",
                    "Confidence": 95.5,
                    "Geometry": {
                        "BoundingBox": {
                            "Width": 0.2,
                            "Height": 0.05,
                            "Left": 0.1,
                            "Top": 0.1,
                        }
                    },
                },
                {
                    "BlockType": "LINE",
                    "Text": "Living Room: 15' x 20'",
                    "Confidence": 92.3,
                    "Geometry": {
                        "BoundingBox": {
                            "Width": 0.25,
                            "Height": 0.05,
                            "Left": 0.1,
                            "Top": 0.2,
                        }
                    },
                },
                {
                    "BlockType": "LINE",
                    "Text": "Bedroom: 12' x 14'",
                    "Confidence": 91.8,
                    "Geometry": {
                        "BoundingBox": {
                            "Width": 0.22,
                            "Height": 0.05,
                            "Left": 0.1,
                            "Top": 0.3,
                        }
                    },
                },
            ]
        }

    def _parse_textract_response(self, response: dict) -> OCRResult:
        """
        Parse Textract response into OCRResult.

        Args:
            response: Textract API response

        Returns:
            Normalized OCR result
        """
        blocks = []
        raw_text_lines = []

        for block in response.get("Blocks", []):
            if block["BlockType"] == "LINE":
                text = block.get("Text", "")
                raw_text_lines.append(text)

                geometry = block.get("Geometry", {})
                bbox = geometry.get("BoundingBox", {})

                blocks.append(
                    TextBlock(
                        text=text,
                        confidence=block.get("Confidence", 0.0) / 100.0,
                        bounding_box={
                            "left": bbox.get("Left", 0.0),
                            "top": bbox.get("Top", 0.0),
                            "width": bbox.get("Width", 0.0),
                            "height": bbox.get("Height", 0.0),
                        },
                        block_type=block["BlockType"],
                    )
                )

        raw_text = "\n".join(raw_text_lines)

        return OCRResult(
            raw_text=raw_text,
            blocks=blocks,
            tables=None,
            forms=None,
            page_count=1,
        )

    async def extract_text(self, file_bytes: bytes, file_type: str = "pdf") -> OCRResult:
        """
        Extract text from blueprint using OCR.

        Args:
            file_bytes: File content as bytes
            file_type: File type (pdf, png, jpg)

        Returns:
            OCR result with extracted text and metadata

        Raises:
            Exception: If OCR processing fails
        """
        try:
            logger.info("starting_ocr_extraction", file_type=file_type)

            # Convert PDF to images if needed
            if file_type.lower() == "pdf":
                images = self._pdf_to_images(file_bytes)
                # For now, process only the first page
                # In production, you'd aggregate results from all pages
                image_bytes = self._image_to_bytes(images[0])
                page_count = len(images)
            else:
                image_bytes = file_bytes
                page_count = 1

            # Call Textract
            response = await self._call_textract(image_bytes)

            # Parse response
            result = self._parse_textract_response(response)
            result.page_count = page_count

            logger.info(
                "ocr_extraction_complete",
                page_count=page_count,
                block_count=len(result.blocks),
                text_length=len(result.raw_text),
            )

            return result

        except Exception as e:
            logger.error("ocr_extraction_failed", error=str(e))
            raise Exception(f"OCR extraction failed: {e}") from e
