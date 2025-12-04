"""S3/MinIO service for file operations."""

import io
from typing import BinaryIO

import aioboto3
from botocore.exceptions import ClientError

from app.core.config import get_settings
from app.core.logging import get_logger

logger = get_logger(__name__)


class S3Service:
    """Service for interacting with S3/MinIO storage."""

    def __init__(self):
        """Initialize S3 service."""
        self.settings = get_settings()
        self.session = aioboto3.Session()

    async def download_file(self, s3_key: str) -> bytes:
        """
        Download file from S3/MinIO.

        Args:
            s3_key: S3 object key

        Returns:
            File content as bytes

        Raises:
            Exception: If download fails
        """
        try:
            async with self.session.client(
                "s3",
                endpoint_url=self.settings.s3_endpoint,
                aws_access_key_id=self.settings.s3_access_key,
                aws_secret_access_key=self.settings.s3_secret_key,
                region_name=self.settings.s3_region,
            ) as s3_client:
                logger.info("downloading_file", s3_key=s3_key, bucket=self.settings.s3_bucket)
                response = await s3_client.get_object(
                    Bucket=self.settings.s3_bucket, Key=s3_key
                )
                content = await response["Body"].read()
                logger.info(
                    "file_downloaded",
                    s3_key=s3_key,
                    size_bytes=len(content),
                )
                return content
        except ClientError as e:
            logger.error("s3_download_failed", s3_key=s3_key, error=str(e))
            raise Exception(f"Failed to download file from S3: {e}") from e
        except Exception as e:
            logger.error("s3_download_error", s3_key=s3_key, error=str(e))
            raise

    async def get_presigned_url(self, s3_key: str, expiration: int = 3600) -> str:
        """
        Get presigned URL for direct access.

        Args:
            s3_key: S3 object key
            expiration: URL expiration time in seconds

        Returns:
            Presigned URL

        Raises:
            Exception: If URL generation fails
        """
        try:
            async with self.session.client(
                "s3",
                endpoint_url=self.settings.s3_endpoint,
                aws_access_key_id=self.settings.s3_access_key,
                aws_secret_access_key=self.settings.s3_secret_key,
                region_name=self.settings.s3_region,
            ) as s3_client:
                logger.info("generating_presigned_url", s3_key=s3_key, expiration=expiration)
                url = await s3_client.generate_presigned_url(
                    "get_object",
                    Params={"Bucket": self.settings.s3_bucket, "Key": s3_key},
                    ExpiresIn=expiration,
                )
                logger.info("presigned_url_generated", s3_key=s3_key)
                return url
        except ClientError as e:
            logger.error("presigned_url_failed", s3_key=s3_key, error=str(e))
            raise Exception(f"Failed to generate presigned URL: {e}") from e
        except Exception as e:
            logger.error("presigned_url_error", s3_key=s3_key, error=str(e))
            raise

    async def upload_file(self, file_content: bytes | BinaryIO, s3_key: str) -> str:
        """
        Upload file to S3/MinIO.

        Args:
            file_content: File content as bytes or file-like object
            s3_key: S3 object key

        Returns:
            S3 key of uploaded file

        Raises:
            Exception: If upload fails
        """
        try:
            async with self.session.client(
                "s3",
                endpoint_url=self.settings.s3_endpoint,
                aws_access_key_id=self.settings.s3_access_key,
                aws_secret_access_key=self.settings.s3_secret_key,
                region_name=self.settings.s3_region,
            ) as s3_client:
                logger.info("uploading_file", s3_key=s3_key, bucket=self.settings.s3_bucket)

                if isinstance(file_content, bytes):
                    file_obj = io.BytesIO(file_content)
                else:
                    file_obj = file_content

                await s3_client.upload_fileobj(
                    file_obj,
                    self.settings.s3_bucket,
                    s3_key,
                )
                logger.info("file_uploaded", s3_key=s3_key)
                return s3_key
        except ClientError as e:
            logger.error("s3_upload_failed", s3_key=s3_key, error=str(e))
            raise Exception(f"Failed to upload file to S3: {e}") from e
        except Exception as e:
            logger.error("s3_upload_error", s3_key=s3_key, error=str(e))
            raise
