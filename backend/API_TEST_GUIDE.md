# API Flow Test Guide

This document describes how to test the complete blueprint upload and analysis flow.

## Prerequisites

1. Services running (PostgreSQL, MinIO, Backend, AI Service)
2. Test data seeded (run `./seed.sh`)
3. Tools: `curl`, `jq` (optional for pretty JSON)

## Full Flow Test

### Step 1: Get Project ID

From the seed script output, note the PROJECT_ID. If you need to find it:

```bash
export PROJECT_ID="<your-project-id>"
```

Or query it from the database:
```bash
psql -h localhost -p 5432 -U platform_user -d construction_platform -c "SELECT id FROM projects LIMIT 1;"
```

### Step 2: Request Upload URL

```bash
curl -X POST http://localhost:8081/projects/$PROJECT_ID/blueprints/upload-url \
  -H 'Content-Type: application/json' \
  -d '{
    "filename": "test-blueprint.pdf",
    "content_type": "application/pdf"
  }' | jq

# Expected response:
# {
#   "blueprint_id": "uuid-here",
#   "upload_url": "http://minio:9000/blueprints/...",
#   "expires_at": "2024-01-01T00:05:00Z"
# }
```

Save the `blueprint_id` and `upload_url`:
```bash
export BLUEPRINT_ID="<blueprint-id-from-response>"
export UPLOAD_URL="<upload-url-from-response>"
```

### Step 3: Create Test File

```bash
# Create a test PDF file
echo "%PDF-1.4 Test Blueprint" > test-blueprint.pdf
```

### Step 4: Upload File to S3

```bash
curl -X PUT "$UPLOAD_URL" \
  -H "Content-Type: application/pdf" \
  --data-binary "@test-blueprint.pdf"

# Expected: Empty response with 200 OK
```

### Step 5: Complete Upload

```bash
curl -X POST http://localhost:8081/blueprints/$BLUEPRINT_ID/complete-upload | jq

# Expected response:
# {
#   "id": "uuid-here",
#   "status": "uploaded",
#   "filename": "test-blueprint.pdf"
# }
```

### Step 6: Start Analysis

```bash
curl -X POST http://localhost:8081/blueprints/$BLUEPRINT_ID/analyze | jq

# Expected response:
# {
#   "job_id": "uuid-here",
#   "status": "queued"
# }
```

Save the job ID:
```bash
export JOB_ID="<job-id-from-response>"
```

### Step 7: Check Job Status

```bash
# Poll the job status (may need to wait for worker to process)
curl http://localhost:8081/jobs/$JOB_ID | jq

# Initial response (queued):
# {
#   "id": "uuid-here",
#   "blueprint_id": "uuid-here",
#   "job_type": "takeoff",
#   "status": "queued",
#   "created_at": "...",
#   "updated_at": "..."
# }

# After processing (wait 5-10 seconds and retry):
# {
#   "id": "uuid-here",
#   "blueprint_id": "uuid-here",
#   "job_type": "takeoff",
#   "status": "completed",
#   "started_at": "...",
#   "completed_at": "...",
#   "result_data": "{...AI response...}",
#   "created_at": "...",
#   "updated_at": "..."
# }
```

### Step 8: Verify in Database

```bash
# Check blueprint status
psql -h localhost -p 5432 -U platform_user -d construction_platform \
  -c "SELECT id, filename, upload_status FROM blueprints WHERE id = '$BLUEPRINT_ID';"

# Check job status
psql -h localhost -p 5432 -U platform_user -d construction_platform \
  -c "SELECT id, status, job_type, completed_at FROM jobs WHERE id = '$JOB_ID';"
```

## Testing Error Cases

### Invalid Project ID
```bash
curl -X POST http://localhost:8081/projects/invalid-uuid/blueprints/upload-url \
  -H 'Content-Type: application/json' \
  -d '{"filename": "test.pdf", "content_type": "application/pdf"}'

# Expected: 400 Bad Request
```

### Missing Required Fields
```bash
curl -X POST http://localhost:8081/projects/$PROJECT_ID/blueprints/upload-url \
  -H 'Content-Type: application/json' \
  -d '{}'

# Expected: 400 Bad Request - "filename and content_type are required"
```

### Analyze Before Upload Complete
```bash
# Request upload URL
RESPONSE=$(curl -s -X POST http://localhost:8081/projects/$PROJECT_ID/blueprints/upload-url \
  -H 'Content-Type: application/json' \
  -d '{"filename": "test2.pdf", "content_type": "application/pdf"}')

BLUEPRINT_ID2=$(echo $RESPONSE | jq -r '.blueprint_id')

# Try to analyze without completing upload
curl -X POST http://localhost:8081/blueprints/$BLUEPRINT_ID2/analyze

# Expected: 400 Bad Request - "Blueprint must be uploaded before analysis"
```

## Health Check

```bash
curl http://localhost:8081/health | jq

# Expected response:
# {
#   "status": "ok",
#   "version": "1.0.0"
# }
```

## Monitoring Worker

Watch the backend logs to see the worker processing jobs:

```bash
docker compose logs -f backend

# Look for messages like:
# - "Processing job"
# - "Job completed successfully"
# - "Job requeued for retry" (if AI service fails)
```

## Clean Up

```bash
# Remove test files
rm test-blueprint.pdf

# Reset database (if needed)
docker compose down -v
docker compose up -d
```

## Troubleshooting

### Upload URL Not Working
- Check MinIO is running: `docker compose ps minio`
- Check MinIO console: http://localhost:9001 (login: minioadmin/minioadmin)
- Verify bucket exists in MinIO console

### Job Stuck in "queued" Status
- Check backend logs: `docker compose logs backend`
- Check AI service is running: `docker compose ps ai_service`
- Verify AI service health: `curl http://localhost:8000/health`

### Database Connection Errors
- Check PostgreSQL is running: `docker compose ps postgres`
- Verify DATABASE_URL in backend environment
- Check database exists: `docker compose exec postgres psql -U platform_user -l`

### S3 Errors
- Check MinIO logs: `docker compose logs minio`
- Verify S3 credentials in backend environment
- Test MinIO connectivity: `curl http://localhost:9000/minio/health/live`
