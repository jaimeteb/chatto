package naivebayes_test

import (
	"testing"

	"github.com/jaimeteb/chatto/internal/clf"
	"github.com/jaimeteb/chatto/internal/testutils"
)

func TestClfPredictions(t *testing.T) {
	classifReloadChan := make(chan clf.Config)
	classifConfig, err := clf.LoadConfig("../../"+testutils.Examples00TestPath, classifReloadChan)
	if err != nil {
		t.Fatal(err)
	}

	classif := clf.New(classifConfig)

	pred1, _ := classif.Model.Predict("on", &classifConfig.Pipeline)
	if pred1 != "turn_on" {
		t.Errorf("pred is incorrect, got: %v, want: %v.", pred1, "turn_on")
	}

	pred2, _ := classif.Model.Predict("foo", &classifConfig.Pipeline)
	if pred2 != "" {
		t.Errorf("pred is incorrect, got: %v, want: %v.", pred2, "")
	}
	t.Cleanup(func() {
		testutils.RemoveFiles("gob")
	})
}
