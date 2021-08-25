package templates

// DeployPipeline defained default jenkins pipeline
const DeployPipeline = `
pipeline {
    agent any
    environment {
    def JENKINS_SLAVE_WORKSPACE = '{{ .JenkinsSlaveWorkspace }}'
    def ATOMCI_SERVER = '{{ .AtomCIServer }}'
    def ACCESS_TOKEN = '{{ .AccessToken }}'
    def USER_TOKEN = '{{ .UserToken }}'
    }
    stages {
        stage('HealthCheck') {
            parallel {
                {{- range $i, $item := .HealthCheckItems }}
                stage('{{ $item.Name }}') {
                    steps {
                        {{ $item.Command }}
                    }
                }
                {{- end }}
            }
        }
        stage('Callback') {
            steps {
                retry(count: 5) {
                    httpRequest acceptType: 'APPLICATION_JSON', contentType: 'APPLICATION_JSON', customHeaders: [[maskValue: true, name: 'Authorization', value: 'Bearer {{ .CallBack.Token }}']], httpMode: 'POST', requestBody: '''{{ .CallBack.Body }}''', responseHandle: 'NONE', timeout: 10, url: '{{ .CallBack.URL }}'
                }
            }
        }
    }
}
`
