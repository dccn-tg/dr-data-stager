pipeline {
    agent any

    environment {
        DOCKER_IMAGE_TAG = "jenkins-${env.BUILD_NUMBER}"
    }
    
    options {
        disableConcurrentBuilds()
    }

    stages {
        // build containers.
        stage('Build') {
            when {
                expression {
                    return !params.PRODUCTION
                }
            }
            agent {
                label 'swarm-manager'
            }
            steps {
                sh 'docker-compose build --parallel'
            }
        }

        // build contains with tags ready for pushing to production registry.
        stage('Build (PRODUCTION)') {
            when {
                expression {
                    return params.PRODUCTION
                }
            }
            agent {
                label 'swarm-manager'
            }
            steps {
                withEnv([
                    "DOCKER_REGISTRY=${params.PRODUCTION_DOCKER_REGISTRY}",
                    "DOCKER_IMAGE_TAG=${params.PRODUCTION_GITHUB_TAG}"
                ]) {
                    sh 'docker-compose build --parallel'
                }
            }
        }

        // perform unit test, and push the containers to the test environment registry on success.
        stage('Unit tests') {
            when {
                expression {
                    return !params.PRODUCTION
                }
            }
            agent {
                docker {
                    image "${env.DOCKER_REGISTRY}/data-stager-api:${env.DOCKER_IMAGE_TAG}"
                    reuseNode true
                }
            }
            steps {
                echo 'Unit tests go here'
            }
            post {
                success {
                    script {
                        if (env.DOCKER_REGISTRY) {
                            sh (
                                label: "Pushing images to repository '${env.DOCKER_REGISTRY}'",
                                script: 'docker-compose push'
                            )
                        }
                    }
                }
            }
        }

        // deploying containers on the test environment.
        stage('Staging') {
             when {
                expression {
                    return !params.PRODUCTION
                }
            }
            stages {
                stage('Deploying containers') {
                    steps {

                        sh (
                            label: 'Removing previous stack',
                            script: 'docker stack rm dr-data-stager'
                        )
                        sleep (
                            time: 30,
                            unit: 'SECONDS'
                        )
                        withCredentials ([
                            usernamePassword (
                                credentialsId: params.AUTH_CLIENT_CREDENTIAL,
                                usernameVariable: 'AUTH_CLIENT_ID',
                                passwordVariable: 'AUTH_CLIENT_SECRET'
                            ),
                        ]) {
                            sh (
                                label: 'Deploying new stack',
                                script: 'docker stack up -c docker-compose.yml -c docker-compose.swarm.yml --with-registry-auth dr-data-stager'
                            )
                        }
                    }
                }

                stage('Health check') {
                    agent {
                        label 'swarm-manager'
                    }
                    steps {
                        withDockerContainer(image: 'jwilder/dockerize', args: '--network dr-data-stager-net') {
                            sh (
                                label: 'Waiting for services to become available',
                                script: 'dockerize \
                                    -timeout 120s \
                                    -wait tcp://db:6379 \
                                    -wait http://api-server:8080 \
                                    -wait http://ui:3080'
                            )
                        }
                    }
                    post {
                        failure {
                            sh (
                                label: 'Displaying service status',
                                script: 'docker stack ps dr-data-stager'
                            )
                            sh (
                                label: 'Displaying service logs',
                                script: 'docker stack services --format \'{{.Name}}\' dr-data-stager | xargs -n 1 docker service logs'
                            )
                        }
                    }
                }
            }
        }

        // making release tag and push containers to production registry.
        stage('Tag and push (PRODUCTION)') {
            when {
                expression {
                    return params.PRODUCTION
                }
            }
            agent {
                docker {
                    image "${env.DOCKER_REGISTRY}/data-stager-api:${env.DOCKER_IMAGE_TAG}"
                    reuseNode true
                }
            }
            steps {
                echo "production: true"
                echo "production github tag: ${params.PRODUCTION_GITHUB_TAG}"

                // Handle Github tags
                withCredentials([
                    usernamePassword (
                        credentialsId: params.GITHUB_CREDENTIAL,
                        usernameVariable: 'GITHUB_USERNAME',
                        passwordVariable: 'GITHUB_PASSWORD'
                    )
                ]) {
                    // Remove local tag (if any)
                    script {
                        def statusCode = sh(script: "git tag --list | grep ${params.PRODUCTION_GITHUB_TAG}", returnStatus: true)
                        if(statusCode == 0) {
                            sh "git tag -d ${params.PRODUCTION_GITHUB_TAG}"
                            echo "Removed existing local tag ${params.PRODUCTION_GITHUB_TAG}"
                        }
                    }
                    
                    // Create local tag
                    sh "git tag -a ${params.PRODUCTION_GITHUB_TAG} -m 'jenkins'"
                    echo "Created local tag ${params.PRODUCTION_GITHUB_TAG}"

                    // Remove remote tag (if any)
                    script {
                        def result = sh(script: "git ls-remote https://${GITHUB_USERNAME}:${GITHUB_PASSWORD}@github.com/dccn-tg/di-data-stager.git refs/tags/${params.PRODUCTION_GITHUB_TAG}", returnStdout: true).trim()
                        if (result != "") {
                            sh "git push --delete https://${GITHUB_USERNAME}:${GITHUB_PASSWORD}@github.com/dccn-tg/di-data-stager.git ${params.PRODUCTION_GITHUB_TAG}"
                            echo "Removed existing remote tag ${params.PRODUCTION_GITHUB_TAG}"
                        }
                    }

                    // Create remote tag
                    sh "git push https://${GITHUB_USERNAME}:${GITHUB_PASSWORD}@github.com/dccn-tg/di-data-stager.git ${params.PRODUCTION_GITHUB_TAG}"
                    echo "Created remote tag ${params.PRODUCTION_GITHUB_TAG}"
                }

                // Override the env variables and 
                // push the Docker images to the production Docker registry
                withEnv([
                    "DOCKER_REGISTRY=${params.PRODUCTION_DOCKER_REGISTRY}",
                    "DOCKER_IMAGE_TAG=${params.PRODUCTION_GITHUB_TAG}"
                ]) {
                    withCredentials([
                        usernamePassword (
                            credentialsId: params.PRODUCTION_DOCKER_REGISTRY_CREDENTIAL,
                            usernameVariable: 'DOCKER_USERNAME',
                            passwordVariable: 'DOCKER_PASSWORD'
                        )
                    ]) {
                        sh "docker login -u ${DOCKER_USERNAME} -p ${DOCKER_PASSWORD} ${params.PRODUCTION_DOCKER_REGISTRY}"
                        sh 'docker-compose push'
                        echo "Pushed images to ${DOCKER_REGISTRY}"
                    }
                } 
            }
        }
    }

    post {
        success {
            script {
                // regenerate env.sh; but strip out the username/password
                def statusCode = sh(returnStatus:true, script: "bash ./print_env.sh | sed 's/^AUTH_CLIENT_ID=.*/AUTH_CLIENT_ID=/' | sed 's/^AUTH_CLIENT_SECRET=.*/AUTH_CLEINT_SECRET=/' > env")
                if ( statusCode != 0 ) {
                    echo "unable to generate env.sh file, check it manually."
                }
            }
            archiveArtifacts "docker-compose.yml, docker-compose.swarm.yml env"
        }
        always {
            echo 'cleaning'
            sh 'docker system prune -f'
        }
    }
}
