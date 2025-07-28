package cli

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mitchellh/cli"
)

type Clips struct {
	Items []Clip `json:"items"`
}

type Clip struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	ClipText  string `json:"clip_text"`
	ClipBrief string `json:"clip_brief"`
	ClipTitle string `json:"clip_title"`
}

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
	client := http.DefaultClient

	// process clip data clientside
	res, err := client.Get("http://localhost:8080/api/clips")
	if err != nil {
		// TODO
		return 1
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		// TODO
		return 1
	}

	clips := Clips{}
	err = json.Unmarshal(data, &clips)
	if err != nil {
		// TODO
		return 1
	}

	// get items to display as charm list

	return 0
}
