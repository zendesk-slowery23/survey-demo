package kubernetes

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/zendesk-slowery23/survey-demo/internal/util"
	"github.com/zendesk-slowery23/survey-demo/pkg/biz/kubernetes"
)

const (
	TypeCronjob     = "cronjob"
	TypeDaemonset   = "daemonset"
	TypeDeployment  = "deployment"
	TypeStatefulset = "statefulset"

	ProbeStartup   = "startup"
	ProbeLiveness  = "liveness"
	ProbeReadiness = "readiness"
)

func Wizard(flags *kubernetes.Flags) error {

	err := survey.Ask([]*survey.Question{
		{
			Name: "name",
			Prompt: &survey.Input{
				Message: "What is the name of your workload?",
				Default: util.CurrentDir(),
			},
			Validate: survey.Required,
		},
		{
			Name: "type",
			Prompt: &survey.Select{
				Message: "What type of workload are you deploying?",
				Options: []string{
					TypeCronjob,
					TypeDaemonset,
					TypeDeployment,
					TypeStatefulset,
				},
				Default: TypeDeployment,
			},
			Validate: survey.Required,
		},
	}, flags)

	if err != nil {
		return err
	}

	switch flags.Type {
	case TypeDeployment, TypeStatefulset:
		err = survey.AskOne(&survey.Input{
			Message: "How many replicas do you require?",
			Default: "2",
		}, &flags.Replicas)
		if err != nil {
			return err
		}
	}

	if flags.Type != TypeCronjob {
		count := 1
		err = survey.AskOne(&survey.Input{
			Message: "How many ports will your workload expose?",
			Default: "1",
		}, &count)

		if err != nil {
			return err
		}

		if count > 0 {
			flags.Ports = []kubernetes.Port{}
		}
		for i := 0; i < count; i++ {
			port := kubernetes.Port{}
			err = survey.Ask([]*survey.Question{
				{
					Name: "name",
					Prompt: &survey.Input{
						Message: "What is the name of your port?",
						Default: "http",
					},
					Validate: survey.Required,
				},
				{
					Name: "containerNumber",
					Prompt: &survey.Input{
						Message: "What is the port number for your CONTAINER?",
						Default: "8080",
					},
				},
			}, &port)

			if err != nil {
				return err
			}

			err = survey.Ask([]*survey.Question{
				{
					Name: "serviceNumber",
					Prompt: &survey.Input{
						Message: "What is the port number for your SERVICE?",
						Default: strconv.Itoa(port.ContainerNumber),
					},
				},
				{
					Name: "protocol",
					Prompt: &survey.Input{
						Message: "What is the protocol of your port?",
						Default: "tcp",
					},
					Validate:  survey.Required,
					Transform: func(ans interface{}) (newAns interface{}) { return strings.ToUpper(ans.(string)) },
				},
			}, &port)

			if err != nil {
				return err
			}

			flags.Ports = append(flags.Ports, port)
		}
	}

	err = surveyImage(flags)
	if err != nil {
		return err
	}

	err = surveyResourceRequirements(flags)
	if err != nil {
		return err
	}

	return nil
}

func surveyImage(flags *kubernetes.Flags) error {

	const (
		gcrStandard         = "Google Container Registry Standard Location"
		gcrStandardLocation = "gcr.io/docker-images-180022/apps"
		ecrStandard         = "ECR Standard Location"
		ecrStandardLocation = "713408432298.dkr.ecr.us-west-2.amazonaws.com/zendesk"
		other               = "Other"
	)

	var loc string

	err := survey.AskOne(&survey.Select{
		Message: "Where can your docker image be found?",
		Options: []string{gcrStandard, ecrStandard, other},
		Default: gcrStandard,
	}, &loc)
	if err != nil {
		return err
	}

	switch loc {
	case gcrStandard:
		flags.Image = fmt.Sprintf("%s/%s", gcrStandardLocation, flags.Name)
	case ecrStandard:
		flags.Image = fmt.Sprintf("%s/%s", ecrStandardLocation, flags.Name)
	case other:
		err = survey.AskOne(&survey.Input{
			Message: "Specify the location of your docker image",
		}, &flags.Image)
		if err != nil {
			return err
		}
	}

	err = survey.Ask([]*survey.Question{
		{
			Name: "imageTag",
			Prompt: &survey.Input{
				Message: "What tag of your image do you wish to deploy?",
				Default: "latest",
			},
			Validate: survey.Required,
		},
		{
			Name: "imagePullPolicy",
			Prompt: &survey.Select{
				Message: "What should the pull policy for the image be?",
				Options: []string{"Always", "IfNotPresent"},
				Default: "IfNotPresent",
			},
		},
	}, flags)
	if err != nil {
		return err
	}

	return nil
}

func surveyResourceRequirements(flags *kubernetes.Flags) error {

	err := survey.Ask([]*survey.Question{
		{
			Name: "cpuRequest",
			Prompt: &survey.Input{
				Message: "How much CPU would you like to request for your workload?",
				Default: "250m",
				Help:    "See xxx for meaningful values",
			},
			Validate: survey.Required,
		},
		{
			Name: "memoryRequest",
			Prompt: &survey.Input{
				Message: "How much memory would you like to request for your workload?",
				Default: "128M",
				Help:    "See xxx for meaningful values",
			},
			Validate: survey.Required,
		},
	}, flags)
	if err != nil {
		return err
	}

	//TODO as the questions for limits; for now just copy
	flags.CpuLimit = flags.CpuRequest
	flags.MemoryLimit = flags.MemoryRequest

	return nil
}
