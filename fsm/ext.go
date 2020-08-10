package fsm

import (
	"log"
	"os/exec"
	"path/filepath"
	"plugin"
)

// Extension interface
type Extension interface {
	GetFunc(string) func(*FSM) interface{}
}

// BuildPlugin builds the extension code as a plugin
func BuildPlugin(path *string) error {
	buildGo := "go"
	buildArgs := []string{
		"build",
		"-buildmode=plugin",
		"-o",
		filepath.Join(*path, "ext.so"),
		filepath.Join(*path, "ext.go"),
	}

	cmd := exec.Command(buildGo, buildArgs...)
	stdout, err := cmd.Output()

	if err != nil {
		return err
	}

	log.Println(string(stdout))
	return nil
}

// LoadExtension creates an extension
func LoadExtension(path *string) Extension {
	if err := BuildPlugin(path); err != nil {
		log.Println(err)
		return nil
	}

	plug, err := plugin.Open(filepath.Join(*path, "ext.so"))
	if err != nil {
		log.Println(err)
		return nil
	}

	echo, err := plug.Lookup("Ext")
	if err != nil {
		log.Println(err)
		return nil
	}

	var extension Extension
	extension, ok := echo.(Extension)
	if !ok {
		log.Println("unexpected type from module symbol")
		return nil
	}

	// greetFunc := extension.GetFunc("greet")
	// greetFunc(&FSM{})
	// goodbyeFunc := extension.GetFunc("goodbye")
	// goodbyeFunc(&FSM{})
	// notaFunc := extension.GetFunc("nota")
	// notaFunc(&FSM{})

	return extension
}
