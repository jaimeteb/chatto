package fsm

import (
	"testing"
)

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

func TestFSM(t *testing.T) {
	path := "../examples/00_test/"
	domain := Create(&path)
	commandList := []string{"turn_on", "turn_off", "hello_universe"}
	if !testEq(domain.CommandList, commandList) {
		t.Errorf("domain.CommandList is incorrect, got: %v, want: %v.", domain.CommandList, commandList)
	}

	machine := FSM{State: 0}
	extension, err := LoadExtension(&path)
	if err != nil {
		t.Errorf(err.Error())
	}

	resp1 := machine.ExecuteCmd("turn_on", "turn_on", domain, extension)
	if resp1 != "Turning on." {
		t.Errorf("resp is incorrect, got: %v, want: %v.", resp1, "Turning on.")
	}

	resp2 := machine.ExecuteCmd("turn_on", "turn_on", domain, extension)
	if resp2 != "Can't do that." {
		t.Errorf("resp is incorrect, got: %v, want: %v.", resp2, "Can't do that.")
	}

	resp3 := machine.ExecuteCmd("hello_universe", "hello", domain, extension)
	if resp3 != "Hello Universe" {
		t.Errorf("resp is incorrect, got: %v, want: %v.", resp3, "Hello Universe")
	}

	// testFuncExists := extension.GetFunc("ext_any")
	// testFuncExistsNot := extension.GetFunc("ext_any_other")
	// if testFuncExists == nil {
	// 	t.Errorf("GetFunc is incorrect, 'ext_any' should exist.")
	// }
	// if testFuncExistsNot != nil {
	// 	t.Errorf("GetFunc is incorrect, 'ext_any_other' should not exist.")
	// }
}
