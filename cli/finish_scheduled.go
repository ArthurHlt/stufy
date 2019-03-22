package main

import (
	"fmt"
	"github.com/ArthurHlt/stufy"
	"github.com/ArthurHlt/stufy/messages"
	"github.com/jessevdk/go-flags"
	"gopkg.in/AlecAivazis/survey.v1"
)

type FinishScheduled struct {
	InlineFlag
	Filename FilenameScheduled `short:"f" long:"filename" description:"Set filename associated to update"`
}

var finishScheduled FinishScheduled

func (c *FinishScheduled) Execute(_ []string) error {
	if c.Inline {
		return c.ExecuteInline()
	}
	return c.ExecuteSurvey()
}

func (c *FinishScheduled) ExecuteInline() error {
	if c.Filename == "" {
		return &flags.Error{
			Type:    flags.ErrRequired,
			Message: "Filename flag required (--filename, -f)",
		}
	}
	err := manager.FinishScheduled(stufy.RequestUnscheduled{
		Filename: string(c.Filename),
		Confirm:  true,
	})
	if err != nil {
		return err
	}
	messages.Printfln("Scheduled task %s has been mark as %s",
		messages.C.Cyan(c.Filename), messages.C.Green("finished"),
	)
	return nil
}

func (c *FinishScheduled) ExecuteSurvey() error {
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
