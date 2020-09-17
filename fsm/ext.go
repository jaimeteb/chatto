package fsm

import (
	"log"
	"os/exec"
	"path/filepath"
	"plugin"
)

// Extension interface
type Extension interface {
	GetFunc(string) func(*FSM, *Domain, string) interface{}
	GetAllFuncs() []string
}

// FuncMap maps function names to functions
type FuncMap map[string]func(*FSM, *Domain, string) interface{}

// GetFunc gets a function from the function map
func (fm FuncMap) GetFunc(action string) func(*FSM, *Domain, string) interface{} {
	if _, ok := fm[action]; ok {
		return fm[action]
	}
	return func(*FSM, *Domain, string) interface{} {
		return nil
	}
}

// GetAllFuncs retreives all functions in function map
func (fm FuncMap) GetAllFuncs() []string {
	allFuncs := make([]string, 0)
	for funcName := range fm {
		allFuncs = append(allFuncs, funcName)
	}
	return allFuncs
}

// BuildPlugin builds the extension code as a plugin
func BuildPlugin(path *string) error {
	buildGo := "go"
	buildArgs := []string{
		"build",
		"-buildmode=plugin",
		"-o",
		filepath.Join(*path, "ext/ext.so"),
		filepath.Join(*path, "ext/ext.go"),
	}

	cmd := exec.Command(buildGo, buildArgs...)
	_, err := cmd.Output()

	if err != nil {
		return err
	}

	return nil
}

// LoadExtension creates an extension
func LoadExtension(path *string) (Extension, error) {
	loadExtErr := func(err error) (Extension, error) {
		log.Println("Error while loading extensions: ", err.Error())
		log.Println("Using bot without extensions.")
		return nil, err
	}

	if err := BuildPlugin(path); err != nil {
		return loadExtErr(err)
	}

	plug, err := plugin.Open(filepath.Join(*path, "ext/ext.so"))
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

	return extension, nil
}
