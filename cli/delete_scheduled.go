package main

import (
	"fmt"
	"github.com/ArthurHlt/stufy"
	"gopkg.in/AlecAivazis/survey.v1"
)

type DeleteScheduled struct {
	All bool `short:"a" long:"all" description:"Show all scheduled tasks (by default it doesn't show finished scheduled)'"`
}

var deleteScheduled DeleteScheduled

func (c *DeleteScheduled) Execute(_ []string) error {
	scheduleds, err := manager.ListScheduled(c.All)
	if err != nil {
		return err
	}
	mType := " non resolved "
	if c.All {
		mType = ""
	}
	if len(scheduleds) == 0 {
		fmt.Printf("There is no %s scheduled tasks\n", mType)
		return nil
	}
	qs := []*survey.Question{
		{
			Name: "Filename",
			Prompt: &survey.Select{
				Message: "What incident do you want to update?",
				Options: scheduleds.Filenames(),
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
	desc := `Delete an existing scheduled task`
	c, err := parser.AddCommand(
		"delete-scheduled",
		desc,
		desc,
		&deleteScheduled)
	if err != nil {
		panic(err)
	}
	c.Aliases = []string{"ds"}
}
