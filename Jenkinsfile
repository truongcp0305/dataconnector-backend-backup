pipeline{
    agent any
    environment{
        SERVICE_NAME = "data-connector.symper.vn"
        SERVICE_ENV = "prod"
    }
    stages{
        stage("build"){
            steps{
                withCredentials([usernamePassword(credentialsId: 'docker-hub', passwordVariable: 'DOCKER_REGISTRY_PWD', usernameVariable: 'DOCKER_REGISTRY_USER')]) {
                    sh 'echo $DOCKER_REGISTRY_PWD | docker login -u $DOCKER_REGISTRY_USER --password-stdin localhost:5000'
                }
                script {
                    latestTag = sh(returnStdout:  true, script: "git tag --sort=-creatordate | head -n 1").trim()
                    env.BUILD_VERSION = latestTag
                    sh "docker build -t localhost:5000/${SERVICE_NAME}:${env.BUILD_VERSION} ."
                    sh "docker push localhost:5000/${SERVICE_NAME}:${env.BUILD_VERSION}"
                    sh "docker image rm localhost:5000/${SERVICE_NAME}:${env.BUILD_VERSION}"
                }
            }
        }
        stage("deploy to k8s"){
            steps{
                withCredentials([usernamePassword(credentialsId: 'data_connector_db', passwordVariable: 'POSTGRES_PASS', usernameVariable: 'POSTGRES_USER')]) {
                    sh "chmod +x changeTag.sh"
                    sh './changeTag.sh $SERVICE_NAME:$BUILD_VERSION $SERVICE_ENV $POSTGRES_USER $POSTGRES_PASS'
                    sshagent(['ssh-remote']) {
                        sh "ssh root@103.148.57.32 rm -rf /root/kubernetes/deployment/${SERVICE_ENV}/${SERVICE_NAME}"
                        sh "ssh root@103.148.57.32 mkdir /root/kubernetes/deployment/${SERVICE_ENV}/${SERVICE_NAME}"
                        sh "scp -o StrictHostKeyChecking=no k8s/* root@103.148.57.32:/root/kubernetes/deployment/${SERVICE_ENV}/${SERVICE_NAME}"
                        sh "ssh root@103.148.57.32 kubectl config set-context --current --namespace=${SERVICE_ENV}"
                        sh "ssh root@103.148.57.32 kubectl apply -f /root/kubernetes/deployment/${SERVICE_ENV}/${SERVICE_NAME}"
                    }
                }
            }
        }
    }
    post{
        always{
            emailext body: '$PROJECT_NAME - Build # $BUILD_NUMBER - $BUILD_STATUS: \nCheck console output at $BUILD_URL to view the results.',
            subject: '$PROJECT_NAME - Build # $BUILD_NUMBER - $BUILD_STATUS!',
            to: 'hoangnd@symper.vn'
        }
        success{
            echo "========pipeline executed successfully ========"
        }
        failure{
            echo "========pipeline execution failed========"
        }
    }
}