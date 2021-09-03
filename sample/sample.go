package sample

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-atomci/workflow"
	"github.com/go-atomci/workflow/jenkins"
	"github.com/go-atomci/workflow/jenkins/templates"
)

// NewWorkFlowProvide new workflow provide
func NewWorkFlowProvide(driver, addr, user, token, jobName string, flowProcessor jenkins.FlowProcessor) (workflow.WorkFlow, error) {
	var err error
	var workFlowProvider workflow.WorkFlow
	switch {
	case driver == workflow.DriverJenkins.String():
		workFlowProvider, err = jenkins.NewJenkinsClient(
			jenkins.URL(addr),
			jenkins.JenkinsUser(user),
			jenkins.JenkinsToken(token),
			jenkins.JenkinsJob(jobName),
			jenkins.Processor(flowProcessor))
		if err != nil {
			log.Print(err)
			return nil, err
		}
		return workFlowProvider, nil
	}
	log.Print("work flow system not configured")
	return nil, fmt.Errorf("work flow system not configured")
}

func GetPipelineXMLStr() string {
	taskPipelineXMLStrArr := []string{}

	appCheckoutItems := []jenkins.StepItem{
		{
			Name:    "app01",
			Command: "sh 'python3  checkout.py'",
		},
		{
			Name:    "app02",
			Command: "sh 'python3 checkout.py'",
		},
	}

	items := map[string]interface{}{"CheckoutItems": appCheckoutItems}

	//　这里仅示例了一个 templates.Checkout 类型
	taskPipelineXMLStr, err := jenkins.GeneratePipelineXMLStr(templates.Checkout, items)
	if err != nil {
		return ""
	}
	taskPipelineXMLStrArr = append(taskPipelineXMLStrArr, taskPipelineXMLStr)

	return strings.Join(taskPipelineXMLStrArr, " ")
}

func SamplePipeline(addr, user, token string) error {
	envVars := []jenkins.EnvItem{
		{Key: "JENKINS_SLAVE_WORKSPACE", Value: "/home/jenkins/agent"},
		{Key: "ACCESS_TOKEN", Value: token},
	}

	containerTemplates := []jenkins.ContainerEnv{
		{
			Name:       "jnlp",
			Image:      "colynn/jenkins-jnlp-agent:latest",
			WorkingDir: "/home/jenkins/agent",
		},
		{
			Name:       "kaniko",
			Image:      "colynn/kaniko-ex",
			WorkingDir: "/home/jenkins/agent",
			CommandArr: commandAndArgSplit("/bin/sh -c"),
			ArgsArr:    commandAndArgSplit("cat"),
		},
	}

	// Notes
	// pipelineStageStr: you can use `jenkins.GeneratePipelineXMLStr()` get the defined of pipelineXmlStr.
	flowProcessor := &jenkins.CIContext{
		RegistryAddr:       "192.168.2.10:5000",
		EnvVars:            envVars,
		ContainerTemplates: containerTemplates,
		Stages:             GetPipelineXMLStr(),
		CommonContext: jenkins.CommonContext{
			JenkinsSlaveWorkspace: "CIInfo[3]",
			AccessToken:           "adminToken",
			AtomCIServer:          "atomciServer",
		},
		CallBack: jenkins.CallbackRequest{
			Token: token,
			URL:   "callBackURL",
			Body:  "callBackRequestBody",
		},
	}

	jenkinsClient, err := NewWorkFlowProvide(workflow.DriverJenkins.String(), addr, user, token, "sample-pipeline-test", flowProcessor)
	if err != nil {
		return err
	}

	jenkinsClient.Ping()

	runID, err := jenkinsClient.Build()
	if err != nil {
		return err
	}

	// abort jenkins client
	return jenkinsClient.Abort(runID)
}

func commandAndArgSplit(itemStr string) (itemArr []string) {
	itemStr = strings.TrimSpace(itemStr)
	if itemStr == "" {
		return
	}
	itemArr = strings.Split(itemStr, " ")
	return
}
