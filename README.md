# Construction Estimation & Bidding Automation SaaS Platform

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

1. **Clone the repository**
   ```bash
   git clone https://github.com/wonbyte/fantastic-octo-memory.git
   cd fantastic-octo-memory
   ```

2. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Start all services**
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

1. Create a feature branch from `main`
2. Make your changes
3. Ensure all tests pass
4. Submit a pull request

## 📄 License

[Your License Here]

## 🆘 Support

For issues and questions, please open a GitHub issue.

## 📚 Additional Documentation

- [BID_EXPORT_GUIDE.md](./BID_EXPORT_GUIDE.md) - Comprehensive bid export and download guide
- [DEPLOYMENT.md](./DEPLOYMENT.md) - Production deployment guide
- [E2E_TESTING.md](./E2E_TESTING.md) - End-to-end testing guide
- [M5_IMPLEMENTATION_SUMMARY.md](./M5_IMPLEMENTATION_SUMMARY.md) - Milestone 5 summary
- [M6_IMPLEMENTATION_SUMMARY.md](./M6_IMPLEMENTATION_SUMMARY.md) - Milestone 6 summary
- [M7_IMPLEMENTATION_SUMMARY.md](./M7_IMPLEMENTATION_SUMMARY.md) - Milestone 7 summary
- [backend/README.md](./backend/README.md) - Backend service documentation
- [ai_service/README.md](./ai_service/README.md) - AI service documentation (if exists)
- [app/README.md](./app/README.md) - Frontend application documentation