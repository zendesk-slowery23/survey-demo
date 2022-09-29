package spinnaker

import (
	"log"

	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	spin := &cobra.Command{
		Use:   "spinnaker",
		Short: "Commands related to spinnaker",
		Long:  "Commands related to spinnaker",
	}

	pipeline := &cobra.Command{
		Use:     "pipeline",
		Aliases: []string{"pipe"},
		Short:   "Commands related to spinnaker pipeline generation",
		Long:    "Commands related to spinnaker pipeline generation",
	}

	gen := &cobra.Command{
		Use:     "generate",
		Aliases: []string{"create"},
		Short:   "Generate a spinnaker pipeline",
		Long:    "Generate a spinnaker pipeline",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("Here's where we generate spinnaker pipelines")
		},
	}

	spin.AddCommand(pipeline)
	pipeline.AddCommand(gen)

	return spin
}
