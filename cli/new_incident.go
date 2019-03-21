package main

import (
	"github.com/ArthurHlt/stufy"
	"github.com/ArthurHlt/stufy/model"
	"gopkg.in/AlecAivazis/survey.v1"
)

type NewIncident struct {
}

var newIncident NewIncident

func (c *NewIncident) Execute(_ []string) error {
	config, err := manager.Config()
	if err != nil {
		return err
	}
	qs := []*survey.Question{
		{
			Name:      "Cause",
			Prompt:    &survey.Input{Message: "What is the cause of the incident?"},
			Validate:  survey.Required,
			Transform: survey.Title,
		},
		{
			Name: "Severity",
			Prompt: &survey.Select{
				Message: "What is the severity of the incident?",
				Options: model.SeveritiesString(),
			},
			Validate: survey.Required,
		},
		{
			Name: "Systems",
			Prompt: &survey.MultiSelect{
				Message: "What are the affected systems?",
				Options: config.Content.Systems,
			},
			Validate: survey.Required,
		},
		{
			Name:     "Description",
			Prompt:   &survey.Input{Message: "Add a concise description of the incident."},
			Validate: survey.Required,
		},
		{
			Name: "Open",
			Prompt: &survey.Confirm{
				Message: "Open the incident for editing?",
				Default: false,
			},
		},
	}

	var req stufy.RequestCreate
	err = survey.Ask(qs, &req)
	if err != nil {
		return err
	}

	return manager.CreateIncident(req)
}

func init() {
	desc := `Create a new incident`
	c, err := parser.AddCommand(
		"new-incident",
		desc,
		desc,
		&newIncident)
	if err != nil {
		panic(err)
	}
	c.Aliases = []string{"n"}
}
