package fsm

import (
	"io/ioutil"
	"os"
	"testing"
)

var fsmYaml = `
states:
  - "off"
  - "on"
commands:
  - "turn_on"
  - "turn_off"
functions:
  - tuple:
      command: "turn_on"
      state: "off"
    transition: "on"
    message: "Turning on."
  - tuple:
      command: "turn_off"
      state: "on"
    transition: "off"
    message: "Turning off."
defaults:
  unknown: "Can't do that."
  unsure: "???"
`

func testEq(a, b []string) bool {
	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func writeDummyFile() error {
	fsmFile := []byte(fsmYaml)
	return ioutil.WriteFile("fsm.yml", fsmFile, 0644)
}

func removeDummyFile() error {
	return os.Remove("fsm.yml")
}

func TestFSM(t *testing.T) {
	if err := writeDummyFile(); err != nil {
		t.Errorf(err.Error())
	}

	here := "."
	domain := Create(&here)
	commandList := []string{"turn_on", "turn_off"}
	if !testEq(domain.CommandList, commandList) {
		t.Errorf("domain.CommandList is incorrect, got: %v, want: %v.", domain.CommandList, commandList)
	}

	machine := FSM{State: 0}
	resp1 := machine.ExecuteCmd("turn_on", "turn_on", domain, *new(Extension))
	if resp1 != "Turning on." {
		t.Errorf("resp is incorrect, got: %v, want: %v.", resp1, "Turning on.")
	}

	resp2 := machine.ExecuteCmd("turn_on", "turn_on", domain, *new(Extension))
	if resp2 != "Can't do that." {
		t.Errorf("resp is incorrect, got: %v, want: %v.", resp2, "Can't do that.")
	}

	if err := removeDummyFile(); err != nil {
		t.Errorf(err.Error())
	}
}
