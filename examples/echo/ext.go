package main

import "fmt"

// FuncMap maps function names to functions
type FuncMap map[string]func()

func greetFunc() {
	fmt.Println("Hello Universe")
}

func goodbyeFunc() {
	fmt.Println("Goodbye Universe")
}

// Echo is exported
var Echo = FuncMap{
	"greet":   greetFunc,
	"goodbye": goodbyeFunc,
}

// Run executes an action
func (fm FuncMap) Run(action string) {
	if _, ok := fm[action]; ok {
		fm[action]()
	} else {
		fmt.Println("...")
	}
}

func main() {}
