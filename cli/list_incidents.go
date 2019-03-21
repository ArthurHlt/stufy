package main

import (
	"encoding/json"
	"fmt"
	"github.com/ArthurHlt/stufy/messages"
	"github.com/ArthurHlt/stufy/model"
	"github.com/olekukonko/tablewriter"
	"strings"
	"time"
)

type ListIncidents struct {
	All  bool `short:"a" long:"all" description:"Show all incidents (by default it doesn't show resolved incident)"`
	Json bool `short:"j" long:"json" description:"Show as json"`
}

var listIncidents ListIncidents

func (c *ListIncidents) Execute(_ []string) error {
	incidents, err := manager.ListIncident(c.All)
	if err != nil {
		return err
	}
	if c.Json {
		b, err := json.MarshalIndent(incidents, "", "\t")
		if err != nil {
			return err
		}
		fmt.Println(string(b))
		return nil
	}
	table := tablewriter.NewWriter(messages.Output())
	table.SetHeader([]string{"Title", "Last Update", "State", "Affected System", "Description"})
	table.SetRowSeparator("-")
	table.SetBorder(false)
	table.SetRowLine(true)
	for _, i := range incidents {
		row := make([]string, 0)
		row = append(row, i.Title)
		lastUpdate := i.Date
		if !i.Modified.IsZero() {
			lastUpdate = i.Modified
		}
		row = append(row, time.Time(lastUpdate).Format(time.RFC822))
		severity := messages.C.Brown(i.Severity).String()
		if i.Severity == model.MajorOutage {
			severity = messages.C.Red(i.Severity).String()
		}
		if i.Resolved {
			severity = messages.C.Green("Resolved").String()
		}
		row = append(row, severity)
		row = append(row, strings.Join(i.AffectedSystems, ", "))
		row = append(row, i.Content)
		table.Append(row)
	}
	table.Render()
	return nil
}

func init() {
	desc := `List incidents`
	c, err := parser.AddCommand(
		"list-incidents",
		desc,
		desc,
		&listIncidents)
	if err != nil {
		panic(err)
	}
	c.Aliases = []string{"li"}
}
