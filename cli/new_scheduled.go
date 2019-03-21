package main

import (
	"github.com/ArthurHlt/stufy"
	"gopkg.in/AlecAivazis/survey.v1"
	"time"
)

type NewScheduled struct {
}

var newScheduled NewScheduled

func (c *NewScheduled) Execute(_ []string) error {
	config, err := manager.Config()
	if err != nil {
		return err
	}
	qs := []*survey.Question{
		{
			Name:      "Title",
			Prompt:    &survey.Input{Message: "What is the title of the scheduled task?"},
			Validate:  survey.Required,
			Transform: survey.Title,
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
			Prompt:   &survey.Input{Message: "Add a concise description of the scheduled task."},
			Validate: survey.Required,
		},
		{
			Name: "Date",
			Prompt: &survey.Input{
				Message: "When will the scheduled task will start (YYYY-mm-ddTHH:MM)?",
				Default: time.Now().Add(7 * 24 * time.Hour).Format("2006-01-02T15:04"),
			},
			Validate: survey.Required,
		},
		{
			Name: "Duration",
			Prompt: &survey.Input{
				Message: "How long the scheduled task will take?",
				Default: "2h",
			},
			Validate: survey.Required,
		},
		{
			Name: "Open",
			Prompt: &survey.Confirm{
				Message: "Open the scheduled task for editing?",
				Default: false,
			},
		},
	}

	var req stufy.RequestScheduled
	err = survey.Ask(qs, &req)
	if err != nil {
		return err
	}

	return manager.CreateScheduled(req)
}

func init() {
	desc := `Create a new scheduled task`
	c, err := parser.AddCommand(
		"new-scheduled",
		desc,
		desc,
		&newScheduled)
	if err != nil {
		panic(err)
	}
	c.Aliases = []string{"ns"}
}
