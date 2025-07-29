package cli

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os/exec"
	"strings"
	"time"

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
		// TODO
		return 1
	}

	text := strings.TrimSpace(string(out))

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
		// TODO
		return 1
	}
	t := <-tCh
	if t.err != nil {
		//TODO
		return 1
	}

	payload := struct {
		ClipText  string `json:"clip_text"`
		ClipBrief string `json:"clip_brief"`
		ClipTitle string `json:"clip_title"`
	}{
		ClipText:  text,
		ClipBrief: d.val,
		ClipTitle: t.val,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		// TODO
		return 1
	}

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Post(BASE_URL+"/api/clips", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return 1
	}
	if res.StatusCode != 200 {
		//TODO
		return 1
	}
	defer res.Body.Close()

	return 0
}
