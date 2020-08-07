package ext

import (
	"log"
	"os/exec"
	"path/filepath"
	"plugin"

	"github.com/jaimeteb/chatto/fsm"
)

// Extension interface
type Extension interface {
	GetFunc(string) func(*fsm.FSM)
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

// Create creates an extension
func Create(path *string) Extension {
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
	// greetFunc(&fsm.FSM{})
	// goodbyeFunc := extension.GetFunc("goodbye")
	// goodbyeFunc(&fsm.FSM{})
	// notaFunc := extension.GetFunc("nota")
	// notaFunc(&fsm.FSM{})

	return extension
}
