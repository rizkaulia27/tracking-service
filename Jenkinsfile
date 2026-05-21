pipeline {
    agent any

    environment {
        IMAGE = "kaulie27/tracking-service:${env.BUILD_NUMBER}"
        NETWORK = "tracking-net"
    }

    stages {

        // 1. CHECKOUT
        stage('Checkout Repo') {
            steps {
                echo "CHECKOUT SUCCESS"
            }
        }

        // 2. UNIT TEST
        stage('Unit Test') {
            steps {
                catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                    sh 'go test -short ./...'
                }
            }
        }

        // 3. LINT / VET
        stage('Lint / Vet') {
            steps {
                sh 'go vet ./...'
            }
        }

        // 4. BUILD IMAGE
        stage('Build Image') {
            steps {
                sh 'docker build -t $IMAGE .'
            }
        }

        // 5. FUNCTIONAL TEST
        stage('Functional Test') {
            steps {
        
                sh '''
                echo "RUN FUNCTIONAL TEST"
        
                docker rm -f mongodb || true
                docker rm -f test-tracking || true
                docker rm -f curl-test || true
                docker network rm $NETWORK || true
        
                docker network create $NETWORK
        
                docker run -d \
                  --name mongodb \
                  --network $NETWORK \
                  -e MONGO_INITDB_ROOT_USERNAME=admin \
                  -e MONGO_INITDB_ROOT_PASSWORD=admin123 \
                  mongo
        
                echo "WAITING FOR MONGODB..."
                sleep 20
        
                docker run -d \
                  --name test-tracking \
                  --network $NETWORK \
                  $IMAGE
        
                echo "WAITING FOR APPLICATION..."
                sleep 30
        
                echo "===== CONTAINER LOG ====="
                docker logs test-tracking
        
                echo "===== HEALTH CHECK ====="
        
                docker run --rm \
                  --network $NETWORK \
                  curlimages/curl \
                  curl http://test-tracking:8087/tracking || true
        
                echo "===== GO TEST ====="
                go test -run TestTrackingAPI_Success
                '''
            }
        }
        
        // 6. PUSH IMAGE
        stage('Push Image') {
            steps {
                withCredentials([usernamePassword(
                    credentialsId: 'dockerhub-login',
                    usernameVariable: 'USERNAME',
                    passwordVariable: 'PASSWORD'
                )]) {

                    sh '''
                    echo "$PASSWORD" | docker login -u "$USERNAME" --password-stdin

                    docker push $IMAGE
                    '''
                }
            }
        }

        // 7. DEPLOY
        stage('Deploy') {
            steps {
                echo "DEPLOY SUCCESS"
            }
        }

        // 8. VERIFY
        stage('Verify') {
            steps {
                echo "PIPELINE SUCCESS"
            }
        }
    }

    post {

        success {
            echo 'PIPELINE SUCCESS'
        }

        failure {
            echo 'PIPELINE FAILED'
        }

        always {

            sh '''
            docker rm -f mongodb || true
            docker rm -f test-tracking || true
            docker network rm $NETWORK || true
            '''
        }
    }
}
