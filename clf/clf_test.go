package clf

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/navossoc/bayesian"
)

var clfYaml = `
classification:
  - command: "turn_on"
    texts:
      - "turn on"
      - "on"

  - command: "turn_off"
    texts:
      - "turn off"
      - "off"
`

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

func writeDummyFile() error {
	clfFile := []byte(clfYaml)
	return ioutil.WriteFile("clf.yml", clfFile, 0644)
}

func removeDummyFile() error {
	return os.Remove("clf.yml")
}

func TestClf(t *testing.T) {
	if err := writeDummyFile(); err != nil {
		t.Errorf(err.Error())
	}

	classif := Create()
	classes := []bayesian.Class{"turn_on", "turn_off"}
	if !testEq(classes, classif.Classes) {
		t.Errorf("classes is incorrect, got: %v, want: %v.", classes, classif.Classes)
	}

	pred, _ := classif.Predict("on")
	if pred != "turn_on" {
		t.Errorf("pred is incorrect, got: %v, want: %v.", pred, "turn_on")
	}

	if err := removeDummyFile(); err != nil {
		t.Errorf(err.Error())
	}
}
