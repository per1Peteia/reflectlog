package agent

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadFile(t *testing.T) {
	t.Run("returns error when outside permitted working directory", func(t *testing.T) {
		in := json.RawMessage(`{"path": "../../../../agent-python/main.py"}`)
		got, err := ReadFile(in)
		want := ""

		assert.Error(t, err)
		assert.ErrorAs(t, err, &PathNotPermittedError{})
		assert.Contains(t, err.Error(), "../../../../agent-python/main.py")
		assert.Equal(t, want, got)
	})

	t.Run("returns no error when happy", func(t *testing.T) {
		in := json.RawMessage(`{"path": "./internal/agent/test.txt"}`)
		got, err := ReadFile(in)
		want := "hello, world\n"

		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("returns error when trying to access a directory with same prefix as wd", func(t *testing.T) {
		in := json.RawMessage(`{"path": "../../../../code-editing-agent-evil/test.txt"}`)
		got, err := ReadFile(in)
		want := ""

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "../../../../code-editing-agent-evil/test.txt")
		assert.ErrorAs(t, err, &PathNotPermittedError{})
		assert.Equal(t, want, got)
	})
}

func TestListFilesInput(t *testing.T) {
	t.Run("returns error when outside of permitted working directory", func(t *testing.T) {
		in := json.RawMessage(`{"path": "../../../../agent-python/"}`)
		got, err := ListFiles(in)
		want := ""

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "../../../../agent-python/")
		assert.ErrorAs(t, err, &PathNotPermittedError{})
		assert.Equal(t, want, got)
	})
}

func TestEditFile(t *testing.T) {
	t.Run("returns error when outside of permitted working directory", func(t *testing.T) {
		in := json.RawMessage(`{"path": "../../../../code-editing-agent-evil/test.txt", "old_str": "hello, world\n", "new_str": "hello, you\n"}`)
		got, err := EditFile(in)
		want := ""

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "../../../../code-editing-agent-evil/test.txt")
		assert.ErrorAs(t, err, &PathNotPermittedError{})
		assert.Equal(t, want, got)
	})
}
