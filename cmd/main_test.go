package main

import "testing"

func TestRun(t *testing.T) {
	got := run()
	want := 0
	if got != want {
		t.Errorf("unexpected failure, got %q, want %q", got, want)
	}
}
