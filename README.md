# Construction Estimation & Bidding Automation SaaS Platform

[![CI Status](https://github.com/wonbyte/fantastic-octo-memory/workflows/CI/badge.svg)](https://github.com/wonbyte/fantastic-octo-memory/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/wonbyte/fantastic-octo-memory/branch/main/graph/badge.svg)](https://codecov.io/gh/wonbyte/fantastic-octo-memory)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go)](https://golang.org/)
[![Python Version](https://img.shields.io/badge/Python-3.12+-3776AB?logo=python)](https://www.python.org/)
[![Node Version](https://img.shields.io/badge/Node-22_LTS-339933?logo=node.js)](https://nodejs.org/)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

A comprehensive platform for automating construction estimation and bidding processes, built as a modern cloud-native application with AI-powered capabilities.

## 🏗️ Project Overview

This platform streamlines the construction bidding process by:
- Automating cost estimations using AI and historical data
- Managing bid submissions and tracking
- Providing real-time analytics and insights
- Integrating with existing construction management tools

## 📁 Monorepo Structure

```
.
├── backend/          # Go backend API service (Go 1.25+)
├── ai_service/       # Python FastAPI AI/ML service (Python 3.12+)
├── app/              # React Native + Web frontend (Expo SDK 52+)
├── infra/            # Docker, CI/CD, and infrastructure configs
├── .github/          # GitHub Actions workflows
└── README.md         # This file
```

### Directory Details

- **`/backend`**: Core business logic API built with Go
  - RESTful API endpoints
  - PostgreSQL integration
  - Redis caching
  - Authentication and authorization

- **`/ai_service`**: AI/ML service built with Python FastAPI
  - Cost estimation models
  - Document processing
  - Prediction APIs
  - Model training pipelines

- **`/app`**: Cross-platform mobile and web application
  - React Native for mobile (iOS/Android)
  - Expo for unified development
  - Web support via React Native Web
  - Shared codebase across platforms

- **`/infra`**: Infrastructure and deployment configurations
  - Docker configurations
  - CI/CD pipelines
  - Kubernetes manifests (future)

## 🚀 Quick Start

### Prerequisites

- **Docker** (with Docker Compose V2)
- **Make** (for convenience commands)
- **Node.js** 22 LTS (for local app development)
- **Go** 1.25+ (for local backend development)
- **Python** 3.12+ (for local AI service development)

### Running the Platform

1. **Validate your setup (recommended)**
   ```bash
   ./scripts/validate-setup.sh
   ```

2. **Clone the repository**
   ```bash
   git clone https://github.com/wonbyte/fantastic-octo-memory.git
   cd fantastic-octo-memory
   ```

3. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Start all services**
   ```bash
   make dev
   ```

   This will start:
   - Backend API at http://localhost:8080
   - AI Service at http://localhost:8000
   - Frontend at http://localhost:3000
   - PostgreSQL at localhost:5432
   - Redis at localhost:6379

4. **Access the application**
   - Web App: http://localhost:3000
   - Backend API: http://localhost:8080
   - AI Service: http://localhost:8000
   - API Docs: http://localhost:8000/docs (FastAPI)

### Available Make Commands

```bash
make dev      # Start all services in development mode with hot-reload
make start    # Start all services in production mode
make stop     # Stop all running services
make build    # Build all Docker images
make clean    # Remove containers, volumes, and clean up
make logs     # Tail logs from all services
make test     # Run tests for all services
```

## 🛠️ Technology Stack

| Component | Technology | Version |
|-----------|-----------|---------|
| **Backend** | Go | 1.25+ |
| **AI Service** | Python | 3.12+ |
| **AI Framework** | FastAPI | 0.115+ |
| **Frontend** | React | 19.0+ |
| **Mobile** | React Native | 0.82 |
| **Mobile Framework** | Expo SDK | 54.0+ |
| **Runtime** | Node.js | 24 LTS |
| **Database** | PostgreSQL | 16+ |
| **Cache** | Redis | 7.4+ |
| **Container** | Docker | Latest |
| **Orchestration** | Docker Compose | V2 |
| **Base Image** | Alpine Linux | 3.19+ |

## 🔧 Development

### Backend Development

```bash
cd backend
go mod download
go run main.go
```

### AI Service Development

```bash
cd ai_service
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate
pip install -r requirements.txt
uvicorn main:app --reload
```

### App Development

```bash
cd app
npm install
npm start
```

## 🧪 Testing

### Run all tests
```bash
make test
```

### Run tests for specific service
```bash
# Backend
cd backend && go test ./...

# AI Service
cd ai_service && pytest

# App
cd app && npm test
```

## 🔒 Environment Configuration

Copy the `.env.example` files to `.env` in each directory and configure:

- **Root `.env`**: Shared configuration
- **Backend `.env`**: Database URLs, API keys
- **AI Service `.env`**: Model paths, ML configurations
- **App `.env`**: API endpoints, feature flags

## 📝 Code Quality

### Linting

```bash
# Go
cd backend && go vet ./...

# Python
cd ai_service && ruff check .

# TypeScript/JavaScript
cd app && npm run lint
```

### Formatting

```bash
# Go
cd backend && go fmt ./...

# Python
cd ai_service && ruff format .

# TypeScript/JavaScript
cd app && npm run format
```

## 🔄 CI/CD Pipelines

The platform uses GitHub Actions for comprehensive CI/CD automation across all services in the monorepo.

### Continuous Integration (CI)

**Workflow**: `.github/workflows/ci.yml`

Automatically runs on every push and pull request to `main` and `develop` branches:

#### Backend (Go)
- ✅ Dependency caching
- ✅ `go vet` linting
- ✅ Unit tests with race detection
- ✅ Integration tests
- ✅ Coverage reports to Codecov

#### AI Service (Python)
- ✅ Dependency caching
- ✅ Ruff linting
- ✅ Unit tests with pytest
- ✅ Integration tests
- ✅ Coverage reports to Codecov

#### Frontend (React Native/Expo)
- ✅ Dependency installation
- ✅ ESLint linting
- ✅ TypeScript type checking
- ✅ Unit tests with Jest
- ✅ Coverage reports to Codecov

#### E2E Tests (Playwright)
- ✅ Full-stack end-to-end tests
- ✅ Test artifacts and screenshots
- ✅ Test reports

#### Docker Builds
- ✅ Build verification for all Docker images
- ✅ Build caching for faster builds
- ✅ Multi-service validation

### Build and Push (Production Images)

**Workflow**: `.github/workflows/build-production.yml`

Builds and pushes production-optimized Docker images to GitHub Container Registry (GHCR):

**Triggers**:
- 🏷️ Automatically on version tags (e.g., `v1.0.0`)
- 🔘 Manual trigger with custom version tag

**What it does**:
1. Builds optimized production images using multi-stage Dockerfiles
2. Pushes to `ghcr.io/$GITHUB_REPOSITORY/[backend|ai-service|frontend]`
3. Tags images with both version and `latest`
4. Uses build caching for efficiency

**Usage**:
```bash
# Automatically triggered by creating a version tag
git tag v1.0.0
git push origin v1.0.0

# Or trigger manually from GitHub Actions UI:
# Actions → Build and Push Production Images → Run workflow
# Enter version: v1.0.0 (or latest)
```

### Deployment Workflows

#### Deploy to Staging

**Workflow**: `.github/workflows/deploy-staging.yml`

Deploys to staging environment for testing before production.

**Features**:
- 🔘 Manual trigger with version selection
- ✅ Health checks after deployment
- ✅ Smoke tests
- ✅ E2E tests against staging
- ✅ Deployment summaries

**Usage**:
```bash
# From GitHub Actions UI:
# Actions → Deploy to Staging → Run workflow
# Select version: v1.0.0 or latest
# Optional: Skip health checks (for debugging)
```

**Environment Variables** (Configure in GitHub Settings → Environments → staging):
- `STAGING_URL` - Frontend URL (e.g., https://staging.example.com)
- `STAGING_API_URL` - Backend API URL
- `STAGING_AI_URL` - AI service URL

#### Deploy to Production

**Workflow**: `.github/workflows/deploy-production.yml`

Production deployment with safety checks and rollback capabilities.

**Features**:
- 🔘 Manual trigger only (safety measure)
- ✅ Pre-deployment validation
- ✅ Version format checking (requires semantic versioning)
- ✅ Backup creation before deployment
- ✅ Blue-green/rolling deployment support
- ✅ Database migration handling
- ✅ Comprehensive health checks
- ✅ Smoke tests on production
- ✅ Error rate monitoring
- ✅ Automatic rollback on failure
- ✅ Post-deployment verification

**Usage**:
```bash
# IMPORTANT: Always deploy to staging first!
# 
# From GitHub Actions UI:
# Actions → Deploy to Production → Run workflow
# Enter version: v1.0.0 (must be semantic version)
# Requires manual approval if enabled
```

**Environment Variables** (Configure in GitHub Settings → Environments → production):
- `PRODUCTION_URL` - Frontend URL (e.g., https://app.example.com)
- `PRODUCTION_API_URL` - Backend API URL
- `PRODUCTION_AI_URL` - AI service URL

### Complete CI/CD Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                     Development Process                          │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│  1. Push/PR to main/develop → CI Workflow                       │
│     ✓ Lint all services                                         │
│     ✓ Run unit tests                                            │
│     ✓ Run integration tests                                     │
│     ✓ Build Docker images                                       │
│     ✓ Run E2E tests                                             │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│  2. Create version tag (e.g., v1.0.0) → Build Production        │
│     ✓ Build optimized images                                    │
│     ✓ Push to GHCR                                              │
│     ✓ Tag with version + latest                                 │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│  3. Manual Deploy to Staging                                    │
│     ✓ Deploy version to staging                                 │
│     ✓ Run health checks                                         │
│     ✓ Run smoke tests                                           │
│     ✓ Run E2E tests                                             │
│     ✓ Manual validation                                         │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│  4. Manual Deploy to Production                                 │
│     ✓ Pre-deployment checks                                     │
│     ✓ Create backups                                            │
│     ✓ Deploy with zero-downtime                                 │
│     ✓ Run migrations                                            │
│     ✓ Health checks                                             │
│     ✓ Smoke tests                                               │
│     ✓ Monitor error rates                                       │
│     ✓ Post-deployment verification                              │
└─────────────────────────────────────────────────────────────────┘
```

### Setting Up Deployments

The deployment workflows are **placeholder implementations** that demonstrate best practices. To make them functional:

1. **Choose your deployment target**:
   - Docker Compose on VPS/VM
   - Kubernetes (EKS, GKE, AKS)
   - AWS ECS/Fargate
   - Fly.io
   - Railway
   - Render
   - Other cloud platforms

2. **Configure environment variables** in GitHub:
   - Go to: Settings → Environments
   - Create environments: `staging` and `production`
   - Add environment variables (URLs, API keys, etc.)
   - Enable required reviewers for production

3. **Add deployment credentials** as GitHub Secrets:
   - SSH keys for VPS
   - Cloud provider credentials (AWS, GCP, Azure)
   - Kubernetes config
   - Container registry tokens (if not using GHCR)

4. **Customize deployment steps** in workflow files:
   - Replace placeholder echo commands with actual deployment commands
   - Add your infrastructure-specific deployment logic
   - Configure health check endpoints
   - Set up monitoring and alerting

5. **Example deployment implementations**:

   **Docker Compose (VPS)**:
   ```yaml
   - name: Deploy via SSH
     uses: appleboy/ssh-action@v1.0.0
     with:
       host: ${{ secrets.DEPLOY_HOST }}
       username: ${{ secrets.DEPLOY_USER }}
       key: ${{ secrets.DEPLOY_SSH_KEY }}
       script: |
         cd /app
         export VERSION=${{ github.event.inputs.version }}
         docker compose -f docker-compose.production.yml pull
         docker compose -f docker-compose.production.yml up -d
   ```

   **Kubernetes**:
   ```yaml
   - name: Deploy to Kubernetes
     run: |
       kubectl set image deployment/backend \
         backend=${{ env.REGISTRY }}/${{ env.NAMESPACE }}/backend:${{ github.event.inputs.version }}
       kubectl rollout status deployment/backend
   ```

   **AWS ECS**:
   ```yaml
   - name: Deploy to ECS
     run: |
       aws ecs update-service \
         --cluster production-cluster \
         --service backend \
         --force-new-deployment
   ```

### Monorepo CI/CD Best Practices

✅ **Implemented in this repository**:

1. **Path-based caching**: Separate caches for Go, Python, and Node.js dependencies
2. **Parallel execution**: All services tested concurrently
3. **Fail fast**: Early validation prevents wasted CI time
4. **Selective testing**: Only affected services are tested (when configured)
5. **Build caching**: Docker layer caching for faster builds
6. **Artifact retention**: Test reports and logs preserved
7. **Coverage tracking**: Per-service coverage with Codecov flags
8. **Security scanning**: Vulnerability checks on dependencies
9. **Semantic versioning**: Enforced for production deployments
10. **Environment isolation**: Separate staging and production environments

### Monitoring and Rollback

**Health Checks**:
All deployments include health checks for:
- Backend API (`/health` endpoint)
- AI Service (`/health` endpoint)
- Frontend (homepage load)
- Database connectivity
- Redis cache connectivity

**Rollback Procedure**:
If a deployment fails or issues are detected:

1. **Automatic Rollback** (on health check failure):
   - Deployment workflow fails
   - Previous version remains active
   
2. **Manual Rollback**:
   ```bash
   # Identify the last known good version
   # Re-run deployment workflow with that version
   
   # Or use your infrastructure's rollback:
   kubectl rollout undo deployment/backend
   docker service rollback backend
   ```

### CI/CD Troubleshooting

**CI Tests Failing**:
- Check the specific job logs in GitHub Actions
- Run tests locally: `make test`
- Verify dependencies are up to date

**Build Failing**:
- Check Docker build logs
- Verify Dockerfiles are correct
- Test builds locally: `make prod-build`

**Deployment Failing**:
- Verify environment variables are configured
- Check deployment credentials/secrets
- Ensure target infrastructure is accessible
- Review health check endpoints

**Health Checks Failing**:
- Verify service URLs are correct
- Check if services started successfully
- Review application logs
- Verify database/Redis connectivity

## 🚢 Production Deployment

### Quick Deploy

The platform includes production-ready Docker images and deployment configurations for various platforms.

#### Build Production Images

```bash
# Build all production images
./scripts/push-images.sh --version 1.0.0 \
  --api-url https://api.yourdomain.com \
  --ai-url https://ai.yourdomain.com

# Or use environment variables
export VERSION=1.0.0
export REGISTRY=ghcr.io
export NAMESPACE=your-username/project-name
./scripts/push-images.sh
```

#### Deploy Options

**1. Docker Compose (Recommended for MVP)**
```bash
# Copy and configure production environment
cp .env.production.example .env.production
# Edit .env.production with your configuration

# Deploy with docker-compose
docker-compose -f docker-compose.production.yml --env-file .env.production up -d
```

**2. Platform-Specific Deployments**

- **AWS ECS**: See [DEPLOYMENT.md - AWS ECS Section](./DEPLOYMENT.md#option-1-aws-ecs)
- **Fly.io**: See [DEPLOYMENT.md - Fly.io Section](./DEPLOYMENT.md#option-2-flyio)
- **Railway**: See [DEPLOYMENT.md - Railway Section](./DEPLOYMENT.md#option-3-railway)
- **Self-Hosted**: See [DEPLOYMENT.md - Self-Hosted Section](./DEPLOYMENT.md#option-4-self-hosted-with-docker-compose)

### Production Files

- **Production Dockerfiles**:
  - `backend/Dockerfile.production` - Multi-stage optimized Go build
  - `ai_service/Dockerfile.production` - Multi-stage optimized Python build
  - `app/Dockerfile.production` - Static web build with Nginx
  
- **Configuration**:
  - `docker-compose.production.yml` - Production Docker Compose setup
  - `.env.production.example` - Production environment template
  
- **Documentation**:
  - [DEPLOYMENT.md](./DEPLOYMENT.md) - Complete deployment guide
  - [E2E_TESTING.md](./E2E_TESTING.md) - End-to-end testing and validation

### Security Checklist

Before deploying to production:

- ✅ Change all default passwords
- ✅ Generate secure JWT_SECRET
- ✅ Configure HTTPS/TLS
- ✅ Set up database backups
- ✅ Enable error tracking (Sentry)
- ✅ Configure monitoring
- ✅ Configure rate limiting
- ✅ Configure security headers (HSTS, CSP)
- ✅ Set CORS allowed origins
- ✅ Review file upload limits
- ✅ Test E2E flow

See [DEPLOYMENT.md - Security Checklist](./DEPLOYMENT.md#security-checklist) for complete list.
See [M7_SECURITY_HARDENING.md](./M7_SECURITY_HARDENING.md) for security hardening details.

### E2E Validation

After deployment, validate the complete user flow:

```bash
# Run automated E2E tests
npx playwright test

# Or follow manual testing guide
# See E2E_TESTING.md for detailed scenarios
```

Test the complete flow:
1. ✅ User signup/login
2. ✅ Project creation
3. ✅ Blueprint upload
4. ✅ AI analysis
5. ✅ Bid generation
6. ✅ PDF download

See [E2E_TESTING.md](./E2E_TESTING.md) for detailed testing procedures.

---

## 📄 Bid Export & Download

The platform provides professional bid export capabilities in multiple formats:

### Supported Formats

- **PDF** - Professional bid proposals with company branding
  - Cover page with logo
  - Itemized cost breakdown
  - Trade breakdown summary
  - Inclusions/exclusions
  - Payment terms and warranty
  
- **CSV** - Data export for analysis and integration
  - Full bid details in comma-separated format
  - Compatible with all spreadsheet applications
  - Easy to import into other systems

- **Excel** - Spreadsheet-ready format
  - UTF-8 BOM encoded for Excel compatibility
  - Opens directly in Microsoft Excel
  - Maintains data integrity

### Features

✅ **Professional PDF Template**
- Company logo and branding support
- Optional cover page with company information
- Itemized trade breakdown
- Detailed inclusions and exclusions
- Payment terms and warranty information

✅ **Multi-Format Consistency**
- Same data across all export formats
- Format-optimized presentation
- Reliable data integrity

✅ **Easy Downloads**
- Simple API endpoints
- Frontend integration ready
- Direct file downloads

### API Endpoints

```bash
# Get PDF URL
GET /bids/{id}/pdf

# Download CSV
GET /bids/{id}/csv

# Download Excel
GET /bids/{id}/excel
```

### Usage Example

```typescript
import { bidsApi } from './api/bids';

// Download PDF
const { pdf_url } = await bidsApi.getBidPDF(bidId);
window.open(pdf_url, '_blank');

// Download CSV
const csvBlob = await bidsApi.downloadBidCSV(bidId);
// Create download link and trigger download

// Download Excel
const excelBlob = await bidsApi.downloadBidExcel(bidId);
// Create download link and trigger download
```

For detailed documentation, see [BID_EXPORT_GUIDE.md](./BID_EXPORT_GUIDE.md)

---

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details on:
- Code of conduct
- Development setup
- Coding standards
- Commit message conventions
- Pull request process

Quick start for contributors:
```bash
# Fork the repository and clone your fork
git clone https://github.com/YOUR-USERNAME/fantastic-octo-memory.git

# Install pre-commit hooks (recommended)
pip install pre-commit
pre-commit install

# Create a feature branch
git checkout -b feature/your-feature-name

# Make your changes and commit
git commit -m "feat: add amazing feature"

# Push and create a pull request
git push origin feature/your-feature-name
```

## 🔒 Security

Security is a top priority. Please see our [Security Policy](SECURITY.md) for:
- Reporting vulnerabilities
- Security best practices
- Supported versions

If you discover a security vulnerability, please report it responsibly through our [Security Advisory](https://github.com/wonbyte/fantastic-octo-memory/security) page.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

For issues and questions, please open a GitHub issue.

## 📚 Additional Documentation

- [CODE_QUALITY_ANALYSIS.md](./CODE_QUALITY_ANALYSIS.md) - Comprehensive code quality analysis and recommendations
- [CONTRIBUTING.md](./CONTRIBUTING.md) - Contribution guidelines
- [SECURITY.md](./SECURITY.md) - Security policy and vulnerability reporting
- [CHANGELOG.md](./CHANGELOG.md) - Version history and release notes
- [BID_EXPORT_GUIDE.md](./BID_EXPORT_GUIDE.md) - Comprehensive bid export and download guide
- [DEPLOYMENT.md](./DEPLOYMENT.md) - Production deployment guide
- [E2E_TESTING.md](./E2E_TESTING.md) - End-to-end testing guide
- [M5_IMPLEMENTATION_SUMMARY.md](./M5_IMPLEMENTATION_SUMMARY.md) - Milestone 5 summary
- [M6_IMPLEMENTATION_SUMMARY.md](./M6_IMPLEMENTATION_SUMMARY.md) - Milestone 6 summary
- [M7_IMPLEMENTATION_SUMMARY.md](./M7_IMPLEMENTATION_SUMMARY.md) - Milestone 7 summary
- [M7_SECURITY_HARDENING.md](./M7_SECURITY_HARDENING.md) - Security hardening details
- [backend/README.md](./backend/README.md) - Backend service documentation
- [ai_service/README.md](./ai_service/README.md) - AI service documentation (if exists)
- [app/README.md](./app/README.md) - Frontend application documentation