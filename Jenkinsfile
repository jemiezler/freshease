pipeline {
    agent any

    environment {
        REGISTRY = 'asia-southeast1-docker.pkg.dev'
        PROJECT_ID = 'nodal-kite-477007-s4'
        REPOSITORY = 'freshease'
        IMAGE_API = "${REGISTRY}/${PROJECT_ID}/${REPOSITORY}/api:latest"
        IMAGE_ADMIN = "${REGISTRY}/${PROJECT_ID}/${REPOSITORY}/admin:latest"
        COMPOSE_FILE = 'docker-compose.yml'
    }

    options {
        buildDiscarder(logRotator(numToKeepStr: '10'))
        timeout(time: 30, unit: 'MINUTES')
        timestamps()
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
                script {
                    env.GIT_COMMIT = sh(
                        script: 'git rev-parse --short HEAD',
                        returnStdout: true
                    ).trim()
                    env.GIT_BRANCH = env.BRANCH_NAME ?: sh(
                        script: 'git rev-parse --abbrev-ref HEAD',
                        returnStdout: true
                    ).trim()
                }
            }
        }

        stage('Authenticate to GCP') {
            steps {
                script {
                    withCredentials([file(credentialsId: 'gcp-service-account-key', variable: 'GCP_KEY')]) {
                        sh '''
                            gcloud auth activate-service-account --key-file=${GCP_KEY}
                            gcloud auth configure-docker ${REGISTRY} --quiet
                        '''
                    }
                }
            }
        }

        stage('Pull Latest Images') {
            steps {
                script {
                    try {
                        sh '''
                            docker pull ${IMAGE_API} || echo "Failed to pull API image, will use existing"
                            docker pull ${IMAGE_ADMIN} || echo "Failed to pull Admin image, will use existing"
                        '''
                    } catch (Exception e) {
                        echo "Warning: Failed to pull images: ${e.message}"
                    }
                }
            }
        }

        stage('Deploy to Production') {
            when {
                anyOf {
                    branch 'main'
                    branch 'master'
                }
            }
            steps {
                script {
                    dir('.') {
                        // Stop existing containers
                        sh '''
                            docker-compose -f ${COMPOSE_FILE} down || true
                        '''

                        // Pull latest images
                        sh '''
                            docker-compose -f ${COMPOSE_FILE} pull
                        '''

                        // Start services
                        sh '''
                            docker-compose -f ${COMPOSE_FILE} up -d
                        '''

                        // Wait for services to be healthy
                        sh '''
                            echo "Waiting for services to be healthy..."
                            sleep 10
                            docker-compose -f ${COMPOSE_FILE} ps
                        '''
                    }
                }
            }
        }

        stage('Health Check') {
            steps {
                script {
                    sh '''
                        echo "Performing health checks..."
                        
                        # Check API health
                        API_HEALTHY=false
                        for i in {1..30}; do
                            if curl -f http://localhost:8000/api/health >/dev/null 2>&1; then
                                echo "✓ API is healthy"
                                API_HEALTHY=true
                                break
                            fi
                            echo "Waiting for API... ($i/30)"
                            sleep 2
                        done
                        
                        if [ "$API_HEALTHY" = false ]; then
                            echo "✗ API health check failed"
                            exit 1
                        fi
                        
                        # Check Admin health
                        ADMIN_HEALTHY=false
                        for i in {1..30}; do
                            if curl -f http://localhost:3000 >/dev/null 2>&1; then
                                echo "✓ Admin is healthy"
                                ADMIN_HEALTHY=true
                                break
                            fi
                            echo "Waiting for Admin... ($i/30)"
                            sleep 2
                        done
                        
                        if [ "$ADMIN_HEALTHY" = false ]; then
                            echo "✗ Admin health check failed"
                            exit 1
                        fi
                    '''
                }
            }
        }

        stage('Cleanup') {
            steps {
                script {
                    sh '''
                        # Remove unused images
                        docker image prune -f
                        
                        # Remove unused volumes (optional, be careful)
                        # docker volume prune -f
                    '''
                }
            }
        }
    }

    post {
        success {
            script {
                echo "Deployment successful!"
                // Optional: Send success notification
                // slackSend(color: 'good', message: "Deployment to production succeeded: ${env.GIT_COMMIT}")
            }
        }
        failure {
            script {
                echo "Deployment failed!"
                // Optional: Send failure notification
                // slackSend(color: 'danger', message: "Deployment to production failed: ${env.GIT_COMMIT}")
                
                // Rollback to previous version
                sh '''
                    docker-compose -f ${COMPOSE_FILE} down
                    # Add rollback logic here if needed
                '''
            }
        }
        always {
            cleanWs()
        }
    }
}

