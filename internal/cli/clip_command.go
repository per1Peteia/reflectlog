package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
		fmt.Fprintf(os.Stderr, "Error running pbpaste: %s\n", err)
		return 1
	}

	text := strings.TrimSpace(string(out))

	// get LLM generated title and description for text
	payload, err := processPayloadConcurrently(text)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing clipped text: %s\n", err)
		return 1
	}

	body, err := json.Marshal(payload)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %s\n", err)
		return 1
	}

	client := &http.Client{Timeout: TIMEOUT}
	res, err := client.Post(BASE_URL+CREATE_CLIP_URL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error posting JSON: %v\n", err)
		return 1
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		data, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading response body: %v\n", err)
			return 1
		}
		fmt.Fprintf(os.Stderr, "Server Error (%d): %s\n", res.StatusCode, string(data))
		return 1 // this could be more accurate
	}

	return 0
}

type ClipParams struct {
	ClipText  string `json:"clip_text"`
	ClipBrief string `json:"clip_brief"`
	ClipTitle string `json:"clip_title"`
}

func processPayloadConcurrently(text string) (ClipParams, error) {
	type Result struct {
		val string
		err error
	}

	dCh := make(chan Result, 1)
	tCh := make(chan Result, 1)

	go func() {
		desc, err := getAnswerByPrompt(text, clipSynopsisPrompt)
		dCh <- Result{val: desc, err: err}
	}()
	go func() {
		title, err := getAnswerByPrompt(text, clipTitlePrompt)
		tCh <- Result{val: title, err: err}
	}()

	d := <-dCh
	if d.err != nil {
		return ClipParams{}, d.err
	}
	t := <-tCh
	if t.err != nil {
		return ClipParams{}, t.err
	}

	return ClipParams{
		ClipText:  text,
		ClipBrief: d.val,
		ClipTitle: t.val,
	}, nil
}
