package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/zendesk-slowery23/survey-demo/internal/cmd/cobra/kubernetes"
	"github.com/zendesk-slowery23/survey-demo/pkg/api"
)

func DoStuff() {
	root := &cobra.Command{
		Use:  "cicd",
		Long: "CICD toolkit command line interface",
	}

	root.AddCommand(kubernetes.New(&api.KubernetesFlags{}))

	err := root.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
