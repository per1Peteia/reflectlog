package agent_test

import (
	"encoding/json"
	agent "github.com/per1Peteia/code-editing-agent/internal/agent"
	"testing"
)

func TestReadFile(t *testing.T) {
	t.Run("returns error when outside permitted working directory", func(t *testing.T) {
		in := json.RawMessage(`{"path": "../../../../agent-python/main.py"}`)
		got, err := agent.ReadFile(in)
		want := ""
		wantErr := agent.PathNotPermittedError{Path: "../../../../agent-python/main.py"}

		AssertError(t, err)
		AssertErrorAs(t, wantErr, err)
		AssertNilString(t, got, want)
	})

	t.Run("returns no error when happy", func(t *testing.T) {
		in := json.RawMessage(`{"path": "./internal/agent/tests/test.txt"}`)
		got, err := agent.ReadFile(in)
		want := "hello, world\n"

		AssertNoError(t, err)
		AssertEqualString(t, got, want)
	})

	t.Run("returns error when trying to access a directory with same prefix as wd", func(t *testing.T) {
		in := json.RawMessage(`{"path": "../../../../code-editing-agent-evil/test.txt"}`)
		got, err := agent.ReadFile(in)
		want := ""
		wantErr := agent.PathNotPermittedError{Path: "../../../../code-editing-agent-evil/test.txt"}

		AssertError(t, err)
		AssertErrorAs(t, wantErr, err)
		AssertNilString(t, got, want)
	})
}

func TestListFilesInput(t *testing.T) {
	t.Run("returns error when outside of permitted working directory", func(t *testing.T) {
		in := json.RawMessage(`{"path": "../../../../agent-python/"}`)
		got, err := agent.ListFiles(in)
		want := ""
		wantErr := agent.PathNotPermittedError{Path: "../../../../agent-python/"}

		AssertError(t, err)
		AssertErrorAs(t, wantErr, err)
		AssertNilString(t, got, want)
	})
}

func TestEditFile(t *testing.T) {
	t.Run("returns error when outside of permitted working directory", func(t *testing.T) {
		in := json.RawMessage(`{"path": "../../../../code-editing-agent-evil/test.txt", "old_str": "hello, world\n", "new_str": "hello, you\n"}`)
		got, err := agent.EditFile(in)
		want := ""
		wantErr := agent.PathNotPermittedError{Path: "../../../../code-editing-agent-evil/test.txt"}

		AssertError(t, err)
		AssertErrorAs(t, wantErr, err)
		AssertNilString(t, got, want)

	})
}
