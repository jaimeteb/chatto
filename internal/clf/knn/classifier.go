package knn

import (
	"math"
	"sort"
)

// EuclideanDistance calculates euclidean distance between two points
func EuclideanDistance(p1, p2 []float64) float64 {
	val := 0.0
	for i := range p1 {
		val += math.Pow(p1[i]-p2[i], 2)
	}
	return math.Sqrt(val)
}

// KNN main structure
type KNN struct {
	K      int
	Data   [][]float64
	Labels []string
}

type neighbor struct {
	distance float64
	label    string
}

type neighborCount struct {
	label string
	count int
}

// PredictMany performs a classification on multiple input vectors
func (knn *KNN) PredictMany(X [][]float64) (predictedLabels []string, probabilities []float64) {
	for _, x := range X {
		pred, prob := knn.PredictOne(x)
		predictedLabels = append(predictedLabels, pred)
		probabilities = append(probabilities, prob)
	}
	return
}

// PredictOne performs a classification on one input vector
func (knn *KNN) PredictOne(X []float64) (predictedLabel string, probability float64) {
	neighs := make([]neighbor, len(knn.Data))

	for i := 0; i < len(knn.Data); i++ {
		neighs[i] = neighbor{
			distance: EuclideanDistance(X, knn.Data[i]),
			label:    knn.Labels[i],
		}
	}

	sort.SliceStable(neighs, func(i, j int) bool {
		return neighs[i].distance < neighs[j].distance
	})

	nearest := neighs[:knn.K]

	labelFreq := map[string]int{}
	for _, nn := range nearest {
		labelFreq[nn.label]++
	}

	labelSort := make([]neighborCount, 0)
	for label, count := range labelFreq {
		labelSort = append(labelSort, neighborCount{
			label: label,
			count: count,
		})
	}

	sort.SliceStable(labelSort, func(i, j int) bool {
		return labelSort[i].count > labelSort[j].count
	})

	predictedLabel = labelSort[0].label
	probability = float64(labelSort[0].count) / float64(knn.K)
	return predictedLabel, probability
}
