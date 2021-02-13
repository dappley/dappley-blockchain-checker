# DappleyWeb_Pipelien

### Pipeline:

```
pipeline {
    agent any
    tools {
        go 'go-1.15.7'
    }
    environment {
        GO1157MODULE = 'on'
    }
    stages {
        stage('SCM Checkout') {
            steps {
                git 'https://github.com/heesooh/DappleyWeb_Pipeline'
            }
        }
        stage('Postman Test') {
            steps {
                sh 'newman run https://www.getpostman.com/collections/761f5d0bf6cc08b8518f --reporters=cli,htmlextra --reporter-htmlextra-export "newman/report.html"'
            }
        }
        stage('Build & Deploy') {
            steps {
                sh 'go build'
                sh './DappleyWeb_Pipeline -fileName "newman/report.html"'
            }
        }
        stage('Close') {
            steps {
                sh 'rm -r newman'
            }
        }
    }
}
```