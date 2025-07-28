package cli

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/mitchellh/cli"
)

func clipCommandFactory() (cli.Command, error) {
	return &clipCommand{}, nil
}

type clipCommand struct{}

func (c *clipCommand) Help() string {
	return ""
}

func (c *clipCommand) Synopsis() string {
	return "clip some information you want to capture, saves it to the db to makes it processable in the future either programmatically or for the agent"
}

func (c *clipCommand) Run(args []string) int {
	out, err := exec.Command("pbpaste").Output()
	if err != nil {
		return 1
	}

	cbContent := strings.TrimSpace(string(out))
	fmt.Println(cbContent)

	description := getAnswer(cbContent)
	fmt.Println(description)

	return 0
}
