package main

import (
	"fmt"
	"github.com/ArthurHlt/stufy"
	"gopkg.in/AlecAivazis/survey.v1"
)

type FinishScheduled struct {
}

var finishScheduled FinishScheduled

func (c *FinishScheduled) Execute(_ []string) error {
	scheduleds, err := manager.ListScheduled(false)
	if err != nil {
		return err
	}

	mType := " non resolved "
	if len(scheduleds) == 0 {
		fmt.Printf("There is no %s scheduled tasks\n", mType)
		return nil
	}

	qs := []*survey.Question{
		{
			Name: "Filename",
			Prompt: &survey.Select{
				Message: "What scheduled do you want to finish?",
				Options: scheduleds.Filenames(),
			},
			Validate: survey.Required,
		},
		{
			Name: "Confirm",
			Prompt: &survey.Confirm{
				Message: "Are you sure you want to finish this scheduled task?",
				Default: true,
			},
		},
	}

	var req stufy.RequestUnscheduled
	err = survey.Ask(qs, &req)
	if err != nil {
		return err
	}
	return manager.FinishScheduled(req)
}

func init() {
	desc := `Finish a scheduled task`
	c, err := parser.AddCommand(
		"finish-scheduled",
		desc,
		desc,
		&finishScheduled)
	if err != nil {
		panic(err)
	}
	c.Aliases = []string{"fs"}
}
