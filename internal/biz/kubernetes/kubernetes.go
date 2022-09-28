package kubernetes

import (
	"log"

	"github.com/zendesk-slowery23/survey-demo/pkg/api"
)

type ScaffoldingService interface {
	Create(flags *api.KubernetesFlags) error
}

type scaffoldingService struct {
}

func New() ScaffoldingService {
	return &scaffoldingService{}
}

func (s *scaffoldingService) Create(flags *api.KubernetesFlags) error {
	log.Printf("Here's where we scaffold a workload of type %s with %d replicas", flags.Type, flags.Replicas)
	return nil
}
