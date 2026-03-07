// Code assisted by GitHub Copilot CLI agent.

package main

import "testing"

// TestGreet verifies that Greet returns the expected greeting message.
func TestGreet(t *testing.T) {
	got := Greet("World")
	want := "Hello, World!"
	if got != want {
		t.Errorf("Greet(%q) = %q, want %q", "World", got, want)
	}
}
