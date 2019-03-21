package main

import (
	"fmt"
	"github.com/ArthurHlt/stufy"
	"github.com/ArthurHlt/stufy/messages"
	"github.com/jessevdk/go-flags"
	"os"
	"strings"
)

var Version string

type Options struct {
	Target  Target `short:"t" long:"target" description:"Set a target, this can be a directory path or a git repo (e.g.: git@github.com:ArthurHlt/stufy-test.git or https://user:password@github.com/ArthurHlt/stufy-test.git)"`
	Version func() `short:"v" long:"version" description:"Show version"`
}

type Target func(string)

func (Target) Complete(match string) []flags.Completion {
	items := make([]flags.Completion, 0)
	aliases := manager.Aliases()
	for _, a := range aliases {
		if !strings.HasPrefix(a.Name, match) {
			continue
		}
		items = append(items, flags.Completion{
			Item: a.Name,
		})
	}
	return items
}

var options Options

var manager *stufy.Manager

var parser = flags.NewParser(&options, flags.HelpFlag|flags.PassDoubleDash|flags.IgnoreUnknown)

func Parse(args []string) error {
	options.Target = func(target string) {
		var err error
		manager, err = stufy.NewManager(target)
		if err != nil {
			messages.Fatal(err.Error())
		}
	}
	askVersion := false
	options.Version = func() {
		askVersion = true
		fmt.Println("Stufy " + Version)
	}
	_, err := parser.ParseArgs(args[1:])

	if err != nil {
		if errFlag, ok := err.(*flags.Error); ok && askVersion && errFlag.Type == flags.ErrCommandRequired {
			return nil
		}
		if errFlag, ok := err.(*flags.Error); ok && errFlag.Type == flags.ErrHelp {
			fmt.Println(err.Error())
			return nil
		}
		return err
	}

	return nil
}

func main() {
	var err error
	manager, err = stufy.NewManager("")
	if err != nil {
		messages.Fatal(err.Error())
	}
	err = Parse(os.Args)
	if err != nil {
		messages.Fatal(err.Error())
	}
}
