# Testing Summary

Comprehensive overview of the test infrastructure for the Construction Estimation & Bidding Automation platform.

## Overview

This document provides a complete guide to the testing capabilities implemented across all layers of the application stack.

## Test Categories

### 1. End-to-End (E2E) Tests

**Location**: `e2e/`
**Framework**: Playwright
**Language**: TypeScript

#### Test Files

- `basic.spec.ts` - Authentication, accessibility, dark mode, offline mode
- `user-journey.spec.ts` - Complete user flow, responsive design, performance
- `revision-comparison.spec.ts` - Revision history and comparison features
- `complete-workflow.spec.ts` - Comprehensive workflow from signup to PDF download

#### Coverage

- ✅ User authentication (signup, login, logout)
- ✅ Project creation and management
- ✅ Blueprint upload and analysis
- ✅ Bid generation and export
- ✅ PDF download
- ✅ Error handling and edge cases
- ✅ Mobile responsiveness (375px to 768px)
- ✅ Cross-browser (Chrome, Firefox, Safari)
- ✅ Accessibility (ARIA labels, keyboard navigation)
- ✅ Dark mode toggle
- ✅ Offline mode handling
- ✅ Revision comparison and history

#### Running E2E Tests

**Prerequisites**: E2E tests require the full application stack to be running:
- Frontend: `http://localhost:3000`
- Backend API: `http://localhost:8080`
- AI Service: `http://localhost:8000`
- Database (PostgreSQL)
- Redis
- S3/MinIO

```bash
# Start all services first
make dev  # or docker-compose up

# Install dependencies (first time)
npm install
npx playwright install --with-deps

# Run all tests
npm run test:e2e

# Run with UI (interactive)
npm run test:e2e:ui

# Run specific browser
npx playwright test --project=chromium

# Run specific file
npx playwright test e2e/complete-workflow.spec.ts

# Debug mode
npx playwright test --debug
```

**CI/CD Note**: E2E tests in CI are configured with `continue-on-error: true` because they require the full stack. They execute but don't block builds. For reliable E2E testing, use a staging environment or run locally with all services.

#### Test Statistics

- **Test Files**: 4
- **Unique Test Cases**: 22+
- **Total Test Cases** (across 5 browsers): 110+
- **Browser Projects**: 5 (Chrome, Firefox, Safari, Mobile Chrome, Mobile Safari)

### 2. Backend Integration Tests

**Location**: `backend/internal/integration/`
**Framework**: Go testing + testify
**Language**: Go

#### Test Files

- `api_integration_test.go` - Complete API workflow integration tests

#### Coverage

- ✅ Complete workflow: project → blueprint → analysis → bid
- ✅ Project creation with database
- ✅ Blueprint upload with S3
- ✅ Analysis workflow with job queue
- ✅ Bid generation with calculations
- ✅ Concurrent users simulation
- ✅ Database operations (connection pool, transactions)
- ✅ Redis cache operations
- ✅ S3 storage operations
- ✅ Authentication flow

#### Running Backend Integration Tests

```bash
cd backend

# Run all tests including integration
go test -v ./...

# Run only integration tests
go test -v ./internal/integration/...

# Run with coverage
go test -v -coverprofile=coverage.out ./internal/integration/...
go tool cover -html=coverage.out

# Skip integration tests (short mode)
go test -v -short ./...
```

#### Test Statistics

- **Test Files**: 1 new + 14 existing
- **Integration Test Functions**: 10+
- **Test Scenarios**: Database, Redis, S3, Auth, API workflows

### 3. AI Service Integration Tests

**Location**: `ai_service/`
**Framework**: pytest + FastAPI TestClient
**Language**: Python

#### Test Files

- `test_integration.py` - Blueprint analysis and bid generation integration
- `test_main.py` - Basic API endpoint tests
- `test_vision_enhancements.py` - Vision/OCR functionality tests

#### Coverage

- ✅ Blueprint analysis complete flow
- ✅ Blueprint analysis with OCR processing
- ✅ Invalid S3 key error handling
- ✅ Concurrent blueprint analysis
- ✅ Bid generation complete flow
- ✅ Bid generation with various markups
- ✅ Empty takeoff data handling
- ✅ Vision API health checks
- ✅ Performance and response time tests
- ✅ Memory usage stability tests
- ✅ Error handling (invalid JSON, missing fields)

#### Running AI Service Integration Tests

```bash
cd ai_service

# Install dependencies
pip install -r requirements.txt -r requirements-dev.txt

# Run all tests
pytest

# Run only integration tests
pytest -v -m integration

# Run with coverage
pytest --cov=. --cov-report=html

# Skip integration tests
pytest -m "not integration"

# Skip slow tests
pytest -m "not slow"
```

#### Test Statistics

- **Test Files**: 3
- **Integration Test Classes**: 5
- **Test Functions**: 20+
- **Markers**: `integration`, `slow`

### 4. Load/Performance Tests

**Location**: `load-tests/`
**Framework**: Artillery
**Format**: YAML

#### Configuration Files

- `artillery-backend.yml` - Backend API load tests
- `artillery-ai-service.yml` - AI service load tests
- `README.md` - Comprehensive load testing guide

#### Test Scenarios

**Backend API (5 scenarios):**
1. Health check monitoring
2. User authentication (signup)
3. Project CRUD operations
4. Blueprint workflow
5. Bid generation and export

**AI Service (5 scenarios):**
1. Health check
2. Root info endpoint
3. Blueprint analysis
4. Bid generation
5. Mixed operations

#### Load Test Phases

**Backend:**
- Warm up: 30s @ 5 RPS
- Ramp up: 60s @ 10→50 RPS
- Sustained: 120s @ 50 RPS
- Spike: 30s @ 100 RPS
- Cool down: 30s @ 10 RPS

**AI Service:**
- Warm up: 20s @ 2 RPS
- Ramp up: 60s @ 5→15 RPS
- Sustained: 90s @ 15 RPS
- Spike: 20s @ 30 RPS

#### Running Load Tests

```bash
# Install Artillery
npm install

# Test backend
npm run test:load:backend

# Test AI service
npm run test:load:ai

# Run all load tests
npm run test:load:all

# Custom target
artillery run --target https://api.yourdomain.com load-tests/artillery-backend.yml

# With report
artillery run --output report.json load-tests/artillery-backend.yml
artillery report report.json
```

#### Performance Targets

| Service | Metric | Target | Critical |
|---------|--------|--------|----------|
| **Backend** | P95 Response Time | < 500ms | < 1000ms |
| **Backend** | P99 Response Time | < 1000ms | < 2000ms |
| **Backend** | Error Rate | < 0.1% | < 1% |
| **Backend** | Throughput | > 100 RPS | > 50 RPS |
| **AI Service** | P95 Response Time | < 2s | < 5s |
| **AI Service** | P99 Response Time | < 5s | < 10s |
| **AI Service** | Error Rate | < 0.5% | < 2% |
| **AI Service** | Throughput | > 20 RPS | > 10 RPS |

## CI/CD Integration

### GitHub Actions Workflow

The CI workflow (`.github/workflows/ci.yml`) includes:

#### Jobs

1. **lint-and-test-backend** - Go linting and unit tests
2. **lint-and-test-ai-service** - Python linting and unit tests
3. **lint-and-test-app** - TypeScript linting and unit tests
4. **e2e-tests** - Playwright E2E tests (continue-on-error: requires full stack)
5. **integration-tests-backend** - Go integration tests
6. **integration-tests-ai-service** - Python integration tests
7. **docker-build** - Docker image build tests

**Note**: E2E tests are configured with `continue-on-error: true` because they require the full application stack (frontend, backend, AI service, database, Redis, S3) to be running. These tests are designed to run locally or in a full staging environment. They will execute in CI but won't block the build if they fail.

#### Coverage Reporting

All test jobs upload coverage to Codecov:
- Backend: `coverage.out`
- AI Service: `coverage.xml`
- App: `coverage-final.json`
- Integration tests: separate coverage files

#### Artifacts

- Playwright HTML reports
- Test screenshots (on failure)
- Load test reports

## Running All Tests

### Local Development

```bash
# Backend unit tests
cd backend && go test ./...

# Backend integration tests
cd backend && go test ./internal/integration/...

# AI service tests
cd ai_service && pytest

# AI service integration tests
cd ai_service && pytest -m integration

# App tests
cd app && npm test

# E2E tests
npm run test:e2e

# Load tests
npm run test:load:all
```

### CI/CD

Tests run automatically on:
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop` branches

## Test Fixtures

### E2E Test Fixtures

**Location**: `e2e/fixtures/`

For blueprint upload testing, add sample PDF files:
- `test-blueprint.pdf` - Small test blueprint (1-2 MB)
- `test-blueprint-large.pdf` - Large test file (8-10 MB)
- `test-blueprint-multipage.pdf` - Multi-page document

See `e2e/fixtures/README.md` for details on creating test fixtures.

## Best Practices

### Writing Tests

1. **E2E Tests**
   - Use descriptive test names
   - Add appropriate timeouts for async operations
   - Use `.catch()` for optional elements
   - Test critical user flows comprehensively

2. **Integration Tests**
   - Use `testing.Short()` to skip in CI if needed
   - Mock external dependencies when appropriate
   - Test error conditions and edge cases
   - Use proper cleanup in test teardown

3. **Load Tests**
   - Start with realistic user flows
   - Include think time between actions
   - Use unique identifiers for test data
   - Monitor system resources during tests

### Running Tests Efficiently

```bash
# Run specific test suites
go test -v -run TestSpecificFunction
pytest -k test_specific_function
npx playwright test -g "specific test name"

# Skip slow tests
go test -short ./...
pytest -m "not slow"

# Parallel execution
go test -parallel 4 ./...
pytest -n auto
npx playwright test --workers=4
```

## Troubleshooting

### Common Issues

**E2E Tests Timing Out**
- Increase timeout values
- Check if services are running
- Verify network connectivity
- Check browser console for errors

**Integration Tests Failing**
- Verify test database is available
- Check Redis connectivity
- Ensure S3/MinIO is configured
- Review service logs

**Load Tests High Error Rate**
- Check service health
- Verify resource limits (CPU, memory)
- Review rate limiting configuration
- Check database connection pool

## Documentation References

- [E2E_TESTING.md](./E2E_TESTING.md) - Detailed E2E testing guide
- [load-tests/README.md](./load-tests/README.md) - Load testing guide
- [e2e/fixtures/README.md](./e2e/fixtures/README.md) - Test fixtures guide
- [README.md](./README.md) - Main project documentation

## Metrics and Reporting

### Test Coverage

Current coverage targets:
- Backend: > 70%
- AI Service: > 70%
- App: > 60%

View coverage reports:
```bash
# Backend
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# AI Service
pytest --cov=. --cov-report=html
open htmlcov/index.html

# App
npm test -- --coverage
open coverage/lcov-report/index.html
```

### Performance Metrics

Monitor in production:
- Response times (P50, P95, P99)
- Error rates
- Throughput (RPS)
- Resource utilization (CPU, memory)

## Continuous Improvement

Regular testing activities:
- ✅ Run E2E tests before each deployment
- ✅ Run load tests after performance optimizations
- ✅ Update tests when adding new features
- ✅ Review and update performance targets quarterly
- ✅ Monitor test flakiness and fix unstable tests
- ✅ Keep test dependencies up to date

## Support

For testing questions or issues:
- Review this document and linked guides
- Check test logs and error messages
- Open an issue with test failure details
- Include environment information

---

**Last Updated**: 2024-12-16
**Version**: 1.0.0
