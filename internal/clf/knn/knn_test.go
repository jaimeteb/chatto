package knn_test

import (
	"path"
	"testing"

	"github.com/jaimeteb/chatto/internal/clf"
	"github.com/jaimeteb/chatto/internal/clf/knn"
	"github.com/jaimeteb/chatto/internal/clf/wordvectors"
	"github.com/jaimeteb/chatto/internal/testutils"
)

func TestClassifier(t *testing.T) {
	reload := make(chan clf.Config)
	cfg, err := clf.LoadConfig(path.Join("../../", testutils.Examples01MoodbotPath), reload)
	if err != nil {
		t.Fatalf("failed to load clf config: %v", err)
	}

	type args struct {
		wv     wordvectors.Config
		params map[string]interface{}
		text   string
	}
	type want struct {
		accuracy   float32
		prediction string
	}
	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				wv: wordvectors.Config{
					WordVectorsFile: path.Join("../../", testutils.TestWordVectors),
					Truncate:        1.0,
				},
				params: map[string]interface{}{"k": 3},
				text:   "hello",
			},
			want: want{
				accuracy:   0.6,
				prediction: "greet",
			},
		},
		{
			name: "test",
			args: args{
				wv: wordvectors.Config{
					WordVectorsFile: path.Join("../../", testutils.TestWordVectors),
					Truncate:        0.1,
				},
				params: map[string]interface{}{"k": 1},
				text:   "hello",
			},
			want: want{
				accuracy:   0.3,
				prediction: "greet",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := knn.NewClassifier(tt.args.wv, tt.args.params)
			gotAccuracy := c.Learn(cfg.Classification, &cfg.Pipeline)
			gotPrediction, _ := c.Predict(tt.args.text, &cfg.Pipeline)
			if !(gotAccuracy >= tt.want.accuracy) {
				t.Errorf("Learn() = %v, want %v", gotAccuracy, tt.want.accuracy)
			}
			if !(gotPrediction >= tt.want.prediction) {
				t.Errorf("Predict() = %v, want %v", gotPrediction, tt.want.prediction)
			}
		})
	}
	t.Cleanup(testutils.RemoveGobFiles)
}
