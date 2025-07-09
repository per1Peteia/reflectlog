package cli

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/mitchellh/cli"
	"github.com/per1Peteia/rfl/internal/agent"
)

func defaultCommandFactory() (cli.Command, error) {
	return &defaultCommand{}, nil
}

type defaultCommand struct{}

func (c *defaultCommand) Help() string {
	return rflHelp
}

func (c *defaultCommand) Synopsis() string {
	return "open a TUI to browse different materials"
}

func (c *defaultCommand) Run(args []string) int {
	return 0
}

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

	client := anthropic.NewClient()
	message := anthropic.NewUserMessage(anthropic.NewTextBlock(cbContent))
	synopsis, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		System:    []anthropic.TextBlockParam{{Text: clipSynopsisPrompt}},
		Model:     anthropic.ModelClaude4Sonnet20250514,
		Messages:  []anthropic.MessageParam{message},
		MaxTokens: int64(1024),
	})

	if err != nil {
		log.Fatal(err)
		return 1
	}

	for _, content := range synopsis.Content {
		switch content.Type {
		case "text":
			fmt.Println(content.Text)
		}
	}

	return 0
}

func InitCLI() int {
	c := cli.NewCLI("app", "1.0.0")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"": defaultCommandFactory,
		// "init": initCommandFactory,
		"chat": chatCommandFactory,
		// "log":  logCommandFactory,
		"clip": clipCommandFactory,
	}

	exitStatus, err := c.Run()
	if err != nil {
		return 1
	}

	return exitStatus
}
