# GitHub Actions Workflows

## build-and-push.yml

This workflow builds and pushes Docker images to Google Artifact Registry.

### Triggers

- **Push to main/develop**: Automatically builds when code is pushed to `main` or `develop` branches
- **Pull Requests**: Builds images for PRs (but doesn't push to registry)
- **Manual Dispatch**: Can be manually triggered via GitHub Actions UI

### What it does

1. **Builds API Image**: Builds Docker image from `backend/Dockerfile`
2. **Builds Admin Image**: Builds Docker image from `frontend-admin/Dockerfile`
3. **Pushes to Registry**: Pushes images to `asia-southeast1-docker.pkg.dev/nodal-kite-477007-s4/freshease/`
4. **Tags Images**: Tags with branch name, commit SHA, and `latest` (for main branch)

### Requirements

- GitHub Secret: `GCP_SA_KEY` (Google Cloud Service Account JSON key)
- GitHub Secret: `JENKINS_WEBHOOK_URL` (Optional, for Jenkins webhook)
- GitHub Secret: `JENKINS_TOKEN` (Optional, for Jenkins authentication)
- GitHub Secret: `NEXT_PUBLIC_API_BASE_URL` (Optional, for frontend build)

### Image Tags

- `latest`: Latest build from main branch
- `main-{sha}`: Specific commit from main branch
- `develop-{sha}`: Specific commit from develop branch
- `pr-{number}`: Pull request builds

### Usage

#### Automatic Build
Push to `main` or `develop` branch and the workflow will automatically run.

#### Manual Build
1. Go to Actions tab in GitHub
2. Select "Build and Push Docker Images"
3. Click "Run workflow"
4. Select service to build (api, admin, or both)
5. Click "Run workflow"

### Troubleshooting

- **Build fails**: Check logs in GitHub Actions
- **Push fails**: Verify GCP credentials are correct
- **Image not found**: Check if Artifact Registry repository exists
- **Authentication errors**: Verify `GCP_SA_KEY` secret is set correctly

