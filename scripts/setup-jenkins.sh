#!/bin/bash
set -e

# Setup script for Jenkins server
# This script should be run on the Jenkins server to install required dependencies

echo "=========================================="
echo "Setting up Jenkins for Freshease CI/CD"
echo "=========================================="

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    echo "Please run as root or with sudo"
    exit 1
fi

# Update system
echo "Updating system packages..."
apt-get update -y

# Install Docker
if ! command -v docker &> /dev/null; then
    echo "Installing Docker..."
    apt-get install -y \
        ca-certificates \
        curl \
        gnupg \
        lsb-release
    
    mkdir -p /etc/apt/keyrings
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
    
    echo \
      "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
      $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null
    
    apt-get update -y
    apt-get install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin
    
    # Add jenkins user to docker group
    usermod -aG docker jenkins
else
    echo "Docker is already installed"
fi

# Install Docker Compose (standalone)
if ! command -v docker-compose &> /dev/null; then
    echo "Installing Docker Compose..."
    curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    chmod +x /usr/local/bin/docker-compose
else
    echo "Docker Compose is already installed"
fi

# Install Google Cloud SDK
if ! command -v gcloud &> /dev/null; then
    echo "Installing Google Cloud SDK..."
    echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] https://packages.cloud.google.com/apt cloud-sdk main" | tee -a /etc/apt/sources.list.d/google-cloud-sdk.list
    curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key --keyring /usr/share/keyrings/cloud.google.gpg add -
    apt-get update -y
    apt-get install -y google-cloud-sdk
else
    echo "Google Cloud SDK is already installed"
fi

# Install curl (for health checks)
if ! command -v curl &> /dev/null; then
    echo "Installing curl..."
    apt-get install -y curl
else
    echo "curl is already installed"
fi

# Install wget (for health checks in containers)
if ! command -v wget &> /dev/null; then
    echo "Installing wget..."
    apt-get install -y wget
else
    echo "wget is already installed"
fi

# Restart Jenkins to apply changes
echo "Restarting Jenkins..."
systemctl restart jenkins || service jenkins restart

echo "=========================================="
echo "Setup completed successfully!"
echo "=========================================="
echo ""
echo "Next steps:"
echo "1. Install required Jenkins plugins:"
echo "   - Docker Pipeline"
echo "   - Docker"
echo "   - CloudBees Docker Build and Publish"
echo "   - Credentials Binding"
echo ""
echo "2. Configure Jenkins credentials:"
echo "   - Add GCP Service Account key as credential ID: 'gcp-service-account-key'"
echo ""
echo "3. Create Jenkins pipeline job:"
echo "   - Create new pipeline job"
echo "   - Point to repository"
echo "   - Select 'Pipeline script from SCM'"
echo "   - Set Script Path to 'Jenkinsfile'"
echo ""
echo "4. Test the pipeline by running a build"

