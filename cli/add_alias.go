package main

import (
	"github.com/ArthurHlt/stufy/messages"
	"strings"
)

type AddAlias struct {
}

var addAlias AddAlias

func (c *AddAlias) Execute(args []string) error {
	if len(args) == 0 {
		messages.Fatal("You must provide an alias name as first argument")
	}
	return manager.AddAlias(strings.Join(args, "-"))
}

func init() {
	desc := `Add an alias to your current target to use instead of plain target`
	c, err := parser.AddCommand(
		"add-alias",
		desc,
		desc,
		&addAlias)
	if err != nil {
		panic(err)
	}
	c.Aliases = []string{"a"}
}
