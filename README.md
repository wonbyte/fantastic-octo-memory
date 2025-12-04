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
| **Frontend** | React | 18.3+ |
| **Mobile** | React Native | 0.76+ |
| **Mobile Framework** | Expo SDK | 52+ |
| **Runtime** | Node.js | 22 LTS |
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

## 🤝 Contributing

1. Create a feature branch from `main`
2. Make your changes
3. Ensure all tests pass
4. Submit a pull request

## 📄 License

[Your License Here]

## 🆘 Support

For issues and questions, please open a GitHub issue.