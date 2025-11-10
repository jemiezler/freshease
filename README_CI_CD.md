# CI/CD Quick Start Guide

## Overview

This project uses GitHub Actions for building Docker images and Jenkins for deployment.

## Quick Setup

### 1. GitHub Actions Setup

1. **Add GitHub Secrets**:
   - Go to repository Settings → Secrets and variables → Actions
   - Add the following secrets:
     - `GCP_SA_KEY`: Google Cloud Service Account JSON key
     - `JENKINS_WEBHOOK_URL`: (Optional) Jenkins webhook URL
     - `JENKINS_TOKEN`: (Optional) Jenkins authentication token

2. **Workflow Triggers**:
   - Automatically builds on push to `main` or `develop` branches
   - Can be manually triggered via GitHub Actions UI
   - Builds both API and Admin images

### 2. Jenkins Setup

1. **Install Required Tools** (run on Jenkins server):
   ```bash
   sudo ./scripts/setup-jenkins.sh
   ```

2. **Install Jenkins Plugins**:
   - Docker Pipeline
   - Docker
   - CloudBees Docker Build and Publish
   - Credentials Binding

3. **Configure Credentials**:
   - Add GCP Service Account key as credential ID: `gcp-service-account-key`

4. **Create Pipeline Job**:
   - New Item → Pipeline
   - Pipeline script from SCM
   - Git repository URL
   - Script Path: `Jenkinsfile`

### 3. Deployment

#### Automated Deployment (via Jenkins)
- Push to `main` branch triggers GitHub Actions
- GitHub Actions builds and pushes images
- Jenkins automatically deploys (if webhook configured)
- Or manually trigger Jenkins pipeline

#### Manual Deployment
```bash
# Deploy to production
./scripts/deploy.sh production

# Rollback to previous version
./scripts/rollback.sh previous
```

## File Structure

```
.
├── .github/
│   └── workflows/
│       └── build-and-push.yml    # GitHub Actions workflow
├── scripts/
│   ├── deploy.sh                 # Deployment script
│   ├── rollback.sh               # Rollback script
│   └── setup-jenkins.sh          # Jenkins setup script
├── Jenkinsfile                   # Jenkins pipeline
├── docker-compose.yml            # Production compose file
├── docker-compose.prod.yml       # Enhanced production compose file
└── CI_CD_SETUP.md               # Detailed setup documentation
```

## Environment Variables

Create `.env.prod` file with required environment variables (see `CI_CD_SETUP.md` for details).

## Health Checks

- API: `http://localhost:8000/api/health`
- Admin: `http://localhost:3000`

## Troubleshooting

See `CI_CD_SETUP.md` for detailed troubleshooting guide.

## Support

For issues, contact DevOps team or create an issue in the repository.

