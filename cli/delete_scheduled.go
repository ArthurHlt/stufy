package main

import (
	"fmt"
	"github.com/ArthurHlt/stufy"
	"github.com/ArthurHlt/stufy/messages"
	"github.com/jessevdk/go-flags"
	"gopkg.in/AlecAivazis/survey.v1"
)

type DeleteScheduled struct {
	InlineFlag
	All      bool              `long:"all" description:"Show all scheduled tasks (by default it doesn't show finished scheduled)"`
	Filename FilenameScheduled `short:"f" long:"filename" description:"Set filename associated to update"`
}

var deleteScheduled DeleteScheduled

func (c *DeleteScheduled) Execute(_ []string) error {
	if c.Inline {
		return c.ExecuteInline()
	}
	return c.ExecuteSurvey()
}

func (c *DeleteScheduled) ExecuteInline() error {
	if c.Filename == "" {
		return &flags.Error{
			Type:    flags.ErrRequired,
			Message: "Filename flag required (--filename, -f)",
		}
	}
	err := manager.DeleteScheduled(stufy.RequestUnscheduled{
		Filename: string(c.Filename),
		Confirm:  true,
	})
	if err != nil {
		return err
	}
	messages.Printfln("Scheduled task %s has been %s",
		messages.C.Cyan(c.Filename), messages.C.Red("deleted"),
	)
	return nil
}

func (c *DeleteScheduled) ExecuteSurvey() error {
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
