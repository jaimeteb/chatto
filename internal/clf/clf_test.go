package clf_test

import (
	"path"
	"reflect"
	"testing"

	"github.com/jaimeteb/chatto/internal/clf"
	"github.com/jaimeteb/chatto/internal/clf/dataset"
	"github.com/jaimeteb/chatto/internal/clf/knn"
	"github.com/jaimeteb/chatto/internal/clf/naivebayes"
	"github.com/jaimeteb/chatto/internal/clf/wordvectors"
	"github.com/jaimeteb/chatto/internal/testutils"
)

func TestNew(t *testing.T) {
	dsOnOff := dataset.DataSet{
		dataset.DataClass{
			Command: "on",
			Texts:   []string{"on", "turn_on"},
		},
		dataset.DataClass{
			Command: "off",
			Texts:   []string{"off", "turn_off"},
		},
	}

	type args struct {
		cfg *clf.Config
	}
	tests := []struct {
		name    string
		args    args
		want    reflect.Type
		wantErr bool
	}{
		{
			name: "naive bayes save",
			args: args{
				cfg: &clf.Config{
					Classification: dsOnOff,
					Model: clf.ModelConfig{
						Classifier: "naive_bayes",
						Parameters: map[string]interface{}{"tfidf": false},
						Directory:  "./test_model_naive_bayes",
						Save:       true,
						Load:       false,
					},
				},
			},
			want: reflect.TypeOf(&naivebayes.Classifier{}),
		},
		{
			name: "naive bayes load",
			args: args{
				cfg: &clf.Config{
					Classification: dsOnOff,
					Model: clf.ModelConfig{
						Classifier: "naive_bayes",
						Parameters: map[string]interface{}{"tfidf": false},
						Directory:  "./test_model_naive_bayes",
						Save:       false,
						Load:       true,
					},
				},
			},
			want: reflect.TypeOf(&naivebayes.Classifier{}),
		},
		{
			name: "knn save",
			args: args{
				cfg: &clf.Config{
					Classification: dsOnOff,
					Model: clf.ModelConfig{
						Classifier: "knn",
						Parameters: map[string]interface{}{"k": 1},
						Directory:  "./test_model_knn",
						Save:       true,
						Load:       false,
						WordVectorsConfig: wordvectors.Config{
							WordVectorsFile: path.Join("../", testutils.TestWordVectors),
							Truncate:        1.0,
							SkipOOV:         false,
						},
					},
				},
			},
			want: reflect.TypeOf(&knn.Classifier{}),
		},
		{
			name: "knn load",
			args: args{
				cfg: &clf.Config{
					Classification: dsOnOff,
					Model: clf.ModelConfig{
						Classifier: "knn",
						Parameters: map[string]interface{}{"k": 1},
						Directory:  "./test_model_knn",
						Save:       false,
						Load:       true,
						WordVectorsConfig: wordvectors.Config{
							WordVectorsFile: path.Join("../", testutils.TestWordVectors),
							Truncate:        1.0,
							SkipOOV:         false,
						},
					},
				},
			},
			want: reflect.TypeOf(&knn.Classifier{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := reflect.TypeOf(clf.New(tt.args.cfg).Model)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
	t.Cleanup(func() {
		testutils.RemoveFiles("gob")
	})
}
