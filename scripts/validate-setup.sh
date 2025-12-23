#!/bin/bash
# Setup validation script for Construction Estimation Platform
# Checks that all required tools and dependencies are installed

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "========================================"
echo "Setup Validation"
echo "========================================"
echo ""

# Track validation status
ERRORS=0
WARNINGS=0

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to check version
check_version() {
    local tool=$1
    local required=$2
    local current=$3
    
    echo -n "  Version check: "
    if [ -n "$current" ]; then
        echo -e "${GREEN}$current${NC}"
    else
        echo -e "${YELLOW}Unable to determine version${NC}"
        ((WARNINGS++))
    fi
}

# Check Docker
echo -n "Checking Docker... "
if command_exists docker; then
    echo -e "${GREEN}✓${NC}"
    DOCKER_VERSION=$(docker --version | grep -oP '\d+\.\d+\.\d+' | head -1)
    check_version "Docker" "20.10+" "$DOCKER_VERSION"
else
    echo -e "${RED}✗ Not found${NC}"
    echo "  Docker is required. Install from: https://docs.docker.com/get-docker/"
    ((ERRORS++))
fi

# Check Docker Compose
echo -n "Checking Docker Compose... "
if docker compose version >/dev/null 2>&1; then
    echo -e "${GREEN}✓${NC}"
    COMPOSE_VERSION=$(docker compose version | grep -oP '\d+\.\d+\.\d+' | head -1)
    check_version "Docker Compose" "2.0+" "$COMPOSE_VERSION"
elif command_exists docker-compose; then
    echo -e "${YELLOW}⚠ Legacy docker-compose found${NC}"
    echo "  Consider upgrading to Docker Compose V2"
    ((WARNINGS++))
else
    echo -e "${RED}✗ Not found${NC}"
    echo "  Docker Compose V2 is required"
    ((ERRORS++))
fi

# Check Make
echo -n "Checking Make... "
if command_exists make; then
    echo -e "${GREEN}✓${NC}"
else
    echo -e "${YELLOW}⚠ Not found${NC}"
    echo "  Make is recommended for convenience commands"
    ((WARNINGS++))
fi

# Check Node.js (for frontend development)
echo -n "Checking Node.js... "
if command_exists node; then
    echo -e "${GREEN}✓${NC}"
    NODE_VERSION=$(node --version | grep -oP '\d+\.\d+\.\d+')
    check_version "Node.js" "22.x LTS" "$NODE_VERSION"
    
    # Check if version is at least 22
    NODE_MAJOR=$(echo "$NODE_VERSION" | cut -d. -f1)
    if [ "$NODE_MAJOR" -lt 22 ]; then
        echo -e "  ${YELLOW}⚠ Node.js 22 LTS recommended${NC}"
        ((WARNINGS++))
    fi
else
    echo -e "${YELLOW}⚠ Not found${NC}"
    echo "  Node.js 22 LTS is required for frontend development"
    echo "  Install from: https://nodejs.org/"
    ((WARNINGS++))
fi

# Check npm
if command_exists npm; then
    echo -n "Checking npm... "
    echo -e "${GREEN}✓${NC}"
    NPM_VERSION=$(npm --version)
    check_version "npm" "10.x" "$NPM_VERSION"
fi

# Check Go (for backend development)
echo -n "Checking Go... "
if command_exists go; then
    echo -e "${GREEN}✓${NC}"
    GO_VERSION=$(go version | grep -oP 'go\d+\.\d+(\.\d+)?' | sed 's/go//')
    check_version "Go" "1.25+" "$GO_VERSION"
else
    echo -e "${YELLOW}⚠ Not found${NC}"
    echo "  Go 1.25+ is required for backend development"
    echo "  Install from: https://golang.org/dl/"
    ((WARNINGS++))
fi

# Check Python (for AI service development)
echo -n "Checking Python... "
if command_exists python3; then
    echo -e "${GREEN}✓${NC}"
    PYTHON_VERSION=$(python3 --version | grep -oP '\d+\.\d+\.\d+')
    check_version "Python" "3.12+" "$PYTHON_VERSION"
    
    # Check if version is at least 3.12
    PYTHON_MINOR=$(echo "$PYTHON_VERSION" | cut -d. -f2)
    if [ "$PYTHON_MINOR" -lt 12 ]; then
        echo -e "  ${YELLOW}⚠ Python 3.12+ recommended${NC}"
        ((WARNINGS++))
    fi
else
    echo -e "${YELLOW}⚠ Not found${NC}"
    echo "  Python 3.12+ is required for AI service development"
    echo "  Install from: https://www.python.org/downloads/"
    ((WARNINGS++))
fi

# Check pip
if command_exists pip3; then
    echo -n "Checking pip... "
    echo -e "${GREEN}✓${NC}"
    PIP_VERSION=$(pip3 --version | grep -oP '\d+\.\d+\.\d+' | head -1)
    check_version "pip" "23.0+" "$PIP_VERSION"
fi

# Check Git
echo -n "Checking Git... "
if command_exists git; then
    echo -e "${GREEN}✓${NC}"
    GIT_VERSION=$(git --version | grep -oP '\d+\.\d+\.\d+')
    check_version "Git" "2.30+" "$GIT_VERSION"
else
    echo -e "${RED}✗ Not found${NC}"
    echo "  Git is required"
    ((ERRORS++))
fi

# Check environment files
echo ""
echo "Checking environment files..."
if [ -f .env ]; then
    echo -e "  .env: ${GREEN}✓${NC}"
else
    echo -e "  .env: ${YELLOW}⚠ Not found${NC}"
    echo "    Run: cp .env.example .env"
    ((WARNINGS++))
fi

if [ -f backend/.env ]; then
    echo -e "  backend/.env: ${GREEN}✓${NC}"
else
    echo -e "  backend/.env: ${YELLOW}⚠ Not found (optional)${NC}"
fi

if [ -f ai_service/.env ]; then
    echo -e "  ai_service/.env: ${GREEN}✓${NC}"
else
    echo -e "  ai_service/.env: ${YELLOW}⚠ Not found (optional)${NC}"
fi

if [ -f app/.env ]; then
    echo -e "  app/.env: ${GREEN}✓${NC}"
else
    echo -e "  app/.env: ${YELLOW}⚠ Not found (optional)${NC}"
fi

# Optional tools
echo ""
echo "Optional Development Tools:"

echo -n "  pre-commit: "
if command_exists pre-commit; then
    echo -e "${GREEN}✓ Installed${NC}"
else
    echo -e "${YELLOW}Not installed${NC}"
    echo "    Install with: pip install pre-commit"
    echo "    Then run: pre-commit install"
fi

echo -n "  golangci-lint: "
if command_exists golangci-lint; then
    echo -e "${GREEN}✓ Installed${NC}"
else
    echo -e "${YELLOW}Not installed${NC}"
    echo "    Install from: https://golangci-lint.run/usage/install/"
fi

echo -n "  ruff (Python linter): "
if command_exists ruff; then
    echo -e "${GREEN}✓ Installed${NC}"
else
    echo -e "${YELLOW}Not installed${NC}"
    echo "    Install with: pip install ruff"
fi

# Summary
echo ""
echo "========================================"
echo "Summary"
echo "========================================"

if [ $ERRORS -eq 0 ] && [ $WARNINGS -eq 0 ]; then
    echo -e "${GREEN}✓ All checks passed!${NC}"
    echo ""
    echo "You're ready to start developing!"
    echo "Run 'make dev' to start all services."
    exit 0
elif [ $ERRORS -eq 0 ]; then
    echo -e "${YELLOW}⚠ $WARNINGS warning(s) found${NC}"
    echo ""
    echo "Your setup is functional, but some optional tools are missing."
    echo "Run 'make dev' to start all services."
    exit 0
else
    echo -e "${RED}✗ $ERRORS error(s) and $WARNINGS warning(s) found${NC}"
    echo ""
    echo "Please install the missing required tools before proceeding."
    exit 1
fi
