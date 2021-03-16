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
                catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                    sh 'newman run https://www.getpostman.com/collections/761f5d0bf6cc08b8518f > log_Main.txt'
                    sh 'newman run https://www.getpostman.com/collections/2087c1d495b54086e676 > log_Mask.txt'
                    sh 'newman run https://www.getpostman.com/collections/8d49cf76196fb998d870 > log_Test.txt'
                }
            }
        }
        stage('Build & Deploy') {
            steps {
                sh 'sudo chown -R $USER:$USER ../lastError'
                sh 'go build'
                sh './DappleyWeb_Pipeline -email <Email Address> -passWord <Email Password> -main log_Main.txt -mask log_Mask.txt -test log_Test.txt'
            }
        }
        stage('Close') {
            steps {
                sh 'rm -r log_Main.txt'
                sh 'rm -r log_Mask.txt'
                sh 'rm -r log_Test.txt'
            }
        }
    }
}
```