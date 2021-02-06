package clf

import (
	"testing"
)

func TestClf1(t *testing.T) {
	path := "../examples/00_test/"
	classif := Create(&path)
	pred1, _ := classif.Model.Predict("on", classif.Pipeline)
	if pred1 != "turn_on" {
		t.Errorf("pred is incorrect, got: %v, want: %v.", pred1, "turn_on")
	}
	pred2, _ := classif.Model.Predict("foo", classif.Pipeline)
	if pred2 != "" {
		t.Errorf("pred is incorrect, got: %v, want: %v.", pred2, "")
	}
}

func TestClf2(t *testing.T) {
	path := "../examples/404/"
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("No panic")
		}
	}()
	Create(&path)
}
