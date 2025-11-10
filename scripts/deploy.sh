#!/bin/bash
set -e

# Deployment script for Freshease
# Usage: ./scripts/deploy.sh [environment]
# Environment: production, staging (default: production)

ENVIRONMENT="${1:-production}"
COMPOSE_FILE="docker-compose.yml"
REGISTRY="asia-southeast1-docker.pkg.dev"
PROJECT_ID="nodal-kite-477007-s4"
REPOSITORY="freshease"

echo "=========================================="
echo "Deploying Freshease to ${ENVIRONMENT}"
echo "=========================================="

# Authenticate to GCP if not already authenticated
if ! gcloud auth list --filter=status:ACTIVE --format="value(account)" | grep -q .; then
    echo "Authenticating to GCP..."
    gcloud auth activate-service-account --key-file="${GCP_SERVICE_ACCOUNT_KEY:-/path/to/service-account-key.json}"
fi

# Configure Docker to use gcloud as credential helper
echo "Configuring Docker authentication..."
gcloud auth configure-docker ${REGISTRY} --quiet

# Pull latest images
echo "Pulling latest images..."
docker-compose -f ${COMPOSE_FILE} pull

# Stop existing containers
echo "Stopping existing containers..."
docker-compose -f ${COMPOSE_FILE} down

# Start services
echo "Starting services..."
docker-compose -f ${COMPOSE_FILE} up -d

# Wait for services to be healthy
echo "Waiting for services to be healthy..."
sleep 10

# Health check
echo "Performing health checks..."
MAX_RETRIES=30
RETRY_COUNT=0

# Check API health
echo "Checking API health..."
while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    if curl -f http://localhost:8000/api/health > /dev/null 2>&1; then
        echo "✓ API is healthy"
        break
    fi
    RETRY_COUNT=$((RETRY_COUNT + 1))
    echo "Waiting for API... ($RETRY_COUNT/$MAX_RETRIES)"
    sleep 2
done

if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
    echo "✗ API health check failed"
    exit 1
fi

# Check Admin health
RETRY_COUNT=0
echo "Checking Admin health..."
while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    if curl -f http://localhost:3000 > /dev/null 2>&1; then
        echo "✓ Admin is healthy"
        break
    fi
    RETRY_COUNT=$((RETRY_COUNT + 1))
    echo "Waiting for Admin... ($RETRY_COUNT/$MAX_RETRIES)"
    sleep 2
done

if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
    echo "✗ Admin health check failed"
    exit 1
fi

# Show running containers
echo "=========================================="
echo "Deployment completed successfully!"
echo "=========================================="
docker-compose -f ${COMPOSE_FILE} ps

# Cleanup old images
echo "Cleaning up old images..."
docker image prune -f

echo "Deployment to ${ENVIRONMENT} completed!"

