package cli

import (
	"os"

	"github.com/mitchellh/cli"
)

func defaultCommandFactory() (cli.Command, error) {
	return &defaultCommand{}, nil
}

type defaultCommand struct{}

func (c *defaultCommand) Help() string {
	return ""
}

func (c *defaultCommand) Synopsis() string {
	return "open a TUI to browse different materials"
}

func (c *defaultCommand) Run(args []string) int {
	return 0
}

func InitCLI() int {
	c := cli.NewCLI("app", "1.0.0")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"": defaultCommandFactory,
		// "init": initCommandFactoy,
		"chat": chatCommandFactory,
		// "log":  logCommandFactory,
		"clip": clipCommandFactory,
		// "review": reviewCommandFactory,
	}

	exitStatus, err := c.Run()
	if err != nil {
		return 1
	}

	return exitStatus
}
