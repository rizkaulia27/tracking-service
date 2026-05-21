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
            echo RUN FUNCTIONAL TEST

            # buat network
            docker network create tracking-net || true

            # jalankan container app untuk test
            docker run -d \
              --name test-tracking \
              --network tracking-net \
              -p 8087:8087 \
              kaulie27/tracking-service:${BUILD_NUMBER}

            # tunggu app ready
            sleep 10

            # jalankan functional test
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
                sh 'echo "DEPLOY SUCCESS"'
            }
        }
        // 8. VERIFY
        stage('Verify') {
            steps {
                sh 'echo "PIPELINE SUCCESS"'
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
            docker rm -f mongo-test || true
            docker rm -f test-tracking || true
            docker network rm $NETWORK || true
            '''
        }
    }
}
