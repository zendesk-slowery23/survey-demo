package cmd

import (
	"github.com/spf13/cobra"
	k8scmd "github.com/zendesk-slowery23/survey-demo/internal/cmd/cobra/kubernetes"
	"github.com/zendesk-slowery23/survey-demo/pkg/biz/kubernetes"
)

func BuildRoot() *cobra.Command {
	root := &cobra.Command{
		Use:  "cicd",
		Long: "CICD toolkit command line interface",
	}

	interactive := &cobra.Command{
		Use:   "interactive",
		Short: "Interactive mode for subsequent commands.",
		Long:  "Interactive mode for subsequent commands. Rather than flags, the user will go through a Survey workflow.",
	}

	root.AddCommand(interactive)

	for parent, it := range map[*cobra.Command]bool{
		root:        false,
		interactive: true,
	} {
		parent.AddCommand(k8scmd.New(&kubernetes.Flags{}, it))
		// parent.AddCommand(spincmd.New())
	}

	return root
}
