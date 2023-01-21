package clf_test

import (
	"reflect"
	"testing"

	"github.com/jaimeteb/chatto/internal/clf"
	"github.com/jaimeteb/chatto/internal/clf/dataset"
)

func TestGetConfusionMatrix(t *testing.T) {
	type args struct {
		cfg *clf.Config
	}
	type want struct {
		confusionMatrix [][]int
		yTrue, yPred    []int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "on off confusion matrix",
			args: args{
				cfg: &clf.Config{
					Classification: dataset.DataSet{
						dataset.DataClass{
							Command: "on",
							Texts:   []string{"on", "turn_on"},
						},
						dataset.DataClass{
							Command: "off",
							Texts:   []string{"off", "turn_off"},
						},
					},
					Model: clf.ModelConfig{
						Classifier: "naive_bayes",
					},
				},
			},
			want: want{
				confusionMatrix: [][]int{{2, 0}, {0, 2}},
				yTrue:           []int{0, 0, 1, 1},
				yPred:           []int{0, 0, 1, 1},
			},
		},
		{
			name: "bad confusion matrix",
			args: args{
				cfg: &clf.Config{
					Classification: dataset.DataSet{
						dataset.DataClass{
							Command: "0",
							Texts:   []string{"a", "b", "c"},
						},
						dataset.DataClass{
							Command: "1",
							Texts:   []string{"a", "b", "d"},
						},
					},
					Model: clf.ModelConfig{
						Classifier: "naive_bayes",
					},
				},
			},
			want: want{
				confusionMatrix: [][]int{{3, 0}, {2, 1}},
				yTrue:           []int{0, 0, 0, 0, 0, 1},
				yPred:           []int{0, 0, 0, 1, 1, 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := want{}
			got.confusionMatrix, got.yPred, got.yTrue = clf.GetConfusionMatrix(tt.args.cfg)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSums(t *testing.T) {
	type args struct {
		confusionMatrix [][]int
	}
	type want struct {
		sumI, sumJ     []int
		sumIJ, sumTrue int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "confusion matrix 1",
			args: args{
				confusionMatrix: [][]int{
					{1, 0},
					{0, 1},
				},
			},
			want: want{
				sumI:    []int{1, 1},
				sumJ:    []int{1, 1},
				sumIJ:   2,
				sumTrue: 2,
			},
		},
		{
			name: "confusion matrix 2",
			args: args{
				confusionMatrix: [][]int{
					{5, 1, 0},
					{0, 5, 2},
					{0, 0, 5},
				},
			},
			want: want{
				sumI:    []int{6, 7, 5},
				sumJ:    []int{5, 6, 7},
				sumIJ:   18,
				sumTrue: 15,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := want{}
			got.sumI, got.sumJ, got.sumIJ, got.sumTrue = clf.GetSums(tt.args.confusionMatrix)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScores(t *testing.T) {
	type args struct {
		confusionMatrix [][]int
	}
	type want struct {
		s clf.Scores
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "confusion matrix 1",
			args: args{
				confusionMatrix: [][]int{
					{1, 0},
					{0, 1},
				},
			},
			want: want{
				s: clf.Scores{
					Precision: []float64{1.0, 1.0},
					Recall:    []float64{1.0, 1.0},
					F1score:   []float64{1.0, 1.0},

					PrecisionAvg: 1.0,
					RecallAvg:    1.0,
					F1scoreAvg:   1.0,

					PrecisionWeightedAvg: 1.0,
					RecallWeightedAvg:    1.0,
					F1scoreWeightedAvg:   1.0,

					Accuracy: 1.0,
				},
			},
		},
		{
			name: "confusion matrix 2",
			args: args{
				confusionMatrix: [][]int{
					{1, 1},
					{0, 2},
				},
			},
			want: want{
				s: clf.Scores{
					Precision: []float64{0.5, 1.0},
					Recall:    []float64{1.0, 0.6666666666666666},
					F1score:   []float64{0.6666666666666666, 0.8},

					PrecisionAvg: 0.75,
					RecallAvg:    0.8333333333333333,
					F1scoreAvg:   0.7333333333333334,

					PrecisionWeightedAvg: 0.75,
					RecallWeightedAvg:    0.8333333333333333,
					F1scoreWeightedAvg:   0.7333333333333334,

					Accuracy: 0.75,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sumI, sumJ, sumIJ, sumTrue := clf.GetSums(tt.args.confusionMatrix)
			got := clf.GetScores(tt.args.confusionMatrix, sumI, sumJ, sumIJ, sumTrue)
			if !reflect.DeepEqual(got, tt.want.s) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
