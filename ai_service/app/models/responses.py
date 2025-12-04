"""Pydantic response models."""

from pydantic import BaseModel, Field


class Room(BaseModel):
    """Room model."""

    name: str = Field(..., description="Room name")
    dimensions: str = Field(..., description="Room dimensions")
    area: float = Field(..., description="Room area in square feet")
    room_type: str | None = Field(None, description="Type of room")


class Opening(BaseModel):
    """Opening (door/window) model."""

    opening_type: str = Field(..., description="Type of opening (door/window)")
    count: int = Field(..., description="Number of openings")
    size: str = Field(..., description="Size specification")
    details: str | None = Field(None, description="Additional details")


class Fixture(BaseModel):
    """Fixture model."""

    fixture_type: str = Field(..., description="Type of fixture")
    category: str = Field(..., description="Category (plumbing/electrical/HVAC)")
    count: int = Field(..., description="Number of fixtures")
    details: str | None = Field(None, description="Additional details")


class Measurement(BaseModel):
    """Measurement model."""

    measurement_type: str = Field(..., description="Type of measurement")
    value: float = Field(..., description="Measurement value")
    unit: str = Field(..., description="Unit of measurement")
    location: str | None = Field(None, description="Location of measurement")


class Material(BaseModel):
    """Material model."""

    material_name: str = Field(..., description="Name of material")
    quantity: float = Field(..., description="Quantity needed")
    unit: str = Field(..., description="Unit of measurement")
    specifications: str | None = Field(None, description="Material specifications")


class AnalyzeBlueprintResponse(BaseModel):
    """Response model for blueprint analysis."""

    blueprint_id: str = Field(..., description="Blueprint identifier")
    status: str = Field(..., description="Processing status")
    rooms: list[Room] = Field(default_factory=list, description="List of identified rooms")
    openings: list[Opening] = Field(default_factory=list, description="List of doors and windows")
    fixtures: list[Fixture] = Field(default_factory=list, description="List of fixtures")
    measurements: list[Measurement] = Field(
        default_factory=list, description="List of measurements"
    )
    materials: list[Material] = Field(default_factory=list, description="List of materials")
    raw_ocr_text: str | None = Field(None, description="Raw OCR extracted text")
    confidence_score: float = Field(..., description="Overall confidence score", ge=0, le=1)
    processing_time_ms: int = Field(..., description="Processing time in milliseconds")


class LineItem(BaseModel):
    """Line item for bid."""

    description: str = Field(..., description="Item description")
    quantity: float = Field(..., description="Quantity")
    unit: str = Field(..., description="Unit of measurement")
    unit_cost: float = Field(..., description="Cost per unit")
    total: float = Field(..., description="Total cost for line item")


class GenerateBidResponse(BaseModel):
    """Response model for bid generation."""

    bid_id: str = Field(..., description="Unique bid identifier")
    project_id: str = Field(..., description="Project identifier")
    status: str = Field(..., description="Bid status")
    scope_of_work: str = Field(..., description="Detailed scope of work")
    line_items: list[LineItem] = Field(default_factory=list, description="Itemized line items")
    labor_cost: float = Field(..., description="Total labor cost")
    material_cost: float = Field(..., description="Total material cost")
    subtotal: float = Field(..., description="Subtotal before markup")
    markup_amount: float = Field(..., description="Markup amount")
    total_price: float = Field(..., description="Total bid price")
    exclusions: list[str] = Field(default_factory=list, description="Items NOT included")
    inclusions: list[str] = Field(default_factory=list, description="Items specifically included")
    schedule: dict = Field(default_factory=dict, description="Project schedule with milestones")
    payment_terms: str = Field(..., description="Payment terms")
    warranty_terms: str = Field(..., description="Warranty terms")
    closing_statement: str = Field(..., description="Professional closing statement")
