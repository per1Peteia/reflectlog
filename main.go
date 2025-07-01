package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/joho/godotenv"

	agent "github.com/per1Peteia/rfl/internal/agent"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}

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

	err = a.Run(context.TODO())
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
}
