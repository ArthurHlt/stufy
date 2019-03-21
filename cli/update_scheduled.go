package main

import (
	"fmt"
	"github.com/ArthurHlt/stufy"
	"gopkg.in/AlecAivazis/survey.v1"
	"time"
)

type UpdateScheduled struct {
	All bool `short:"a" long:"all" description:"Show all scheduled tasks (by default it doesn't show finished scheduled)'"`
}

var updateScheduled UpdateScheduled

func (c *UpdateScheduled) Execute(_ []string) error {
	config, err := manager.Config()
	if err != nil {
		return err
	}
	scheduled, err := manager.ListScheduled(c.All)
	if err != nil {
		return err
	}

	mType := " non resolved "
	if c.All {
		mType = ""
	}
	if len(scheduled) == 0 {
		fmt.Printf("There is no %s scheduled tasks\n", mType)
		return nil
	}

	filename := ""
	err = survey.AskOne(&survey.Select{
		Message: "What scheduled task do you want to update?",
		Options: scheduled.Filenames(),
	}, &filename, survey.Required)
	if err != nil {
		return err
	}

	currentScheduled, err := manager.FindIncident(filename)
	if err != nil {
		return err
	}
	qs := []*survey.Question{
		{
			Name: "Systems",
			Prompt: &survey.MultiSelect{
				Message: "What are the affected systems?",
				Options: config.Content.Systems,
				Default: currentScheduled.AffectedSystems,
			},
			Validate: survey.Required,
		},
		{
			Name:   "Description",
			Prompt: &survey.Input{Message: "Add a concise description of the scheduled task (empty do not override actual description)."},
		},
		{
			Name: "Date",
			Prompt: &survey.Input{
				Message: "When will the scheduled task will start (YYYY-mm-ddTHH:MM)?",
				Default: time.Time(*currentScheduled.Scheduled).Format("2006-01-02T15:04"),
			},
			Validate: survey.Required,
		},
		{
			Name: "Duration",
			Prompt: &survey.Input{
				Message: "How long the scheduled task will take?",
				Default: fmt.Sprintf("%dm", currentScheduled.Duration),
			},
			Validate: survey.Required,
		},
		{
			Name: "Confirm",
			Prompt: &survey.Confirm{
				Message: "Are you sure you want to update this scheduled task?",
				Default: true,
			},
		},
		{
			Name: "Open",
			Prompt: &survey.Confirm{
				Message: "Open the scheduled task for editing?",
				Default: false,
			},
		},
	}

	var req stufy.RequestUpdateScheduled
	err = survey.Ask(qs, &req)
	if err != nil {
		return err
	}
	req.Filename = filename
	return manager.UpdateScheduled(req)
}

func init() {
	desc := `Update an existing scheduled task`
	c, err := parser.AddCommand(
		"update-scheduled",
		desc,
		desc,
		&updateScheduled)
	if err != nil {
		panic(err)
	}
	c.Aliases = []string{"us"}
}
