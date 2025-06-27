package main

import (
	"testing"
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

