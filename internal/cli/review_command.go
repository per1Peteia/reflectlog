package cli

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mitchellh/cli"
)

func reviewCommandFactory() (cli.Command, error) {
	return &reviewCommand{}, nil
}

type reviewCommand struct{}

func (c *reviewCommand) Help() string {
	return ""
}

func (c *reviewCommand) Synopsis() string {
	return "review by browsing, searching, and editing clips in your clip library"
}

func (c *reviewCommand) Run(args []string) int {
	items :=  
}
