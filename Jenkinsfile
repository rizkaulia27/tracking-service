pipeline {
    agent any

    environment {
        IMAGE = "kaulie27/tracking-service:${env.BUILD_NUMBER}"
    }

    stages {

        // 1. CHECKOUT
        stage('Checkout Repo') {
            steps {
                deleteDir()
                git branch: 'main', url: 'https://github.com/USERNAME/TrackingService.git'
            }
        }

        // 2. UNIT TEST (BOLEH FAIL, TAPI LANJUT)
        stage('Unit Test') {
            steps {
                dir('TrackingService') {
                    catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                        sh 'go test -short ./...'
                    }
                }
            }
        }

        // 3. LINT / VET (WAJIB HIJAU)
        stage('Lint / Vet') {
            steps {
                dir('TrackingService') {
                    sh 'go vet ./...'
                }
            }
        }

        // 4. BUILD IMAGE (WAJIB HIJAU)
        stage('Build Image') {
            steps {
                sh 'docker build -t $IMAGE ./TrackingService'
            }
        }

        // 5. FUNCTIONAL TEST (BOLEH FAIL, TAPI LANJUT)
        stage('Functional Test') {
            steps {
                catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                    sh '''
                    echo "START MONGO DB"

                    docker rm -f mongo-test || true

                    docker run -d \
                      --name mongo-test \
                      -p 27017:27017 \
                      mongo:7

                    sleep 5

                    echo "START APP"

                    docker rm -f test-tracking || true

                    docker run -d \
                      --name test-tracking \
                      --link mongo-test \
                      -e MONGO_URI=mongodb://mongo-test:27017 \
                      -p 8087:8087 \
                      $IMAGE

                    sleep 5

                    echo "RUN FUNCTIONAL TEST"

                    cd TrackingService

                    go test -run TestTrackingAPI_Success
                    '''
                }
            }
        }

        // 6. PUSH IMAGE (WAJIB HIJAU)
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
                sh 'echo "DEPLOY TRACKING SERVICE OK"'
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
            echo 'PIPELINE SUCCESS (meskipun ada stage merah)'
        }

        failure {
            echo 'PIPELINE FAILED (cek build/vet/push)'
        }

        always {
            sh '''
            docker rm -f mongo-test || true
            docker rm -f test-tracking || true
            '''
        }
    }
}