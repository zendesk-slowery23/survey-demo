package kubernetes

import (
	"log"
)

type Flags struct {
	Interactive bool
	Type        string
	Replicas    int
}

type ScaffoldingService interface {
	Create(flags *Flags) error
}

type scaffoldingService struct {
}

func New() ScaffoldingService {
	return &scaffoldingService{}
}

func (s *scaffoldingService) Create(flags *Flags) error {
	log.Printf("Here's where we scaffold a workload of type %s with %d replicas", flags.Type, flags.Replicas)
	return nil
}
