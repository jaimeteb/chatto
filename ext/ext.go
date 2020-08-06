package ext

import (
	"log"
	"os/exec"
	"path/filepath"
	"plugin"
)

// Extension interface
type Extension interface {
	Run(string)
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

	echo, err := plug.Lookup("Echo")
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

	return extension
}
