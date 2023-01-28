package naivebayes_test

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/jaimeteb/chatto/internal/clf"
	"github.com/jaimeteb/chatto/internal/clf/naivebayes"
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

func TestSaveAndLoad(t *testing.T) {
	classifConfig, err := clf.LoadConfig("../../"+testutils.Examples00TestPath, nil)
	if err != nil {
		t.Fatal(err)
	}

	classifOld := naivebayes.NewClassifier(map[string]interface{}{"tfidf": true})
	classifOld.Learn(classifConfig.Classification, &classifConfig.Pipeline)
	if err = classifOld.Save("./"); err != nil {
		t.Error(err)
	}

	classifNew, err := naivebayes.Load("./")
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(classifOld, classifNew) {
		t.Errorf("Classifier saved %v, Classifier loaded %v", spew.Sprint(classifOld), spew.Sprint(classifNew))
	}

	t.Cleanup(func() {
		testutils.RemoveFiles("gob")
	})
}
