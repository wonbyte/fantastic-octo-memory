# Production Deployment Guide

Complete guide for deploying the Construction Estimation & Bidding Platform to production.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Building Docker Images](#building-docker-images)
- [Container Registry Setup](#container-registry-setup)
- [Deployment Platforms](#deployment-platforms)
  - [Option 1: AWS ECS](#option-1-aws-ecs)
  - [Option 2: Fly.io](#option-2-flyio)
  - [Option 3: Railway](#option-3-railway)
  - [Option 4: Self-Hosted with Docker Compose](#option-4-self-hosted-with-docker-compose)
- [Environment Configuration](#environment-configuration)
- [Database Setup](#database-setup)
- [Security Checklist](#security-checklist)
- [Monitoring & Maintenance](#monitoring--maintenance)
- [Troubleshooting](#troubleshooting)

---

## Prerequisites

### Required Tools
- Docker 24+ with Buildx
- Docker Compose V2
- Git
- Your preferred container registry account (GitHub, Docker Hub, AWS ECR)

### Required Credentials
- Database credentials (PostgreSQL)
- Redis password
- S3/MinIO credentials
- JWT secret key
- (Optional) Sentry DSN for error tracking

---

## Building Docker Images

### 1. Production Images

The project includes optimized production Dockerfiles:

- `backend/Dockerfile.production` - Multi-stage Go backend build
- `ai_service/Dockerfile.production` - Multi-stage Python AI service build  
- `app/Dockerfile.production` - Static web build with Nginx

### 2. Build All Images

```bash
# Build backend
docker build -f backend/Dockerfile.production -t construction-backend:latest ./backend

# Build AI service
docker build -f ai_service/Dockerfile.production -t construction-ai-service:latest ./ai_service

# Build frontend
docker build -f app/Dockerfile.production \
  --build-arg EXPO_PUBLIC_API_URL=https://api.yourdomain.com \
  --build-arg EXPO_PUBLIC_AI_SERVICE_URL=https://ai.yourdomain.com \
  -t construction-frontend:latest ./app
```

### 3. Build with Docker Compose

```bash
# Set environment variables
export VERSION=1.0.0
export DOCKER_REGISTRY=ghcr.io
export DOCKER_NAMESPACE=your-username/project-name

# Build all services
docker-compose -f docker-compose.production.yml build
```

---

## Container Registry Setup

### Option 1: GitHub Container Registry (GHCR)

**Recommended for open-source projects**

1. **Create Personal Access Token**
   - Go to GitHub Settings → Developer settings → Personal access tokens
   - Generate new token with `write:packages` scope
   - Save the token securely

2. **Login to GHCR**
   ```bash
   echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin
   ```

3. **Tag and Push Images**
   ```bash
   export REGISTRY=ghcr.io/your-username/project-name
   export VERSION=1.0.0

   # Tag images
   docker tag construction-backend:latest $REGISTRY/backend:$VERSION
   docker tag construction-ai-service:latest $REGISTRY/ai-service:$VERSION
   docker tag construction-frontend:latest $REGISTRY/frontend:$VERSION

   # Also tag as latest
   docker tag construction-backend:latest $REGISTRY/backend:latest
   docker tag construction-ai-service:latest $REGISTRY/ai-service:latest
   docker tag construction-frontend:latest $REGISTRY/frontend:latest

   # Push all tags
   docker push $REGISTRY/backend:$VERSION
   docker push $REGISTRY/backend:latest
   docker push $REGISTRY/ai-service:$VERSION
   docker push $REGISTRY/ai-service:latest
   docker push $REGISTRY/frontend:$VERSION
   docker push $REGISTRY/frontend:latest
   ```

### Option 2: Docker Hub

**Good for public projects**

1. **Login to Docker Hub**
   ```bash
   docker login
   ```

2. **Tag and Push**
   ```bash
   export REGISTRY=docker.io/your-username
   # Follow similar tagging/pushing steps as GHCR
   ```

### Option 3: AWS ECR

**Best for AWS deployments**

1. **Create Repositories**
   ```bash
   aws ecr create-repository --repository-name construction/backend
   aws ecr create-repository --repository-name construction/ai-service
   aws ecr create-repository --repository-name construction/frontend
   ```

2. **Login to ECR**
   ```bash
   aws ecr get-login-password --region us-east-1 | \
     docker login --username AWS --password-stdin \
     ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com
   ```

3. **Tag and Push**
   ```bash
   export REGISTRY=ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com
   # Follow similar tagging/pushing steps
   ```

---

## Deployment Platforms

### Option 1: AWS ECS

**Best for: Enterprise deployments with high scalability needs**

#### Prerequisites
- AWS Account with ECS access
- AWS CLI configured
- Images pushed to ECR

#### Steps

1. **Create ECS Cluster**
   ```bash
   aws ecs create-cluster --cluster-name construction-platform
   ```

2. **Create Task Definitions**

   Create `ecs-task-backend.json`:
   ```json
   {
     "family": "construction-backend",
     "networkMode": "awsvpc",
     "requiresCompatibilities": ["FARGATE"],
     "cpu": "512",
     "memory": "1024",
     "containerDefinitions": [
       {
         "name": "backend",
         "image": "ACCOUNT_ID.dkr.ecr.REGION.amazonaws.com/construction/backend:latest",
         "portMappings": [
           {
             "containerPort": 8080,
             "protocol": "tcp"
           }
         ],
         "environment": [
           {"name": "ENV", "value": "production"},
           {"name": "PORT", "value": "8080"}
         ],
         "secrets": [
           {
             "name": "DATABASE_URL",
             "valueFrom": "arn:aws:secretsmanager:REGION:ACCOUNT:secret:construction/database-url"
           },
           {
             "name": "JWT_SECRET",
             "valueFrom": "arn:aws:secretsmanager:REGION:ACCOUNT:secret:construction/jwt-secret"
           }
         ],
         "logConfiguration": {
           "logDriver": "awslogs",
           "options": {
             "awslogs-group": "/ecs/construction-backend",
             "awslogs-region": "us-east-1",
             "awslogs-stream-prefix": "ecs"
           }
         }
       }
     ]
   }
   ```

3. **Register Task Definitions**
   ```bash
   aws ecs register-task-definition --cli-input-json file://ecs-task-backend.json
   aws ecs register-task-definition --cli-input-json file://ecs-task-ai-service.json
   aws ecs register-task-definition --cli-input-json file://ecs-task-frontend.json
   ```

4. **Create Services**
   ```bash
   aws ecs create-service \
     --cluster construction-platform \
     --service-name backend \
     --task-definition construction-backend \
     --desired-count 2 \
     --launch-type FARGATE \
     --network-configuration "awsvpcConfiguration={subnets=[subnet-xxx],securityGroups=[sg-xxx],assignPublicIp=ENABLED}"
   ```

5. **Configure Application Load Balancer**
   - Create ALB targeting ECS services
   - Configure health checks pointing to `/health` endpoints
   - Set up SSL certificates

6. **Database Setup**
   - Use Amazon RDS for PostgreSQL
   - Configure security groups for ECS access
   - Store connection string in AWS Secrets Manager

---

### Option 2: Fly.io

**Best for: Quick MVP deployment with global edge network**

#### Prerequisites
- Fly.io account
- Fly CLI installed (`brew install flyctl` or download from fly.io)

#### Steps

1. **Login to Fly.io**
   ```bash
   flyctl auth login
   ```

2. **Create Fly Apps**
   ```bash
   # Backend
   cd backend
   flyctl apps create construction-backend

   # AI Service
   cd ../ai_service
   flyctl apps create construction-ai-service

   # Frontend
   cd ../app
   flyctl apps create construction-frontend
   ```

3. **Create fly.toml for Backend**

   Create `backend/fly.toml`:
   ```toml
   app = "construction-backend"
   primary_region = "iad"

   [build]
     dockerfile = "Dockerfile.production"

   [env]
     PORT = "8080"
     ENV = "production"

   [http_service]
     internal_port = 8080
     force_https = true
     auto_stop_machines = false
     auto_start_machines = true
     min_machines_running = 1

   [[http_service.checks]]
     grace_period = "10s"
     interval = "30s"
     method = "GET"
     timeout = "5s"
     path = "/health"

   [[vm]]
     cpu_kind = "shared"
     cpus = 1
     memory_mb = 1024
   ```

4. **Set Secrets**
   ```bash
   # Backend secrets
   flyctl secrets set \
     DATABASE_URL="postgres://..." \
     JWT_SECRET="..." \
     REDIS_URL="redis://..." \
     --app construction-backend

   # AI Service secrets
   flyctl secrets set \
     DATABASE_URL="postgres://..." \
     REDIS_URL="redis://..." \
     --app construction-ai-service
   ```

5. **Create PostgreSQL Database**
   ```bash
   flyctl postgres create --name construction-db
   flyctl postgres attach construction-db --app construction-backend
   ```

6. **Deploy Services**
   ```bash
   # Backend
   cd backend
   flyctl deploy

   # AI Service
   cd ../ai_service
   flyctl deploy

   # Frontend
   cd ../app
   flyctl deploy
   ```

7. **Scale as Needed**
   ```bash
   flyctl scale count 2 --app construction-backend
   flyctl scale vm shared-cpu-2x --memory 2048 --app construction-ai-service
   ```

---

### Option 3: Railway

**Best for: Simple deployment with managed database**

#### Prerequisites
- Railway account
- Railway CLI installed

#### Steps

1. **Login to Railway**
   ```bash
   railway login
   ```

2. **Create New Project**
   ```bash
   railway init
   ```

3. **Add PostgreSQL Database**
   ```bash
   railway add --database postgresql
   ```

4. **Add Redis**
   ```bash
   railway add --database redis
   ```

5. **Deploy Backend**
   ```bash
   cd backend
   railway up -d Dockerfile.production
   ```

6. **Set Environment Variables**
   - Go to Railway dashboard
   - Add environment variables for each service
   - Use Railway's built-in `DATABASE_URL` and `REDIS_URL`

7. **Deploy AI Service and Frontend**
   ```bash
   cd ../ai_service
   railway up -d Dockerfile.production

   cd ../app
   railway up -d Dockerfile.production
   ```

8. **Configure Domains**
   - Add custom domains in Railway dashboard
   - Update EXPO_PUBLIC_API_URL to point to backend domain

---

### Option 4: Self-Hosted with Docker Compose

**Best for: Full control, private infrastructure**

#### Prerequisites
- Linux server (Ubuntu 22.04+ recommended)
- Docker and Docker Compose installed
- Domain names pointed to your server
- SSL certificates (use Let's Encrypt)

#### Steps

1. **Prepare Server**
   ```bash
   # Update system
   sudo apt update && sudo apt upgrade -y

   # Install Docker
   curl -fsSL https://get.docker.com -o get-docker.sh
   sudo sh get-docker.sh

   # Install Docker Compose V2
   sudo apt install docker-compose-plugin -y
   ```

2. **Clone Repository**
   ```bash
   git clone https://github.com/wonbyte/fantastic-octo-memory.git
   cd fantastic-octo-memory
   ```

3. **Configure Environment**
   ```bash
   cp .env.production.example .env.production
   nano .env.production
   # Fill in all required values
   ```

4. **Deploy with Docker Compose**
   ```bash
   docker-compose -f docker-compose.production.yml --env-file .env.production up -d
   ```

5. **Set Up Reverse Proxy (Nginx)**

   Create `/etc/nginx/sites-available/construction-platform`:
   ```nginx
   server {
       listen 80;
       server_name api.yourdomain.com;
       
       location / {
           proxy_pass http://localhost:8081;
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
           proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
           proxy_set_header X-Forwarded-Proto $scheme;
       }
   }

   server {
       listen 80;
       server_name ai.yourdomain.com;
       
       location / {
           proxy_pass http://localhost:8000;
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
       }
   }

   server {
       listen 80;
       server_name yourdomain.com www.yourdomain.com;
       
       location / {
           proxy_pass http://localhost:80;
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
       }
   }
   ```

6. **Enable SSL with Certbot**
   ```bash
   sudo apt install certbot python3-certbot-nginx -y
   sudo certbot --nginx -d api.yourdomain.com -d ai.yourdomain.com -d yourdomain.com
   ```

7. **Set Up Auto-Restart**
   ```bash
   # Create systemd service
   sudo nano /etc/systemd/system/construction-platform.service
   ```

   ```ini
   [Unit]
   Description=Construction Platform
   Requires=docker.service
   After=docker.service

   [Service]
   Type=oneshot
   RemainAfterExit=yes
   WorkingDirectory=/path/to/fantastic-octo-memory
   ExecStart=/usr/bin/docker-compose -f docker-compose.production.yml --env-file .env.production up -d
   ExecStop=/usr/bin/docker-compose -f docker-compose.production.yml down
   
   [Install]
   WantedBy=multi-user.target
   ```

   ```bash
   sudo systemctl enable construction-platform
   sudo systemctl start construction-platform
   ```

---

## Environment Configuration

### Required Environment Variables

**Security-Critical Variables:**
```bash
# MUST be changed from defaults
JWT_SECRET=<generate-with-openssl-rand-base64-32>
POSTGRES_PASSWORD=<strong-password>
REDIS_PASSWORD=<strong-password>
MINIO_ROOT_USER=<strong-username>
MINIO_ROOT_PASSWORD=<strong-password>
```

**Database Configuration:**
```bash
POSTGRES_DB=construction_platform_prod
POSTGRES_USER=platform_user_prod
DATABASE_URL=postgres://user:pass@host:5432/db?sslmode=require
```

**Service URLs:**
```bash
# Internal (container-to-container)
BACKEND_URL=http://backend:8080
AI_SERVICE_URL=http://ai_service:8000

# External (for frontend)
EXPO_PUBLIC_API_URL=https://api.yourdomain.com
EXPO_PUBLIC_AI_SERVICE_URL=https://ai.yourdomain.com
```

### Generating Secure Secrets

```bash
# Generate JWT secret
openssl rand -base64 32

# Generate strong password
openssl rand -base64 24
```

---

## Database Setup

### Initial Migration

When deploying for the first time:

```bash
# The backend entrypoint.sh automatically runs migrations
# But you can also run manually:

docker exec construction-backend-prod migrate \
  -path /app/migrations \
  -database "$DATABASE_URL" \
  up
```

### Seed Data (Optional)

```bash
# Run seed script if needed
docker exec construction-backend-prod /app/seed.sh
```

### Backup Strategy

**Automated Daily Backups:**

```bash
#!/bin/bash
# backup-db.sh

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR=/backups/postgres
mkdir -p $BACKUP_DIR

docker exec construction-postgres-prod pg_dump \
  -U $POSTGRES_USER \
  -d $POSTGRES_DB \
  -F c \
  -b \
  -v \
  -f /tmp/backup_$DATE.dump

docker cp construction-postgres-prod:/tmp/backup_$DATE.dump \
  $BACKUP_DIR/backup_$DATE.dump

# Keep only last 7 days
find $BACKUP_DIR -name "backup_*.dump" -mtime +7 -delete
```

**Schedule with cron:**
```bash
0 2 * * * /path/to/backup-db.sh
```

---

## Security Checklist

### Before Going Live

- [ ] **Secrets Management**
  - [ ] All default passwords changed
  - [ ] JWT_SECRET is randomly generated and secure
  - [ ] Database credentials are strong and unique
  - [ ] Secrets stored in environment variables, not in code
  - [ ] .env.production added to .gitignore

- [ ] **Network Security**
  - [ ] HTTPS/TLS enabled on all public endpoints
  - [ ] Database not exposed to public internet
  - [ ] Redis not exposed to public internet
  - [ ] Firewall rules configured properly

- [ ] **Container Security**
  - [ ] Running containers as non-root users
  - [ ] Production Dockerfiles use minimal base images
  - [ ] No unnecessary packages installed
  - [ ] Security updates applied to base images

- [ ] **Application Security**
  - [ ] CORS configured correctly
  - [ ] Rate limiting enabled
  - [ ] SQL injection protection (using parameterized queries)
  - [ ] XSS protection headers set
  - [ ] Content Security Policy configured

- [ ] **Database Security**
  - [ ] SSL/TLS required for connections (sslmode=require)
  - [ ] Least privilege access (application user can't drop tables)
  - [ ] Regular backups configured
  - [ ] Backup restoration tested

- [ ] **Monitoring**
  - [ ] Sentry or error tracking configured
  - [ ] Health check endpoints working
  - [ ] Log aggregation set up
  - [ ] Alerts configured for critical errors

### Security Headers

Ensure these headers are set (handled by nginx.conf):

```
X-Frame-Options: SAMEORIGIN
X-Content-Type-Options: nosniff
X-XSS-Protection: 1; mode=block
Referrer-Policy: strict-origin-when-cross-origin
Content-Security-Policy: default-src 'self'
```

---

## Monitoring & Maintenance

### Health Checks

All services expose `/health` endpoints:

```bash
# Backend
curl https://api.yourdomain.com/health

# AI Service
curl https://ai.yourdomain.com/health

# Frontend
curl https://yourdomain.com/health
```

### Log Monitoring

**View logs:**
```bash
# All services
docker-compose -f docker-compose.production.yml logs -f

# Specific service
docker-compose -f docker-compose.production.yml logs -f backend

# Last 100 lines
docker-compose -f docker-compose.production.yml logs --tail=100 backend
```

**Set up log rotation:**
```json
// /etc/docker/daemon.json
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  }
}
```

### Performance Monitoring

**Resource Usage:**
```bash
docker stats
```

**Database Performance:**
```bash
docker exec -it construction-postgres-prod psql -U $POSTGRES_USER -d $POSTGRES_DB

-- Check slow queries
SELECT * FROM pg_stat_statements ORDER BY mean_exec_time DESC LIMIT 10;

-- Check database size
SELECT pg_size_pretty(pg_database_size('construction_platform_prod'));
```

### Updates and Rollbacks

**Update to new version:**
```bash
# Pull new images
docker-compose -f docker-compose.production.yml pull

# Restart services with zero downtime
docker-compose -f docker-compose.production.yml up -d
```

**Rollback to previous version:**
```bash
# Use specific version tag
export VERSION=1.0.0
docker-compose -f docker-compose.production.yml up -d
```

---

## Troubleshooting

### Common Issues

**Issue: Backend won't start**
```bash
# Check logs
docker logs construction-backend-prod

# Common causes:
# - Database not ready (wait for postgres health check)
# - Missing JWT_SECRET environment variable
# - Migration failures
```

**Issue: Database connection failed**
```bash
# Verify database is running
docker ps | grep postgres

# Check database logs
docker logs construction-postgres-prod

# Test connection
docker exec construction-backend-prod pg_isready -h postgres -U $POSTGRES_USER
```

**Issue: Frontend can't connect to backend**
```bash
# Check EXPO_PUBLIC_API_URL is correct
docker exec construction-frontend-prod env | grep EXPO_PUBLIC

# Verify backend is accessible
curl https://api.yourdomain.com/health
```

**Issue: Out of disk space**
```bash
# Check disk usage
df -h

# Clean up Docker
docker system prune -a --volumes

# Check image sizes
docker images
```

### Getting Help

1. Check service logs
2. Verify environment variables
3. Check health endpoints
4. Review security groups/firewall rules
5. Test database connectivity
6. Check DNS resolution

---

## Additional Resources

- [Docker Production Best Practices](https://docs.docker.com/develop/dev-best-practices/)
- [AWS ECS Documentation](https://docs.aws.amazon.com/ecs/)
- [Fly.io Documentation](https://fly.io/docs/)
- [Railway Documentation](https://docs.railway.app/)
- [Postgres Performance Tuning](https://wiki.postgresql.org/wiki/Performance_Optimization)

---

## Next Steps

After successful deployment:

1. ✅ Run E2E validation tests (see [E2E_TESTING.md](./E2E_TESTING.md))
2. ✅ Set up monitoring and alerts
3. ✅ Configure automated backups
4. ✅ Document incident response procedures
5. ✅ Set up CI/CD for automated deployments
6. ✅ Plan for scaling strategy
