package kubernetes

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/zendesk-slowery23/survey-demo/pkg/api"
)

func Wizard(flags *api.KubernetesFlags) error {

	return survey.Ask([]*survey.Question{
		{
			Name: "type",
			Prompt: &survey.Select{
				Message: "What type of workload are you deploying?",
				Options: []string{
					"deployment",
					"statefulset",
				},
			},
			Validate: survey.Required,
		},
		{
			Name: "replicas",
			Prompt: &survey.Input{
				Message: "How many replicas do you require?",
				Default: "2",
			},
			Validate: survey.Required,
		},
	}, flags)

}
