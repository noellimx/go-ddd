// Code assisted by GitHub Copilot CLI agent.

package main

import "fmt"

// Greet returns a greeting message for the given name.
func Greet(name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}

// main is the entry point of the hello world program.
func main() {
	fmt.Println(Greet("World"))
}
