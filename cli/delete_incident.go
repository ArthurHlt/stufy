package main

import (
	"fmt"
	"github.com/ArthurHlt/stufy"
	"gopkg.in/AlecAivazis/survey.v1"
)

type DeleteIncident struct {
	All bool `short:"a" long:"all" description:"Show all incidents (by default it doesn't show resolved incident)'"`
}

var deleteIncident DeleteIncident

func (c *DeleteIncident) Execute(_ []string) error {
	incidents, err := manager.ListIncident(c.All)
	if err != nil {
		return err
	}

	mType := " non resolved "
	if c.All {
		mType = ""
	}
	if len(incidents) == 0 {
		fmt.Printf("There is no %s scheduled tasks\n", mType)
		return nil
	}

	qs := []*survey.Question{
		{
			Name: "Filename",
			Prompt: &survey.Select{
				Message: "What incident do you want to update?",
				Options: incidents.Filenames(),
			},
			Validate: survey.Required,
		},
		{
			Name: "Confirm",
			Prompt: &survey.Confirm{
				Message: "Are you sure you want to delete the incident?",
				Default: true,
			},
		},
	}

	var req stufy.RequestDelete
	err = survey.Ask(qs, &req)
	if err != nil {
		return err
	}
	return manager.DeleteIncident(req)
}

func init() {
	desc := `Delete an existing incident`
	c, err := parser.AddCommand(
		"delete-incident",
		desc,
		desc,
		&deleteIncident)
	if err != nil {
		panic(err)
	}
	c.Aliases = []string{"d"}
}
