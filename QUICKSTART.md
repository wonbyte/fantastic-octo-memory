# Quick Start: Production Deployment

This is a quick reference guide for deploying to production. For complete details, see [DEPLOYMENT.md](./DEPLOYMENT.md).

## Prerequisites

- Docker 24+ installed
- Docker Compose V2 installed
- Container registry account (GitHub, Docker Hub, or AWS ECR)
- Production server or cloud platform account

## 1. Configure Environment

```bash
# Copy production environment template
cp .env.production.example .env.production

# Edit with your production values
nano .env.production
```

**Critical values to change:**
- `JWT_SECRET` - Generate with: `openssl rand -base64 32`
- `POSTGRES_PASSWORD` - Strong password
- `REDIS_PASSWORD` - Strong password
- `MINIO_ROOT_USER` and `MINIO_ROOT_PASSWORD` - Strong credentials
- `EXPO_PUBLIC_API_URL` - Your production API URL
- `EXPO_PUBLIC_AI_SERVICE_URL` - Your production AI service URL

## 2. Build Production Images

```bash
# Build all production images locally
make prod-build

# Or build and push to container registry
./scripts/push-images.sh --version 1.0.0 \
  --api-url https://api.yourdomain.com \
  --ai-url https://ai.yourdomain.com
```

## 3. Deploy

### Option A: Docker Compose (Simple Self-Hosted)

```bash
# Start all services
make prod-start

# View logs
make prod-logs

# Stop services
make prod-stop
```

### Option B: Cloud Platforms

Choose your platform and follow the detailed guide:

- **AWS ECS**: See [DEPLOYMENT.md - AWS ECS](./DEPLOYMENT.md#option-1-aws-ecs)
- **Fly.io**: See [DEPLOYMENT.md - Fly.io](./DEPLOYMENT.md#option-2-flyio)
- **Railway**: See [DEPLOYMENT.md - Railway](./DEPLOYMENT.md#option-3-railway)

## 4. Validate Deployment

### Health Checks

```bash
# Check all services are healthy
curl https://api.yourdomain.com/health
curl https://ai.yourdomain.com/health
curl https://yourdomain.com/health
```

### E2E Testing

Follow the complete E2E test procedure in [E2E_TESTING.md](./E2E_TESTING.md):

1. ✅ User signup/login
2. ✅ Project creation
3. ✅ Blueprint upload
4. ✅ AI analysis
5. ✅ Bid generation
6. ✅ PDF download

## 5. Monitoring

### View Logs

```bash
# All services
make prod-logs

# Specific service
docker logs construction-backend-prod -f
docker logs construction-ai-service-prod -f
docker logs construction-frontend-prod -f
```

### Check Resource Usage

```bash
docker stats
```

## 6. Common Issues

### Services Won't Start

```bash
# Check if .env.production exists
ls -la .env.production

# Verify environment variables
docker-compose -f docker-compose.production.yml --env-file .env.production config

# Check logs for errors
docker logs construction-backend-prod
```

### Database Connection Failed

```bash
# Check if PostgreSQL is running
docker ps | grep postgres

# Test database connection
docker exec construction-backend-prod pg_isready -h postgres -U $POSTGRES_USER
```

### Frontend Can't Reach Backend

1. Verify `EXPO_PUBLIC_API_URL` is set correctly in `.env.production`
2. Check CORS configuration
3. Verify backend is accessible: `curl https://api.yourdomain.com/health`

## 7. Backup & Restore

### Backup Database

```bash
docker exec construction-postgres-prod pg_dump \
  -U $POSTGRES_USER -d $POSTGRES_DB \
  -F c -b -v \
  -f /tmp/backup.dump
  
docker cp construction-postgres-prod:/tmp/backup.dump ./backup.dump
```

### Restore Database

```bash
docker cp ./backup.dump construction-postgres-prod:/tmp/backup.dump

docker exec construction-postgres-prod pg_restore \
  -U $POSTGRES_USER -d $POSTGRES_DB \
  -c -v /tmp/backup.dump
```

## 8. Update Deployment

### Pull New Version

```bash
# Pull new images
export VERSION=1.1.0
docker-compose -f docker-compose.production.yml pull

# Restart services (zero downtime)
docker-compose -f docker-compose.production.yml up -d
```

### Rollback

```bash
# Use previous version
export VERSION=1.0.0
docker-compose -f docker-compose.production.yml up -d
```

## Security Checklist

Before going live:

- [ ] All default passwords changed
- [ ] JWT_SECRET is randomly generated
- [ ] HTTPS/TLS configured
- [ ] Firewall rules configured
- [ ] Database not publicly accessible
- [ ] Redis not publicly accessible
- [ ] Backups configured and tested
- [ ] Error tracking (Sentry) configured
- [ ] Monitoring alerts set up

## Support

For detailed documentation:
- [DEPLOYMENT.md](./DEPLOYMENT.md) - Complete deployment guide
- [E2E_TESTING.md](./E2E_TESTING.md) - Testing procedures
- [M7_IMPLEMENTATION_SUMMARY.md](./M7_IMPLEMENTATION_SUMMARY.md) - Implementation details

For issues: Open a GitHub issue or contact support.
