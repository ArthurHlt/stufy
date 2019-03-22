package main

import (
	"github.com/ArthurHlt/stufy"
	"github.com/ArthurHlt/stufy/messages"
	"github.com/ArthurHlt/stufy/model"
	"github.com/jessevdk/go-flags"
	"gopkg.in/AlecAivazis/survey.v1"
)

type NewIncident struct {
	InlineFlag
	SeverityFlag
	Cause       string `short:"c" long:"cause" description:"Cause of the incident"`
	Description string `short:"d" long:"description" description:"Description of the incident"`
}

var newIncident NewIncident

func (c *NewIncident) Execute(_ []string) error {
	if c.Inline {
		return c.ExecuteInline()
	}
	return c.ExecuteSurvey()
}

func (c *NewIncident) ExecuteInline() error {
	if len(c.Systems) == 0 {
		return &flags.Error{
			Type:    flags.ErrRequired,
			Message: "Affected system flag required (--system, -s)",
		}
	}
	if c.Cause == "" {
		return &flags.Error{
			Type:    flags.ErrRequired,
			Message: "Cause flag required (--cause, -c)",
		}
	}
	if c.Severity == "" {
		return &flags.Error{
			Type:    flags.ErrRequired,
			Message: "Severity flag required (--severity, -e)",
		}
	}
	if c.Description == "" {
		return &flags.Error{
			Type:    flags.ErrRequired,
			Message: "Description flag, (--description, -d)",
		}
	}
	err := manager.CreateIncident(stufy.RequestCreate{
		Systems:     c.Systems.StringSlice(),
		Description: c.Description,
		Severity:    c.Severity,
		Cause:       c.Cause,
	})
	if err != nil {
		return err
	}
	messages.Println("Incident has been", messages.C.Green("added"))
	return nil
}

func (c *NewIncident) ExecuteSurvey() error {

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
