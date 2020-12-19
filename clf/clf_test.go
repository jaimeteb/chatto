package clf

import (
	"testing"

	"github.com/navossoc/bayesian"
)

func testEq(a, b []bayesian.Class) bool {
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

func testEqStr(a, b []string) bool {
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

func TestClf(t *testing.T) {
	path := "../examples/00_test/"
	classif := Create(&path)
	classes := []bayesian.Class{"turn_on", "turn_off", "hello_universe", "any"}
	if !testEq(classes, classif.Classes) {
		t.Errorf("classes is incorrect, got: %v, want: %v.", classes, classif.Classes)
	}

	pred, _ := classif.Predict("on")
	if pred != "turn_on" {
		t.Errorf("pred is incorrect, got: %v, want: %v.", pred, "turn_on")
	}
}

func TestPreprocess(t *testing.T) {
	testString := "fOo BaR"
	resultString := Pipeline(&testString, &PipelineConfig{true, true})
	expectedResult := []string{"foo", "bar"}

	if !testEqStr(resultString, expectedResult) {
		t.Errorf("resultString is incorrect, got: %v, want: %v.", resultString, expectedResult)
	}
}
