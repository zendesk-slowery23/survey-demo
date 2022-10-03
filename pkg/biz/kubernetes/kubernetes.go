package kubernetes

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ghodss/yaml"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type Flags struct {
	Name            string
	Type            string
	Replicas        int
	Image           string
	ImageTag        string
	ImagePullPolicy string
	Ports           []Port

	CpuRequest    string
	CpuLimit      string
	MemoryRequest string
	MemoryLimit   string
}

type Port struct {
	Name            string
	ContainerNumber int
	ServiceNumber   int
	Protocol        string
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

	var workload interface{}

	switch flags.Type {
	case "deployment":
		workload = s.toDeployment(flags)
	case "statefulset":
		workload = s.toStatefulset(flags)
	}

	svc := s.toService(flags)

	// NOTE: we are converting to json first and then to yaml instead of directly to yaml because
	// the upstream kubernetes api objects contain `json:omitempty` tags, but none for yaml so encoding
	// straight to yaml yields a lot of null/empty fields

	f, err := os.Create(fmt.Sprintf("kubernetes/%s.yml", flags.Name))
	if err != nil {
		return err
	}
	defer f.Close()

	for _, o := range []interface{}{workload, svc} {
		j, err := json.Marshal(o)
		if err != nil {
			return err
		}

		y, err := yaml.JSONToYAML(j)
		if err != nil {
			return err
		}

		_, err = f.Write(y)
		if err != nil {
			return err
		}

		_, err = f.WriteString("---\n")
		if err != nil {
			return err
		}

	}

	return nil
}

func (s *scaffoldingService) toDeployment(flags *Flags) *apps.Deployment {

	repl := int32(flags.Replicas)

	return &apps.Deployment{
		TypeMeta: meta.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: *s.toObjectMeta(flags),
		Spec: apps.DeploymentSpec{
			Selector: &meta.LabelSelector{
				MatchLabels: map[string]string{
					"app.kubernetes.io/name": flags.Name,
				},
			},
			Replicas: &repl,
			Strategy: apps.DeploymentStrategy{
				RollingUpdate: &apps.RollingUpdateDeployment{
					MaxUnavailable: &intstr.IntOrString{IntVal: 1},
				},
			},
			Template: s.toPodTemplateSpec(flags),
		},
	}
}

func (s *scaffoldingService) toStatefulset(flags *Flags) *apps.StatefulSet {

	repl := int32(flags.Replicas)

	return &apps.StatefulSet{
		TypeMeta: meta.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "StatefulSet",
		},
		ObjectMeta: *s.toObjectMeta(flags),
		Spec: apps.StatefulSetSpec{
			Replicas: &repl,
			Selector: &meta.LabelSelector{
				MatchLabels: map[string]string{
					"app.kubernetes.io/name": flags.Name,
				},
			},
			Template: s.toPodTemplateSpec(flags),
		},
	}
}

func (s *scaffoldingService) toPodTemplateSpec(flags *Flags) core.PodTemplateSpec {
	return core.PodTemplateSpec{
		ObjectMeta: *s.toObjectMeta(flags),
		Spec: core.PodSpec{
			Containers: []core.Container{
				*s.toContainer(flags),
			},
		},
	}
}

func (s *scaffoldingService) toService(flags *Flags) *core.Service {
	return &core.Service{
		TypeMeta: meta.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: *s.toObjectMeta(flags),
		Spec: core.ServiceSpec{
			Selector: map[string]string{
				"app.kubernetes.io/name": flags.Name,
			},
			Ports: s.toServicePorts(flags),
		},
	}
}

func (s *scaffoldingService) toContainer(flags *Flags) *core.Container {

	return &core.Container{
		Name:            flags.Name,
		Image:           fmt.Sprintf("%s:%s", flags.Image, flags.ImageTag),
		ImagePullPolicy: core.PullPolicy(flags.ImagePullPolicy),
		Ports:           s.toContainerPorts(flags),
		Resources: core.ResourceRequirements{
			Requests: core.ResourceList{
				core.ResourceCPU:    resource.MustParse(flags.CpuRequest),
				core.ResourceMemory: resource.MustParse(flags.MemoryRequest),
			},
			Limits: core.ResourceList{
				core.ResourceCPU:    resource.MustParse(flags.CpuLimit),
				core.ResourceMemory: resource.MustParse(flags.MemoryLimit),
			},
		},
	}
}

func (s *scaffoldingService) toContainerPorts(flags *Flags) []core.ContainerPort {
	ports := []core.ContainerPort{}

	for _, p := range flags.Ports {
		ports = append(ports, core.ContainerPort{
			Name:          p.Name,
			ContainerPort: int32(p.ContainerNumber),
			Protocol:      core.Protocol(p.Protocol),
		})
	}

	return ports
}

func (s *scaffoldingService) toServicePorts(flags *Flags) []core.ServicePort {
	ports := []core.ServicePort{}

	for _, p := range flags.Ports {
		ports = append(ports, core.ServicePort{
			Name:       p.Name,
			Port:       int32(p.ServiceNumber),
			TargetPort: intstr.FromInt(p.ContainerNumber),
			Protocol:   core.Protocol(p.Protocol),
		})
	}

	return ports
}

func (s *scaffoldingService) toObjectMeta(flags *Flags) *meta.ObjectMeta {
	return &meta.ObjectMeta{
		Name:      flags.Name,
		Namespace: flags.Name,
		Labels: map[string]string{
			"app.kubernetes.io/name":       flags.Name,
			"app.kubernetes.io/version":    flags.ImageTag,
			"app.kubernetes.io/managed-by": "cicd-toolkit",
		},
	}
}
