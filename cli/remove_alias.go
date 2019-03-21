package main

import (
	"github.com/ArthurHlt/stufy/messages"
	"strings"
)

type RemoveAlias struct {
}

var removeAlias RemoveAlias

func (c *RemoveAlias) Execute(args []string) error {
	if len(args) == 0 {
		messages.Fatal("You must provide an alias name as first argument")
	}
	return manager.RemoveAlias(strings.Join(args, "-"))
}

func init() {
	desc := `Remove an alias`
	c, err := parser.AddCommand(
		"remove-alias",
		desc,
		desc,
		&removeAlias)
	if err != nil {
		panic(err)
	}
	c.Aliases = []string{"ra"}
}
