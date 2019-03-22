package main

import (
	"flag"
	"github.com/ArthurHlt/stufy"
	"github.com/ArthurHlt/stufy/messages"
	"github.com/ArthurHlt/stufy/model"
	"github.com/jessevdk/go-flags"
	"strings"
)

type SeverityFlag struct {
	Severity string `short:"e" long:"severity" description:"Define severity of the incident" choice:"major-outage" choice:"partial-outage" choice:"degraded-performance" choice:"under-maintenance"`
}

type InlineFlag struct {
	Inline  bool    `short:"i" long:"inline" description:"Inline request by using flag instead of survey"`
	Systems Systems `short:"s" long:"system" description:"Set systems affected"`
}

type Systems []System

func (ss Systems) StringSlice() []string {
	sString := make([]string, 0)
	for _, s := range ss {
		sString = append(sString, string(s))
	}
	return sString
}

type System string

func (System) Complete(match string) []flags.Completion {
	tPtr := flag.String("t", "", "")
	tLPtr := flag.String("target", "", "")
	flag.Parse()
	t := *tPtr
	if t == "" {
		t = *tLPtr
	}
	var err error
	manager, err = stufy.NewManager(t)
	if err != nil {
		panic(err)
	}
	config, err := manager.Config()
	if err != nil {
		panic(err)
	}
	items := make([]flags.Completion, 0)
	for _, s := range config.Content.Systems {
		if !strings.HasPrefix(s, match) {
			continue
		}
		items = append(items, flags.Completion{
			Item: s,
		})
	}
	return items
}

type FilenameIncident string

func (FilenameIncident) Complete(match string) []flags.Completion {
	return completeFilename(match, false)
}

type FilenameScheduled string

func (FilenameScheduled) Complete(match string) []flags.Completion {
	return completeFilename(match, true)
}

func completeFilename(match string, scheduleds bool) []flags.Completion {
	tPtr := flag.String("t", "", "")
	tLPtr := flag.String("target", "", "")
	tAllPtr := flag.Bool("all", false, "")
	flag.Parse()
	t := *tPtr
	if t == "" {
		t = *tLPtr
	}
	var err error
	manager, err = stufy.NewManager(t)
	if err != nil {
		panic(err)
	}
	messages.SetStopShow(true)
	var incidents model.Incidents
	if !scheduleds {
		incidents, err = manager.ListIncident(*tAllPtr)
	} else {
		incidents, err = manager.ListScheduled(*tAllPtr)
	}
	messages.SetStopShow(false)
	if err != nil {
		panic(err)
	}
	items := make([]flags.Completion, 0)
	for _, i := range incidents {
		filename := i.Filename()
		if strings.HasPrefix(filename, match) {
			items = append(items, flags.Completion{
				Item: filename,
			})
		}
		splitFilename := strings.SplitN(filename, "-", 2)
		if len(splitFilename) != 2 || strings.HasPrefix(splitFilename[1], match) {
			continue
		}

		items = append(items, flags.Completion{
			Item: filename,
		})
	}
	return items
}
