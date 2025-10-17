#!/bin/bash

# IU-MiCert Issuer Deployment Script for GCP VM
# This script handles the deployment of the issuer backend with Docker

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if .env file exists
if [ ! -f .env ]; then
    log_error ".env file not found!"
    log_info "Please copy .env.production to .env and configure it:"
    echo "  cp .env.production .env"
    echo "  nano .env  # Edit with your values"
    exit 1
fi

# Source environment variables
source .env

# Validate required environment variables
if [ -z "$ISSUER_PRIVATE_KEY" ] || [ "$ISSUER_PRIVATE_KEY" = "YOUR_SEPOLIA_PRIVATE_KEY_WITHOUT_0x_PREFIX" ]; then
    log_error "ISSUER_PRIVATE_KEY is not configured in .env"
    exit 1
fi

if [[ "$SEPOLIA_RPC_URL" == *"YOUR_INFURA_API_KEY"* ]]; then
    log_error "SEPOLIA_RPC_URL is not properly configured in .env"
    exit 1
fi

log_info "Starting IU-MiCert Issuer deployment..."

# Pull latest changes (if deploying from git)
if [ -d ".git" ]; then
    log_info "Pulling latest changes from git..."
    git pull
fi

# Stop existing containers
log_info "Stopping existing containers..."
docker compose down

# Remove old images to ensure fresh build
log_warn "Removing old images..."
docker compose rm -f

# Build and start containers
log_info "Building and starting containers..."
docker compose up -d --build

# Wait for services to be healthy
log_info "Waiting for services to start..."
sleep 10

# Check service health
log_info "Checking service health..."
if docker compose ps | grep -q "iumicert-postgres.*Up.*healthy"; then
    log_info "PostgreSQL is healthy"
else
    log_error "PostgreSQL failed to start properly"
    docker compose logs postgres
    exit 1
fi

if docker compose ps | grep -q "iumicert-issuer.*Up"; then
    log_info "Issuer backend is running"
else
    log_error "Issuer backend failed to start"
    docker compose logs issuer-backend
    exit 1
fi

# Show running containers
log_info "Running containers:"
docker compose ps

# Show logs
log_info "Recent logs:"
docker compose logs --tail=50

log_info "Deployment complete!"
log_info "API is available at: http://localhost:8080"
log_info "Health check: curl http://localhost:8080/api/health"
log_info ""
log_info "To view logs: docker compose logs -f"
log_info "To stop: docker compose down"
