#!/bin/bash

# Container Registry Push Script
# ================================
# Script to build and push Docker images to a container registry
# Supports: GitHub Container Registry (GHCR), Docker Hub, AWS ECR

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default values
REGISTRY="${REGISTRY:-ghcr.io}"
NAMESPACE="${NAMESPACE:-wonbyte/fantastic-octo-memory}"
VERSION="${VERSION:-latest}"
BUILD_ARGS=""

# Parse command line arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    --registry)
      REGISTRY="$2"
      shift 2
      ;;
    --namespace)
      NAMESPACE="$2"
      shift 2
      ;;
    --version)
      VERSION="$2"
      shift 2
      ;;
    --api-url)
      BUILD_ARGS="$BUILD_ARGS --build-arg EXPO_PUBLIC_API_URL=$2"
      shift 2
      ;;
    --ai-url)
      BUILD_ARGS="$BUILD_ARGS --build-arg EXPO_PUBLIC_AI_SERVICE_URL=$2"
      shift 2
      ;;
    --help)
      echo "Usage: $0 [OPTIONS]"
      echo ""
      echo "Options:"
      echo "  --registry REGISTRY        Container registry (default: ghcr.io)"
      echo "  --namespace NAMESPACE      Namespace/organization (default: wonbyte/fantastic-octo-memory)"
      echo "  --version VERSION          Image version tag (default: latest)"
      echo "  --api-url URL             Frontend API URL"
      echo "  --ai-url URL              Frontend AI Service URL"
      echo "  --help                    Show this help message"
      echo ""
      echo "Environment Variables:"
      echo "  REGISTRY                  Same as --registry"
      echo "  NAMESPACE                 Same as --namespace"
      echo "  VERSION                   Same as --version"
      echo ""
      echo "Examples:"
      echo "  $0 --version 1.0.0"
      echo "  $0 --registry docker.io --namespace mycompany/construction --version 1.2.3"
      echo "  $0 --registry ghcr.io --namespace username/repo --version latest --api-url https://api.example.com"
      exit 0
      ;;
    *)
      echo -e "${RED}Error: Unknown option $1${NC}"
      echo "Run '$0 --help' for usage information"
      exit 1
      ;;
  esac
done

# Display configuration
echo -e "${GREEN}=== Docker Build and Push Configuration ===${NC}"
echo "Registry:  $REGISTRY"
echo "Namespace: $NAMESPACE"
echo "Version:   $VERSION"
echo ""

# Verify Docker is running
if ! docker info > /dev/null 2>&1; then
  echo -e "${RED}Error: Docker is not running${NC}"
  exit 1
fi

# Build and push backend
echo -e "${GREEN}=== Building Backend Image ===${NC}"
docker build \
  -f backend/Dockerfile.production \
  -t $REGISTRY/$NAMESPACE/backend:$VERSION \
  -t $REGISTRY/$NAMESPACE/backend:latest \
  ./backend

echo -e "${GREEN}=== Pushing Backend Image ===${NC}"
docker push $REGISTRY/$NAMESPACE/backend:$VERSION
docker push $REGISTRY/$NAMESPACE/backend:latest
echo -e "${GREEN}✓ Backend pushed successfully${NC}"
echo ""

# Build and push AI service
echo -e "${GREEN}=== Building AI Service Image ===${NC}"
docker build \
  -f ai_service/Dockerfile.production \
  -t $REGISTRY/$NAMESPACE/ai-service:$VERSION \
  -t $REGISTRY/$NAMESPACE/ai-service:latest \
  ./ai_service

echo -e "${GREEN}=== Pushing AI Service Image ===${NC}"
docker push $REGISTRY/$NAMESPACE/ai-service:$VERSION
docker push $REGISTRY/$NAMESPACE/ai-service:latest
echo -e "${GREEN}✓ AI Service pushed successfully${NC}"
echo ""

# Build and push frontend
echo -e "${GREEN}=== Building Frontend Image ===${NC}"
docker build \
  -f app/Dockerfile.production \
  $BUILD_ARGS \
  -t $REGISTRY/$NAMESPACE/frontend:$VERSION \
  -t $REGISTRY/$NAMESPACE/frontend:latest \
  ./app

echo -e "${GREEN}=== Pushing Frontend Image ===${NC}"
docker push $REGISTRY/$NAMESPACE/frontend:$VERSION
docker push $REGISTRY/$NAMESPACE/frontend:latest
echo -e "${GREEN}✓ Frontend pushed successfully${NC}"
echo ""

# Summary
echo -e "${GREEN}=== Build and Push Complete ===${NC}"
echo ""
echo "Images pushed:"
echo "  Backend:  $REGISTRY/$NAMESPACE/backend:$VERSION"
echo "  AI Service: $REGISTRY/$NAMESPACE/ai-service:$VERSION"
echo "  Frontend: $REGISTRY/$NAMESPACE/frontend:$VERSION"
echo ""
echo "Next steps:"
echo "  1. Update your deployment configuration with the new image tags"
echo "  2. Deploy to your target environment"
echo "  3. Run E2E tests to validate deployment"
echo ""
echo -e "${GREEN}Done!${NC}"
