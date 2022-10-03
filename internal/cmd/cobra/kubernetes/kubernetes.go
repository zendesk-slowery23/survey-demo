package kubernetes

import (
	"github.com/spf13/cobra"
	survey "github.com/zendesk-slowery23/survey-demo/internal/cmd/survey/kubernetes"
	"github.com/zendesk-slowery23/survey-demo/internal/util"
	"github.com/zendesk-slowery23/survey-demo/pkg/biz/kubernetes"
)

func New(flags *kubernetes.Flags, interactive bool) *cobra.Command {

	k8s := &cobra.Command{
		Use:     "kubernetes",
		Aliases: []string{"k8s"},
		Long:    "Commands related to kubernetes",
		Short:   "Commands related to kuberetes",
	}

	man := &cobra.Command{
		Use:     "manifests",
		Aliases: []string{"man"},
		Long:    "Commands related to kubernetes manifests",
		Short:   "Commands related to kuberetes manifests",
	}

	k8s.AddCommand(man)

	gen := &cobra.Command{
		Use:   "generate",
		Short: "Generate kubernetes manifests for your workload",
		Long:  "Generate kubernetes manifests for your workload",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if interactive {
				err := survey.Wizard(flags)
				if err != nil {
					return err
				}
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			kubernetes.New().Create(flags)
			return nil
		},
	}

	if !interactive {
		gen.Flags().StringVarP(&flags.Name, "name", "n", util.CurrentDir(), "Name of the workload")
		gen.Flags().StringVarP(&flags.Type, "type", "t", "deployment", "Type (one of deployment or statefulset)")
		gen.Flags().IntVarP(&flags.Replicas, "replicas", "r", 2, "# of replicas")
		gen.Flags().StringVarP(&flags.Image, "image", "i", util.CurrentDir(), "Image of your main container")
		gen.Flags().StringVar(&flags.ImageTag, "imageTag", "latest", "Tag of the image")
		gen.Flags().StringVar(&flags.ImagePullPolicy, "pul-policy", "IfNotPresent", "Pull Policy of the Image.  One of 'IfNotPresent' or 'Always'")
		gen.Flags().StringVar(&flags.CpuRequest, "cpu-request", "125m", "CPU Request")
		gen.Flags().StringVar(&flags.CpuLimit, "cpu-limit", "500m", "CPU Limit")
		gen.Flags().StringVar(&flags.MemoryRequest, "mem-request", "128M", "CPU Limit")
		gen.Flags().StringVar(&flags.MemoryLimit, "mem-limit", "128M", "CPU Limit")

		for _, f := range []string{"imageTag"} {
			gen.MarkFlagRequired(f)
		}
	}

	man.AddCommand(gen)
	return k8s
}
