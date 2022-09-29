package kubernetes

import (
	"github.com/spf13/cobra"
	survey "github.com/zendesk-slowery23/survey-demo/internal/cmd/survey/kubernetes"
	"github.com/zendesk-slowery23/survey-demo/pkg/biz/kubernetes"
)

func New(flags *kubernetes.Flags) *cobra.Command {

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

	k8s.PersistentFlags().BoolVarP(&flags.Interactive, "interactive", "i", false, "Interactive Mode")
	k8s.AddCommand(man)

	gen := &cobra.Command{
		Use:   "generate",
		Short: "Generate kubernetes manifests for your workload",
		RunE: func(cmd *cobra.Command, args []string) error {
			if it, _ := cmd.Flags().GetBool("interactive"); it {
				err := survey.Wizard(flags)
				if err != nil {
					return err
				}
			}

			kubernetes.New().Create(flags)
			return nil
		},
	}

	gen.Flags().StringVarP(&flags.Type, "type", "t", "deployment", "Type (one of deployment or statefulset)")
	gen.Flags().IntVarP(&flags.Replicas, "replicas", "r", 2, "# of replicas")

	man.AddCommand(gen)
	return k8s
}
