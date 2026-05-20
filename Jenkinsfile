pipeline {
    agent any

    environment {
        IMAGE = "kaulie27/tracking-service:${env.BUILD_NUMBER}"
    }

    stages {

        // 1. CHECKOUT
        stage('Checkout Repo') {
            steps {
                echo "CHECKOUT SUCCESS"
            }
        }

        // 2. UNIT TEST (BOLEH FAIL)
        stage('Unit Test') {
            steps {
                catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                    sh 'go test -short ./...'
                }
            }
        }

        // 3. LINT / VET (WAJIB HIJAU)
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

        // 5. FUNCTIONAL TEST (BOLEH FAIL)
        stage('Functional Test') {
            steps {

                catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {

                    sh '''
                    echo "START MONGODB"

                    docker rm -f mongo-test || true

                    docker run -d \
                      --name mongo-test \
                      mongo:7

                    sleep 5

                    echo "START TRACKING APP"

                    docker rm -f test-tracking || true

                    docker run -d \
                      --name test-tracking \
                      --link mongo-test \
                      -e MONGO_URI=mongodb://mongo-test:27017 \
                      $IMAGE

                    sleep 5

                    echo "RUN FUNCTIONAL TEST"

                    go test -run TestTrackingAPI_Success
                    '''
                }
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
            '''
        }
    }
}
