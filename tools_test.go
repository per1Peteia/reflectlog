package main

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestReadFile(t *testing.T) {
	t.Run("returns error when outside permitted working directory", func(t *testing.T) {
		in := json.RawMessage(`{"path": "../agent-python/main.py"}`)
		got, err := ReadFile(in)
		want := ""

		var gotErr PathNotPermittedError
		isPermErr := errors.As(err, &gotErr)
		if !isPermErr {
			t.Fatalf("was not PathNotPermittedError, is %T", err)
		}

		wantErr := PathNotPermittedError{Path: "../agent-python/main.py"}

		if gotErr != wantErr {
			t.Errorf("got %v, want %v", gotErr, wantErr)
		}

		AssertError(t, err)
		AssertNilString(t, got, want)
	})

	t.Run("returns no error when happy", func(t *testing.T) {
		in := json.RawMessage(`{"path": "./test.txt"}`)
		got, err := ReadFile(in)
		want := "hello, world\n"

		AssertNoError(t, err)
		AssertEqualString(t, got, want)
	})

	t.Run("returns error when trying to access a directory that with same prefix as wd", func(t *testing.T) {
		in := json.RawMessage(`{"path": "../code-editing-agent-evil/test.txt"}`)
		got, err := ReadFile(in)
		want := ""

		AssertError(t, err)
		AssertNilString(t, got, want)
	})
}
