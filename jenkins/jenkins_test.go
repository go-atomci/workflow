/*
Copyright 2021 The AtomCI Group Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package jenkins

import (
	"strings"
	"testing"

	"github.com/go-atomci/workflow/jenkins/templates"
	"github.com/stretchr/testify/assert"
)

func TestGetStatus(t *testing.T) {
	appCheckoutItems := []StepItem{
		{
			Name:    "app01",
			Command: "sh 'echo hello'",
		},
	}
	items := map[string]interface{}{"CheckoutItems": appCheckoutItems}

	taskPipelineXMLStrArr := []string{}

	buildItems := map[string]interface{}{"BuildItems": []StepItem{}}
	imageItems := map[string]interface{}{"ImageItems": []StepItem{}}

	taskPipelineXMLStr, _ := GeneratePipelineXMLStr(templates.Checkout, items)
	compileTasks, _ := GeneratePipelineXMLStr(templates.Compile, buildItems)
	buildTasks, _ := GeneratePipelineXMLStr(templates.BuildImage, imageItems)
	taskPipelineXMLStrArr = append(taskPipelineXMLStrArr, taskPipelineXMLStr)
	taskPipelineXMLStrArr = append(taskPipelineXMLStrArr, compileTasks)
	taskPipelineXMLStrArr = append(taskPipelineXMLStrArr, buildTasks)
	pipelineStagesStr := strings.Join(taskPipelineXMLStrArr, " ")

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
		Stages: pipelineStagesStr,
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
		token:   "xxx",
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
