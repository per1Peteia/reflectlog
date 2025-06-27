package agent_test

import (
	"errors"
	"testing"

	agent "github.com/per1Peteia/code-editing-agent/internal/agent"
)

func AssertError(t testing.TB, got error) {
	t.Helper()

	if got == nil {
		t.Fatalf(`expected error but go none`)
	}
}

func AssertNoError(t testing.TB, got error) {
	t.Helper()

	if got != nil {
		t.Fatalf(`expected no error but got error: %v`, got)
	}
}

func AssertErrorAs(t testing.TB, wantErr, err error) {
	t.Helper()

	var gotErr agent.PathNotPermittedError
	isPermErr := errors.As(err, &gotErr)
	if !isPermErr {
		t.Fatalf("was not PathNotPermittedError, is %T: %s", err, err)
	}

	if gotErr != wantErr {
		t.Errorf("got %v, want %v", gotErr, wantErr)
	}
}

func AssertNilString(t testing.TB, got, want string) {
	t.Helper()

	if got != want {
		t.Fatalf(`got %q, want %q`, got, want)
	}
}

func AssertEqualString(t testing.TB, got, want string) {
	t.Helper()

	if got != want {
		t.Fatalf(`got %q, want %q`, got, want)
	}
}
