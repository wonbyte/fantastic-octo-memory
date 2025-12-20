# Backend Service

Go-based backend API for the Construction Estimation & Bidding Automation platform.

## Architecture

This backend follows clean architecture principles with the following structure:

```
/backend
  /cmd/server/         # Application entry point
  /internal/
    /config/          # Configuration management
    /handlers/        # HTTP request handlers
    /middleware/      # HTTP middleware (logging, CORS, recovery)
    /models/          # Domain models
    /repository/      # Database access layer
    /services/        # Business logic and external services
  /migrations/        # Database migrations
  /pkg/               # Shared utilities (future use)
```

## Features

- ✅ Clean architecture with separation of concerns
- ✅ PostgreSQL with pgx connection pooling
- ✅ S3-compatible storage (MinIO for development)
- ✅ Database migrations with golang-migrate
- ✅ Structured logging with slog
- ✅ Graceful shutdown
- ✅ Health check endpoint
- ✅ Async job worker for AI processing
- ✅ Pre-signed URL generation for file uploads

## Prerequisites

- Go 1.22+
- PostgreSQL 16+
- MinIO or S3-compatible storage
- (Optional) Docker and Docker Compose

## Quick Start

### Using Docker Compose

```bash
# From repository root
docker compose up backend
```

### Local Development

1. **Install dependencies**
   ```bash
   go mod download
   ```

2. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your local configuration
   ```

3. **Start PostgreSQL and MinIO**
   ```bash
   # Using Docker Compose (from repo root)
   docker compose up postgres minio
   ```

4. **Run migrations**
   ```bash
   export DATABASE_URL="postgres://platform_user:change_me_in_production@localhost:5432/construction_platform?sslmode=disable"
   migrate -path migrations -database "$DATABASE_URL" up
   ```

5. **Start the server**
   ```bash
   go run cmd/server/main.go
   ```

The server will start on `http://localhost:8080` (or port 8081 when running with docker-compose for development)

## API Endpoints

### Health Check
```http
GET /health
```

Returns service health status and version.

### Blueprint Upload Flow

#### 1. Request Upload URL
```http
POST /projects/{id}/blueprints/upload-url
Content-Type: application/json

{
  "filename": "floor-plan.pdf",
  "content_type": "application/pdf"
}
```

Response:
```json
{
  "blueprint_id": "uuid",
  "upload_url": "https://...",
  "expires_at": "2024-01-01T00:00:00Z"
}
```

#### 2. Upload File (Client-side)
Use the returned `upload_url` to upload the file directly to S3/MinIO.

```bash
curl -X PUT -T file.pdf \
  -H "Content-Type: application/pdf" \
  "<upload_url>"
```

#### 3. Complete Upload
```http
POST /blueprints/{id}/complete-upload
```

Response:
```json
{
  "id": "uuid",
  "status": "uploaded",
  "filename": "floor-plan.pdf"
}
```

### Job Processing

#### Start Analysis
```http
POST /blueprints/{id}/analyze
```

Response:
```json
{
  "job_id": "uuid",
  "status": "queued"
}
```

#### Check Job Status
```http
GET /jobs/{id}
```

Response:
```json
{
  "id": "uuid",
  "blueprint_id": "uuid",
  "job_type": "takeoff",
  "status": "completed",
  "started_at": "2024-01-01T00:00:00Z",
  "completed_at": "2024-01-01T00:01:00Z",
  "result_data": "{...}",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:01:00Z"
}
```

## Environment Variables

See `.env.example` for all available configuration options.

Key variables:
- `PORT` - Server port (default: 8080)
- `DATABASE_URL` - PostgreSQL connection string
- `S3_ENDPOINT` - S3/MinIO endpoint
- `S3_BUCKET` - S3 bucket name
- `AI_SERVICE_URL` - AI service endpoint
- `JOB_POLL_INTERVAL` - Worker polling interval

## Database Migrations

### Create a new migration
```bash
migrate create -ext sql -dir migrations -seq description_of_migration
```

### Run migrations
```bash
migrate -path migrations -database "$DATABASE_URL" up
```

### Rollback migration
```bash
migrate -path migrations -database "$DATABASE_URL" down 1
```

## Testing

### Run all tests
```bash
go test -v ./...
```

### Run with coverage
```bash
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Run specific tests
```bash
go test -v ./internal/handlers
```

## Development Tools

### Seed Test Data
```bash
./seed.sh
```

### Format Code
```bash
go fmt ./...
```

### Lint Code
```bash
go vet ./...
```

### Build Binary
```bash
go build -o bin/server ./cmd/server
```

## Worker Process

The backend includes an async worker that:
1. Polls for queued jobs every 5 seconds (configurable)
2. Processes jobs by calling the AI service
3. Retries failed jobs up to 3 times (configurable)
4. Updates job status and stores results

The worker runs in the same process as the HTTP server and handles graceful shutdown.

## Architecture Decisions

### Why Chi Router?
- Lightweight and idiomatic to Go's net/http
- Good performance
- Minimal dependencies
- Context-based middleware

### Why pgx over database/sql?
- Better performance with PostgreSQL
- Built-in connection pooling
- Support for PostgreSQL-specific features (JSONB, arrays, etc.)
- Type-safe queries

### Why slog for logging?
- Standard library (Go 1.21+)
- Structured logging out of the box
- Good performance
- Zero dependencies

## Production Considerations

Before deploying to production:

1. Change all default passwords and secrets
2. Use a managed PostgreSQL instance
3. Use a managed S3 service (AWS S3, GCS, etc.)
4. Enable TLS/SSL for all connections
5. Set up proper monitoring and alerting
6. Configure proper CORS origins
7. Set up rate limiting
8. Use a reverse proxy (nginx, Traefik, etc.)
9. Enable database connection pooling tuning
10. Review and adjust worker configuration

## Troubleshooting

### Database connection issues
- Verify PostgreSQL is running
- Check DATABASE_URL is correct
- Ensure database exists
- Verify network connectivity

### S3/MinIO connection issues
- Verify MinIO is running
- Check S3_ENDPOINT is accessible
- Verify credentials are correct
- Check bucket exists

### Migration failures
- Check migration files are valid SQL
- Verify database permissions
- Review migration logs
- Use `down` migrations to rollback

## License

[Your License Here]
