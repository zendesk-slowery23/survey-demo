package main

import (
	"log"

	"github.com/zendesk-slowery23/survey-demo/internal/cmd"
)

func main() {

	rootCmd := cmd.BuildRoot()

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
