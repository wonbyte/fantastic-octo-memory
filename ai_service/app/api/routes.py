"""API route handlers."""

import time

from fastapi import APIRouter, HTTPException, status

from app.core.logging import get_logger
from app.models.requests import AnalyzeBlueprintRequest, GenerateBidRequest
from app.models.responses import AnalyzeBlueprintResponse, GenerateBidResponse
from app.services.bid_service import BidService
from app.services.ocr_service import OCRService
from app.services.s3_service import S3Service
from app.services.vision_service import VisionService

logger = get_logger(__name__)

router = APIRouter()


@router.get("/health")
async def health() -> dict[str, str]:
    """
    Health check endpoint.

    Returns:
        Health status
    """
    return {"status": "ok", "version": "1.0.0"}


@router.post(
    "/analyze-blueprint",
    response_model=AnalyzeBlueprintResponse,
    status_code=status.HTTP_200_OK,
)
async def analyze_blueprint(request: AnalyzeBlueprintRequest) -> AnalyzeBlueprintResponse:
    """
    Analyze blueprint using OCR and vision models.

    Args:
        request: Blueprint analysis request

    Returns:
        Blueprint analysis result with rooms, openings, fixtures, etc.

    Raises:
        HTTPException: If analysis fails
    """
    start_time = time.time()

    try:
        logger.info(
            "analyze_blueprint_request",
            blueprint_id=request.blueprint_id,
            s3_key=request.s3_key,
        )

        # Initialize services
        s3_service = S3Service()
        ocr_service = OCRService()
        vision_service = VisionService()

        # Download blueprint from S3
        logger.info("downloading_blueprint", s3_key=request.s3_key)
        file_bytes = await s3_service.download_file(request.s3_key)

        # Determine file type from s3_key
        file_type = "pdf" if request.s3_key.lower().endswith(".pdf") else "image"

        # Extract text using OCR
        logger.info("extracting_text_ocr")
        ocr_result = await ocr_service.extract_text(file_bytes, file_type)

        # Analyze blueprint with vision model
        logger.info("analyzing_with_vision_model")
        # For PDFs, convert first page to image using OCR service
        if file_type == "pdf":
            # Convert PDF to image for vision analysis
            from io import BytesIO

            from pdf2image import convert_from_bytes
            images = convert_from_bytes(file_bytes, dpi=200, first_page=1, last_page=1)
            if images:
                img_byte_arr = BytesIO()
                images[0].save(img_byte_arr, format="PNG")
                vision_bytes = img_byte_arr.getvalue()
            else:
                raise HTTPException(
                    status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
                    detail="Failed to convert PDF to image for vision analysis",
                )
        else:
            vision_bytes = file_bytes

        analysis = await vision_service.analyze_blueprint(
            vision_bytes,
            ocr_result.raw_text,
            request.options,
        )

        # Calculate processing time
        processing_time_ms = int((time.time() - start_time) * 1000)

        # Build response
        response = AnalyzeBlueprintResponse(
            blueprint_id=request.blueprint_id,
            status="completed",
            rooms=analysis.rooms,
            openings=analysis.openings,
            fixtures=analysis.fixtures,
            measurements=analysis.measurements,
            materials=analysis.materials,
            raw_ocr_text=ocr_result.raw_text,
            confidence_score=analysis.confidence_score,
            processing_time_ms=processing_time_ms,
        )

        logger.info(
            "blueprint_analysis_complete",
            blueprint_id=request.blueprint_id,
            processing_time_ms=processing_time_ms,
            confidence=analysis.confidence_score,
        )

        return response

    except Exception as e:
        logger.error("blueprint_analysis_error", blueprint_id=request.blueprint_id, error=str(e))
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Blueprint analysis failed: {str(e)}",
        ) from e


@router.post(
    "/generate-bid",
    response_model=GenerateBidResponse,
    status_code=status.HTTP_200_OK,
)
async def generate_bid(request: GenerateBidRequest) -> GenerateBidResponse:
    """
    Generate professional bid package from takeoff data.

    Args:
        request: Bid generation request

    Returns:
        Complete bid package with line items, costs, terms, etc.

    Raises:
        HTTPException: If bid generation fails
    """
    try:
        logger.info(
            "generate_bid_request",
            project_id=request.project_id,
            blueprint_id=request.blueprint_id,
        )

        # Initialize bid service
        bid_service = BidService()

        # Prepare project info
        project_info = {
            "project_id": request.project_id,
            "blueprint_id": request.blueprint_id,
        }

        # Generate bid
        logger.info("generating_bid_package")
        bid_package = await bid_service.generate_bid(
            takeoff_data=request.takeoff_data,
            pricing_rules=request.pricing_rules,
            company_info=request.company_info,
            project_info=project_info,
            markup_percentage=request.markup_percentage,
        )

        # Build response
        response = GenerateBidResponse(
            bid_id=bid_package.bid_id,
            project_id=request.project_id,
            status="completed",
            scope_of_work=bid_package.scope_of_work,
            line_items=bid_package.line_items,
            labor_cost=bid_package.labor_cost,
            material_cost=bid_package.material_cost,
            subtotal=bid_package.subtotal,
            markup_amount=bid_package.markup_amount,
            total_price=bid_package.total_price,
            exclusions=bid_package.exclusions,
            inclusions=bid_package.inclusions,
            schedule=bid_package.schedule,
            payment_terms=bid_package.payment_terms,
            warranty_terms=bid_package.warranty_terms,
            closing_statement=bid_package.closing_statement,
        )

        logger.info(
            "bid_generation_complete",
            bid_id=bid_package.bid_id,
            total_price=bid_package.total_price,
        )

        return response

    except Exception as e:
        logger.error("bid_generation_error", project_id=request.project_id, error=str(e))
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Bid generation failed: {str(e)}",
        ) from e
