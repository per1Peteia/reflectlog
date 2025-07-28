package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/mitchellh/cli"
	"github.com/per1Peteia/rfl/internal/agent"
)

func chatCommandFactory() (cli.Command, error) {
	return &chatCommand{}, nil
}

type chatCommand struct{}

func (c *chatCommand) Help() string {
	return chatHelp
}

func (c *chatCommand) Synopsis() string {
	return "chat with the agent"
}

func (c *chatCommand) Run(args []string) int {
	client := anthropic.NewClient()

	scanner := bufio.NewScanner(os.Stdin)
	getUserMessage := func() (string, bool) {
		if !scanner.Scan() {
			return "", false
		}
		return scanner.Text(), true
	}

	tools := []agent.ToolDefinition{agent.ReadFileDefinition, agent.ListFilesDefinition, agent.EditFileDefinition}
	a := agent.NewAgent(&client, getUserMessage, tools)

	err := a.Run(context.TODO())
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return 1
	}

	return 0
}
