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
	GetAllFuncs() []string
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
	loadExtErr := func(err error) Extension {
		log.Println("Error while loading extensions: ", err.Error())
		log.Println("Using bot without extensions.")
		return nil
	}

	if err := BuildPlugin(path); err != nil {
		return loadExtErr(err)
	}

	plug, err := plugin.Open(filepath.Join(*path, "ext.so"))
	if err != nil {
		return loadExtErr(err)
	}

	echo, err := plug.Lookup("Ext")
	if err != nil {
		return loadExtErr(err)
	}

	var extension Extension
	extension, ok := echo.(Extension)
	if !ok {
		return loadExtErr(err)
	}

	log.Println("Loaded extensions for FSM:")
	for i, f := range extension.GetAllFuncs() {
		log.Printf("%v\t%v\n", i, f)
	}

	return extension
}
