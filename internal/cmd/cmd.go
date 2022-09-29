package cmd

import (
	"github.com/spf13/cobra"
	k8scmd "github.com/zendesk-slowery23/survey-demo/internal/cmd/cobra/kubernetes"
	spincmd "github.com/zendesk-slowery23/survey-demo/internal/cmd/cobra/spinnaker"
	"github.com/zendesk-slowery23/survey-demo/pkg/biz/kubernetes"
)

func BuildRoot() *cobra.Command {
	root := &cobra.Command{
		Use:  "cicd",
		Long: "CICD toolkit command line interface",
	}

	root.AddCommand(k8scmd.New(&kubernetes.Flags{}))
	root.AddCommand(spincmd.New())

	return root
}
