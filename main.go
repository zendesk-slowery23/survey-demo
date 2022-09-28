package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/AlecAivazis/survey/v2"
)

type answers struct {
	Name          string
	Quest         string
	FavoriteColor string `survey:"color"`
	Beatles       []string
}

func main() {

	var ready bool

	err := survey.AskOne(&survey.Confirm{
		Message: "Are you ready?",
		Default: true,
	}, &ready)

	if err != nil {
		log.Fatal(err)
	}

	if !ready {
		log.Fatal("Try again when you're ready")
	}

	qs := []*survey.Question{
		{
			Name: "name",
			Prompt: &survey.Input{
				Message: "What is your name?",
				Default: os.Getenv("USER"),
			},
			Validate:  survey.Required,
			Transform: survey.Title,
		},
		{
			Name: "quest",
			Prompt: &survey.Input{
				Message: "What is your quest?",
			},
			Validate: survey.Required,
		},
		{
			Name: "color",
			Prompt: &survey.Select{
				Message: "What is your favorite Color?",
				Options: []string{
					"🔴 Red",
					"🟠 Orange",
					"🟡 Yellow",
					"🟢 Green",
					"🔵 Blue",
					"🟣 Purple",
					"⚫ Black",
					"⚪ White",
				},
			},
			Validate: survey.Required,
		},
		{
			Name: "beatles",
			Prompt: &survey.MultiSelect{
				Message: "Which Beatles do you like?",
				Options: []string{"John", "Paul", "George", "Ringo"},
			},
			Validate: survey.Required,
		},
	}

	ans := &answers{}
	err = survey.Ask(qs, ans)
	if err != nil {
		log.Fatal(err)
	}

	ansJson, _ := json.Marshal(ans)
	log.Printf("Answers given were: %s", ansJson)
}
