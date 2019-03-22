package main

import (
	"fmt"
	"github.com/ArthurHlt/stufy"
	"github.com/ArthurHlt/stufy/messages"
	"github.com/ArthurHlt/stufy/model"
	"github.com/jessevdk/go-flags"
	"gopkg.in/AlecAivazis/survey.v1"
)

type UpdateIncident struct {
	InlineFlag
	SeverityFlag
	All           bool             `long:"all" description:"Show all incidents (by default it doesn't show resolved incident)'"`
	Resolved      bool             `short:"r" long:"resolved" description:"Set has resolved"`
	Filename      FilenameIncident `short:"f" long:"filename" description:"Set filename associated to update"`
	UpdateType    string           `short:"y" long:"up-type" description:"Update type" choice:"resolved" choice:"monitoring" choice:"status"`
	UpdateContent string           `short:"c" long:"up-content" description:"Update content"`
}

var updateIncident UpdateIncident

func (c *UpdateIncident) Execute(_ []string) error {
	if c.Inline {
		return c.ExecuteInline()
	}
	return c.ExecuteSurvey()
}

func (c *UpdateIncident) ExecuteInline() error {
	if c.Filename == "" {
		return &flags.Error{
			Type:    flags.ErrRequired,
			Message: "Filename flag required (--filename, -f)",
		}
	}
	err := manager.UpdateIncident(stufy.RequestUpdate{
		Filename:      string(c.Filename),
		Severity:      c.Severity,
		Systems:       c.Systems.StringSlice(),
		Resolved:      c.Resolved,
		UpdateType:    c.UpdateType,
		UpdateContent: c.UpdateContent,
		Confirm:       true,
	})
	if err != nil {
		return err
	}
	messages.Printfln("Incident %s has been %s", messages.C.Cyan(c.Filename), messages.C.Brown("updated"))
	return nil
}

func (c *UpdateIncident) ExecuteSurvey() error {
	config, err := manager.Config()
	if err != nil {
		return err
	}
	incidents, err := manager.ListIncident(c.All)
	if err != nil {
		return err
	}

	mType := " non resolved "
	if c.All {
		mType = ""
	}
	if len(incidents) == 0 {
		fmt.Printf("There is no %s incidents\n", mType)
		return nil
	}

	filename := ""
	err = survey.AskOne(&survey.Select{
		Message: "What incident do you want to update?",
		Options: incidents.Filenames(),
	}, &filename, survey.Required)
	if err != nil {
		return err
	}

	currentIncident, err := manager.FindIncident(filename)
	if err != nil {
		return err
	}

	qs := []*survey.Question{
		{
			Name: "Resolved",
			Prompt: &survey.Confirm{
				Message: "The incident has been resolved?",
				Default: false,
			},
		},
		{
			Name: "Severity",
			Prompt: &survey.Select{
				Message: "What is the severity of the incident?",
				Options: model.SeveritiesString(),
				Default: string(currentIncident.Severity),
			},
			Validate: survey.Required,
		},
		{
			Name: "Systems",
			Prompt: &survey.MultiSelect{
				Message: "What are the affected systems?",
				Options: config.Content.Systems,
				Default: currentIncident.AffectedSystems,
			},
			Validate: survey.Required,
		},
		{
			Name: "UpdateType",
			Prompt: &survey.Select{
				Message: "Choose a type if you want to add update container.",
				Options: []string{"resolved", "monitoring", "status", "no"},
				Default: "no",
			},
		},
		{
			Name:   "UpdateContent",
			Prompt: &survey.Input{Message: "Set an update content if you want to add update container."},
		},
		{
			Name: "Confirm",
			Prompt: &survey.Confirm{
				Message: "Are you sure you want to update the incident?",
				Default: true,
			},
		},
		{
			Name: "Open",
			Prompt: &survey.Confirm{
				Message: "Open the incident for editing?",
				Default: false,
			},
		},
	}

	var req stufy.RequestUpdate
	err = survey.Ask(qs, &req)
	if err != nil {
		return err
	}
	req.Filename = filename
	return manager.UpdateIncident(req)
}

func init() {
	desc := `Update an existing incident`
	c, err := parser.AddCommand(
		"update-incident",
		desc,
		desc,
		&updateIncident)
	if err != nil {
		panic(err)
	}
	c.Aliases = []string{"u"}
}
