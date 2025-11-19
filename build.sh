#!/bin/bash
set -e

# ---------------------------
# Config
# ---------------------------
DOCKERHUB_USER="jemiezler"

BACKEND_DIR="./backend"
FRONTEND_ADMIN_DIR="./frontend-admin"
FRONTEND_DIR="./frontend"

BACKEND_IMAGE="freshease-backend"
FRONTEND_ADMIN_IMAGE="freshease-frontend-admin"
FRONTEND_IMAGE="freshease-frontend"

PLATFORMS="linux/amd64,linux/arm64"

# GitHub Release config
GITHUB_REPO="jemiezler/freshease"         # Replace with your repo
RELEASE_TAG="v1.0.0"             # Replace with your release tag
APK_FILE="./frontend/build/app/outputs/flutter-apk/app-release.apk"  # Path to APK

# ---------------------------
# Build & Push Function
# ---------------------------
build_and_push() {
  local SERVICE_DIR=$1
  local IMAGE_NAME=$2

  echo "=== Building $IMAGE_NAME ==="

  docker buildx build \
    --platform $PLATFORMS \
    -t docker.io/${DOCKERHUB_USER}/${IMAGE_NAME}:latest \
    --push \
    $SERVICE_DIR

  echo "=== $IMAGE_NAME Done ==="
  echo ""
}

# ---------------------------
# Parse Arguments
# ---------------------------
BUILD_BACKEND=false
BUILD_FRONTEND=false
BUILD_FRONTEND_ADMIN=false
DEPLOY_APK=false

if [[ $# -eq 0 ]]; then
  BUILD_BACKEND=true
  BUILD_FRONTEND=true
  BUILD_FRONTEND_ADMIN=true
  DEPLOY_APK=true
else
  for arg in "$@"; do
    case $arg in
      backend) BUILD_BACKEND=true ;;
      frontend) BUILD_FRONTEND=true ;;
      frontend-admin) BUILD_FRONTEND_ADMIN=true ;;
      apk) DEPLOY_APK=true ;;
      all) BUILD_BACKEND=true; BUILD_FRONTEND=true; BUILD_FRONTEND_ADMIN=true; DEPLOY_APK=true ;;
      *) echo "Unknown argument: $arg"; echo "Usage: $0 [backend|frontend|frontend-admin|apk|all]"; exit 1 ;;
    esac
  done
fi

# ---------------------------
# Login to Docker Hub
# ---------------------------
echo "Checking Docker credentials..."

if docker info | grep -q "Username:"; then
    echo "Already logged in to Docker Hub."
else
    echo "Logging in to Docker Hub..."
    if [ -z "$DOCKER_USERNAME" ] || [ -z "$DOCKER_PASSWORD" ]; then
        echo "Error: DOCKER_USERNAME or DOCKER_PASSWORD not set."
        exit 1
    fi
    echo "$DOCKER_PASSWORD" | docker login docker.io -u "$DOCKER_USERNAME" --password-stdin
fi

# ---------------------------
# Build Selected Services
# ---------------------------
if $BUILD_BACKEND; then
  build_and_push "$BACKEND_DIR" "$BACKEND_IMAGE"
fi

if $BUILD_FRONTEND; then
  build_and_push "$FRONTEND_DIR" "$FRONTEND_IMAGE"
fi

if $BUILD_FRONTEND_ADMIN; then
  build_and_push "$FRONTEND_ADMIN_DIR" "$FRONTEND_ADMIN_IMAGE"
fi

# ---------------------------
# Deploy APK to GitHub Release
# ---------------------------
if $DEPLOY_APK; then
  echo "=== Deploying APK to GitHub Release ==="

  if [ -z "$GITHUB_TOKEN" ]; then
    echo "Error: GITHUB_TOKEN is not set."
    exit 1
  fi

  if [ ! -f "$APK_FILE" ]; then
    echo "Error: APK file not found at $APK_FILE"
    exit 1
  fi

  # Create or get release ID
  RELEASE_ID=$(curl -s -H "Authorization: token $GITHUB_TOKEN" \
      "https://api.github.com/repos/$GITHUB_REPO/releases/tags/$RELEASE_TAG" | jq -r '.id')

  if [ "$RELEASE_ID" == "null" ]; then
    echo "Release $RELEASE_TAG not found, creating..."
    RELEASE_ID=$(curl -s -X POST -H "Authorization: token $GITHUB_TOKEN" \
      -d "{\"tag_name\":\"$RELEASE_TAG\",\"name\":\"$RELEASE_TAG\",\"body\":\"Automated release\"}" \
      "https://api.github.com/repos/$GITHUB_REPO/releases" | jq -r '.id')
  fi

  echo "Uploading APK..."
  curl -s -X POST \
    -H "Authorization: token $GITHUB_TOKEN" \
    -H "Content-Type: application/vnd.android.package-archive" \
    --data-binary @"$APK_FILE" \
    "https://uploads.github.com/repos/$GITHUB_REPO/releases/$RELEASE_ID/assets?name=$(basename $APK_FILE)"

  echo "APK deployed successfully!"
fi

echo "All done!"
