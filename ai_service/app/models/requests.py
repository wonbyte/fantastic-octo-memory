"""Pydantic request models."""

from pydantic import BaseModel, Field


class AnalyzeBlueprintRequest(BaseModel):
    """Request model for blueprint analysis."""

    blueprint_id: str = Field(..., description="Unique identifier for the blueprint")
    s3_key: str = Field(..., description="S3 key where blueprint is stored")
    project_name: str | None = Field(None, description="Optional project name")
    options: dict | None = Field(None, description="Optional analysis options")


class GenerateBidRequest(BaseModel):
    """Request model for bid generation."""

    project_id: str = Field(..., description="Project identifier")
    blueprint_id: str = Field(..., description="Blueprint identifier")
    takeoff_data: dict = Field(..., description="Material takeoff data from analysis")
    pricing_rules: dict | None = Field(None, description="Optional pricing rules")
    company_info: dict | None = Field(None, description="Optional company information")
    markup_percentage: float = Field(default=20.0, description="Markup percentage", ge=0, le=100)
