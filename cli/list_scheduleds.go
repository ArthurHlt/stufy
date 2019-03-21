package main

import (
	"encoding/json"
	"fmt"
	"github.com/ArthurHlt/stufy/messages"
	"github.com/olekukonko/tablewriter"
	"strings"
	"time"
)

type ListScheduleds struct {
	All  bool `short:"a" long:"all" description:"Show all incidents (by default it doesn't show resolved incident)"`
	Json bool `short:"j" long:"json" description:"Show as json"`
}

var listScheduleds ListScheduleds

func (c *ListScheduleds) Execute(_ []string) error {
	incidents, err := manager.ListScheduled(c.All)
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
	table.SetHeader([]string{"Title", "Last Update", "State", "Affected System", "When", "Duration", "Description"})
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
		state := messages.C.Brown("Scheduled").String()
		if i.Resolved {
			state = messages.C.Green("Finished").String()
		}
		if time.Time(*i.Scheduled).Before(time.Now()) &&
			time.Time(*i.Scheduled).Add(time.Duration(i.Duration) * time.Minute).After(time.Now()) {
			state = messages.C.Blue("In Progress").String()
		}
		row = append(row, state)
		row = append(row, strings.Join(i.AffectedSystems, ", "))
		row = append(row, time.Time(*i.Scheduled).Format(time.RFC822))
		row = append(row, fmt.Sprintf("%d min", i.Duration))
		row = append(row, i.Content)
		table.Append(row)
	}
	table.Render()
	return nil
}

func init() {
	desc := `List scheduleds`
	c, err := parser.AddCommand(
		"list-scheduleds",
		desc,
		desc,
		&listScheduleds)
	if err != nil {
		panic(err)
	}
	c.Aliases = []string{"ls"}
}
