package main

type Resync struct {
}

var resync Resync

func (c *Resync) Execute(args []string) error {
	return manager.Resync()
}

func init() {
	desc := `Resynchronize your target (useful when merging issue on git repo)`
	c, err := parser.AddCommand(
		"resync",
		desc,
		desc,
		&resync)
	if err != nil {
		panic(err)
	}
	c.Aliases = []string{"r"}
}
