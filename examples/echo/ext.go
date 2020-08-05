package main

import "fmt"

// Extension type
type Extension string

// Run executes an action
func (e Extension) Run(action string) {
	switch action {
	case "greet":
		fmt.Println("Hello Universe")
	case "goodbye":
		fmt.Println("Goodbye Universe")
	default:
		fmt.Println("...")
	}
}

// Echo is exported
var Echo Extension

// package echo
