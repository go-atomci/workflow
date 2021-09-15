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
            {{ if .HealthCheckItems }}
            parallel {
                {{- range $i, $item := .HealthCheckItems }}
                stage('{{ $item.Name }}') {
                    steps {
                        {{ $item.Command }}
                    }
                }
                {{- end }}
            }
            {{ else }}
                steps {
                    sh "there was no health ckeck items"
                }
            {{ end }}
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
