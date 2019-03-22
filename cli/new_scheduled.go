package main

import (
	"github.com/ArthurHlt/stufy"
	"github.com/ArthurHlt/stufy/messages"
	"github.com/jessevdk/go-flags"
	"gopkg.in/AlecAivazis/survey.v1"
	"time"
)

type NewScheduled struct {
	InlineFlag
	Title       string `short:"l" long:"title" description:"Title of the scheduled task"`
	Description string `short:"d" long:"description" description:"Description of the incident"`
	Date        string `long:"date" description:"Date when task will start (YYYY-mm-ddTHH:MM)"`
	Duration    string `short:"u" long:"duration" description:"Duration of your task"`
}

var newScheduled NewScheduled

func (c *NewScheduled) Execute(_ []string) error {
	if c.Inline {
		return c.ExecuteInline()
	}
	return c.ExecuteSurvey()
}

func (c *NewScheduled) ExecuteInline() error {
	if len(c.Systems) == 0 {
		return &flags.Error{
			Type:    flags.ErrRequired,
			Message: "Affected system flag required (--system, -s)",
		}
	}
	if c.Title == "" {
		return &flags.Error{
			Type:    flags.ErrRequired,
			Message: "Title flag required (--title, -l)",
		}
	}
	if c.Date == "" {
		return &flags.Error{
			Type:    flags.ErrRequired,
			Message: "Date flag required (--date)",
		}
	}
	if c.Duration == "" {
		return &flags.Error{
			Type:    flags.ErrRequired,
			Message: "Duration flag required (--date)",
		}
	}
	if c.Description == "" {
		return &flags.Error{
			Type:    flags.ErrRequired,
			Message: "Description flag, (--description, -d)",
		}
	}
	err := manager.CreateScheduled(stufy.RequestScheduled{
		Systems:     c.Systems.StringSlice(),
		Description: c.Description,
		Title:       c.Title,
		Date:        c.Date,
		Duration:    c.Duration,
	})
	if err != nil {
		return err
	}
	messages.Println("Scheduled task has been", messages.C.Green("added"))
	return nil
}

func (c *NewScheduled) ExecuteSurvey() error {
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
