"""Application configuration using Pydantic settings."""

from functools import lru_cache

from pydantic import Field
from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    """Application settings."""

    model_config = SettingsConfigDict(
        env_file=".env",
        env_file_encoding="utf-8",
        case_sensitive=False,
        extra="ignore",
    )

    # Server
    port: int = Field(default=8000, description="Server port")
    env: str = Field(default="development", description="Environment")
    log_level: str = Field(default="INFO", description="Log level")

    # S3/MinIO
    s3_endpoint: str = Field(default="http://minio:9000", description="S3 endpoint")
    s3_access_key: str = Field(default="minioadmin", description="S3 access key")
    s3_secret_key: str = Field(default="minioadmin", description="S3 secret key")
    s3_bucket: str = Field(default="blueprints", description="S3 bucket name")
    s3_region: str = Field(default="us-east-1", description="S3 region")

    # OpenAI
    openai_api_key: str = Field(default="", description="OpenAI API key")
    openai_model: str = Field(default="gpt-4o", description="OpenAI text model")
    openai_vision_model: str = Field(default="gpt-4o", description="OpenAI vision model")

    # AWS (for Textract)
    aws_access_key_id: str = Field(default="", description="AWS access key ID")
    aws_secret_access_key: str = Field(default="", description="AWS secret access key")
    aws_region: str = Field(default="us-east-1", description="AWS region")

    # Google Cloud (alternative OCR)
    google_application_credentials: str = Field(
        default="", description="Path to Google credentials JSON"
    )

    # Redis (for caching)
    redis_url: str = Field(default="redis://redis:6379/0", description="Redis URL")


@lru_cache
def get_settings() -> Settings:
    """Get cached settings instance."""
    return Settings()
