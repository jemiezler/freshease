# CI/CD Setup Guide

This document describes the CI/CD pipeline setup for Freshease using GitHub Actions and Jenkins.

## Overview

- **GitHub Actions**: Builds and pushes Docker images to Google Artifact Registry
- **Jenkins**: Deploys the application to production using Docker Compose

## Architecture

```
GitHub (Code) 
    ↓ (Push to main/develop)
GitHub Actions (Build & Push Images)
    ↓ (Images pushed to)
Google Artifact Registry
    ↓ (Notify)
Jenkins (Deploy)
    ↓ (Pull & Deploy)
Production Server (Docker Compose)
```

## Prerequisites

### 1. Google Cloud Platform Setup

1. Create a service account with the following roles:
   - Artifact Registry Writer
   - Storage Admin (if needed)

2. Create a JSON key for the service account

3. Store the key in GitHub Secrets as `GCP_SA_KEY`

### 2. GitHub Secrets

Add the following secrets to your GitHub repository:

- `GCP_SA_KEY`: Google Cloud Service Account JSON key
- `JENKINS_WEBHOOK_URL`: Jenkins webhook URL (optional)
- `JENKINS_TOKEN`: Jenkins authentication token (optional)
- `NEXT_PUBLIC_API_BASE_URL`: API base URL for frontend build (optional, defaults to production URL)

### 3. Jenkins Setup

1. Install required Jenkins plugins:
   - Docker Pipeline
   - Docker
   - CloudBees Docker Build and Publish
   - Google Cloud Storage
   - Credentials Binding

2. Configure Jenkins credentials:
   - Add GCP Service Account key as credential ID: `gcp-service-account-key`
   - Configure Docker credentials if needed

3. Install required tools on Jenkins server:
   - Docker
   - Docker Compose
   - Google Cloud SDK (gcloud)
   - curl (for health checks)

4. Configure Jenkins pipeline:
   - Create a new pipeline job
   - Point it to the repository
   - Select "Pipeline script from SCM"
   - Set SCM to Git
   - Set Script Path to `Jenkinsfile`

## GitHub Actions Workflow

### Trigger Events

- Push to `main` or `develop` branches
- Pull requests to `main` or `develop` branches
- Manual workflow dispatch

### Workflow Steps

1. **Build and Push API Image**
   - Builds Docker image from `backend/Dockerfile`
   - Pushes to `asia-southeast1-docker.pkg.dev/nodal-kite-477007-s4/freshease/api:latest`
   - Tags images with branch name, commit SHA, and semantic version

2. **Build and Push Admin Image**
   - Builds Docker image from `frontend-admin/Dockerfile`
   - Pushes to `asia-southeast1-docker.pkg.dev/nodal-kite-477007-s4/freshease/admin:latest`
   - Tags images with branch name, commit SHA, and semantic version

3. **Notify Jenkins** (optional)
   - Sends webhook to Jenkins to trigger deployment
   - Only runs on successful builds (not on PRs)

## Jenkins Pipeline

### Pipeline Stages

1. **Checkout**: Clones the repository
2. **Authenticate to GCP**: Authenticates using service account
3. **Pull Latest Images**: Pulls latest images from Artifact Registry
4. **Deploy to Production**: Deploys using Docker Compose
5. **Health Check**: Verifies services are healthy
6. **Cleanup**: Removes unused images

### Deployment Process

1. Stop existing containers
2. Pull latest images from registry
3. Start services with Docker Compose
4. Wait for services to be healthy
5. Perform health checks
6. Clean up old images

## Manual Deployment

### Using Deployment Script

```bash
# Deploy to production
./scripts/deploy.sh production

# Deploy to staging
./scripts/deploy.sh staging
```

### Using Docker Compose Directly

```bash
# Pull latest images
docker-compose pull

# Stop existing containers
docker-compose down

# Start services
docker-compose up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f
```

## Rollback

### Using Rollback Script

```bash
# Rollback to previous version
./scripts/rollback.sh previous

# Rollback to specific tag
./scripts/rollback.sh main-abc1234
```

### Manual Rollback

1. Update `docker-compose.yml` with previous image tag
2. Run deployment script
3. Or manually pull and deploy specific image tag

## Environment Variables

### Production Environment

Create `.env.prod` file with the following variables:

```env
# Database
DATABASE_URL=postgres://user:password@host:5432/dbname?sslmode=require

# JWT
JWT_SECRET=your-jwt-secret

# OIDC Google
OIDC_GOOGLE_ISSUER=https://accounts.google.com
OIDC_GOOGLE_CLIENT_ID=your-client-id
OIDC_GOOGLE_CLIENT_SECRET=your-client-secret
OIDC_GOOGLE_REDIRECT_URI=https://your-domain.com/api/auth/oidc/google/callback

# GenAI
GENAI_APIKEY=your-genai-api-key

# MinIO
MINIO_ENDPOINT=minio.example.com:9000
MINIO_ACCESS_KEY_ID=your-access-key
MINIO_SECRET_ACCESS_KEY=your-secret-key
MINIO_BUCKET=freshease
MINIO_USE_SSL=true

# HTTP
HTTP_PORT=:8080

# Ent Debug
ENT_DEBUG=false
```

## Health Checks

### API Health Endpoint

The API should have a health endpoint at `/api/health`. If it doesn't exist, add one or modify the health check scripts.

### Admin Health Check

The admin frontend is checked by making a request to `http://localhost:3000`.

## Troubleshooting

### Images Not Building

1. Check GitHub Actions logs
2. Verify GCP credentials are correct
3. Ensure Artifact Registry repository exists
4. Check Dockerfile syntax

### Deployment Fails

1. Check Jenkins logs
2. Verify Docker and Docker Compose are installed
3. Check GCP authentication
4. Verify environment variables are set
5. Check health endpoints are accessible

### Services Not Starting

1. Check Docker Compose logs: `docker-compose logs`
2. Verify environment variables
3. Check port conflicts
4. Verify images are pulled successfully

### Health Checks Failing

1. Check if services are running: `docker-compose ps`
2. Check service logs: `docker-compose logs [service-name]`
3. Verify health endpoints are accessible
4. Check network configuration

## Monitoring

### View Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f api
docker-compose logs -f admin
```

### Check Status

```bash
# Container status
docker-compose ps

# Resource usage
docker stats

# Image information
docker images | grep freshease
```

## Best Practices

1. **Always test in staging first**: Test deployments in a staging environment before production
2. **Use semantic versioning**: Tag images with semantic versions for better rollback capabilities
3. **Monitor deployments**: Set up monitoring and alerting for deployments
4. **Backup before deployment**: Always backup data before major deployments
5. **Use health checks**: Always verify services are healthy after deployment
6. **Keep secrets secure**: Never commit secrets to repository
7. **Use environment-specific configs**: Use different configs for different environments

## Security Considerations

1. **Service Account Permissions**: Use least privilege principle for service accounts
2. **Secrets Management**: Store secrets in secure vaults (GitHub Secrets, Jenkins Credentials)
3. **Image Scanning**: Regularly scan Docker images for vulnerabilities
4. **Network Security**: Use private networks and firewalls
5. **Authentication**: Use strong authentication for all services
6. **Encryption**: Use TLS/SSL for all external communications

## Support

For issues or questions, please contact the DevOps team or create an issue in the repository.

