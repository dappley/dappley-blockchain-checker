# dappley-blockchain-checker

Test the blockchains from the following servers every hour: 
- Main   : dappley.com 
- Mask   : 35.80.10.175 
- Test   : 3.16.250.102 

When one of the test cases return an error, send out the email to the DappWorks staff.

### Pipeline:

```
pipeline {
    agent any
    tools {
        go 'go-1.16.3'
    }
    environment {
        GO1163MODULE = 'on'
    }
    stages {
        stage('SCM Checkout') {
            steps {
                git 'https://github.com/heesooh/dappley-blockchain-checker'
            }
        }
        stage('Main Server') {
            steps {
                catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                    sh 'newman run https://www.getpostman.com/collections/761f5d0bf6cc08b8518f > log_Main.txt'
                }
            }
        }
        stage('Mask Server') {
            steps {
                catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                    sh 'newman run https://www.getpostman.com/collections/2087c1d495b54086e676 > log_Mask.txt'
                }
            }
        }
        stage('Test Server') {
            steps {
                catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                    sh 'newman run https://www.getpostman.com/collections/8d49cf76196fb998d870 > log_Test.txt'
                }
            }
        }
        stage('Build & Deploy') {
            steps {
                sh 'go build blockchain-checker.go'
                sh './blockchain-checker -email <EMAIL ADDRESS> -passWord <PASS WORD> -main log_Main.txt -mask log_Mask.txt -test log_Test.txt'
            }
        }
        stage('Close') {
            steps {
                sh 'rm -r log_Main.txt'
                sh 'rm -r log_Mask.txt'
                sh 'rm -r log_Test.txt'
                sh 'rm -r blockchain-checker'
            }
        }
    }
}
```
