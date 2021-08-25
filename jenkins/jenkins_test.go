package jenkins

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetStatus(t *testing.T) {
	CIContext := CIContext{
		RegistryAddr: "",
		EnvVars: []EnvItem{
			{Key: "JENKINS_SLAVE_WORKSPACE", Value: "/tmp"},
			{Key: "ACCESS_TOKEN", Value: "dddxxx"},
			{Key: "REPO_CNF", Value: "{\"rsgitlab.example.com\": [\"colynn\",\"jkW7i-xxxx-\"]}"},
			{Key: "DOCKER_AUTH", Value: "xxxddsxxw="},
			{Key: "DOCKER_CONFIG", Value: "/kaniko/.docker"},
		},

		// TODO: add container image to here
		ContainerTemplates: []ContainerEnv{
			{
				Name:       "jnlp",
				Image:      "172.16.1.8:5000/library/jenkins:jnlp-debug",
				WorkingDir: "/home/jenkins/agent",
				CommandArr: []string{},
			},
			{
				Name:       "kaniko",
				Image:      "172.16.18:5000/library/kaniko-executor:latest",
				WorkingDir: "/home/jenkins/agent",
				CommandArr: []string{"/bin/sh", "-c"},
				ArgsArr:    []string{"cat"},
			},
		},
		CallBack: CallbackRequest{
			Token: "xxx",
			URL:   "http://atomci-server",
			Body:  "{\"publish_job_id\": 1}",
		},
	}
	var err error

	worker := jenkinsWorker{
		url:     "http://127.0.0.1:8080",
		user:    "admin",
		token:   "110e88c6c52e02b5bdea82f73151832df9",
		jobName: "atomci-unit-test",
	}
	var xmlConfig string

	worker.jobName = "atomci-unit-test-ci-pipeline"
	t.Run("Build Default", func(t *testing.T) {
		xmlConfig, err = CIContext.GetCIPipelineXML(CIContext)
		assert.Nil(t, err)
		if err != nil {
			t.Fatalf("parse pipelinexml failed: %v", err)
		}
	})
	t.Run("Create Jenkins Job", func(t *testing.T) {
		err := worker.CreateOrUpdateJob(xmlConfig)
		assert.Nil(t, err)
		if err != nil {
			t.Fatalf("create jenkins job failed: %v", err)
		}
	})
}
