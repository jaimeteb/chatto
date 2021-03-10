package knn

import (
	"math"
	"sort"
)

// Dist calculates euclidean distance betwee two slices
func Dist(source, dest []float64) float64 {
	val := 0.0
	for i := range source {
		val += math.Pow(source[i]-dest[i], 2)
	}
	return math.Sqrt(val)
}

// Slice argument sort
type Slice struct {
	sort.Interface
	idx []int
}

// Swap swaps
func (s Slice) Swap(i, j int) {
	s.Interface.Swap(i, j)
	s.idx[i], s.idx[j] = s.idx[j], s.idx[i]
}

// NewSlice creates a new Slice
func NewSlice(n sort.Interface) *Slice {
	s := &Slice{Interface: n, idx: make([]int, n.Len())}
	for i := range s.idx {
		s.idx[i] = i
	}
	return s
}

// NewFloat64Slice creates a new slice of float64
func NewFloat64Slice(n []float64) *Slice {
	return NewSlice(sort.Float64Slice(n))
}

// Entry map sort
type Entry struct {
	name  string
	value int
}

// List of entries
type List []Entry

func (l List) Len() int {
	return len(l)
}

func (l List) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l List) Less(i, j int) bool {
	if l[i].value == l[j].value {
		return l[i].name < l[j].name
	}
	return l[i].value > l[j].value
}

// Counter item frequence in slice
func Counter(target []string) map[string]int {
	counter := map[string]int{}
	for _, elem := range target {
		counter[elem]++
	}
	return counter
}

// KNN main structure
type KNN struct {
	k      int
	data   [][]float64
	labels []string
}

func (knn *KNN) fit(X [][]float64, Y []string) {
	//read data
	knn.data = X
	knn.labels = Y
}

func (knn *KNN) predict(X [][]float64) []string {

	predictedLabel := []string{}
	for _, source := range X {
		var (
			distList   []float64
			nearLabels []string
		)
		//calculate distance between predict target data and surpervised data
		for _, dest := range knn.data {
			distList = append(distList, Dist(source, dest))
		}
		//take top k nearest item's index
		s := NewFloat64Slice(distList)
		sort.Sort(s)
		targetIndex := s.idx[:knn.k]

		//get the index's label
		for _, ind := range targetIndex {
			nearLabels = append(nearLabels, knn.labels[ind])
		}

		//get label frequency
		labelFreq := Counter(nearLabels)

		//the most frequent label is the predict target label
		a := List{}
		for k, v := range labelFreq {
			e := Entry{k, v}
			a = append(a, e)
		}
		sort.Sort(a)
		predictedLabel = append(predictedLabel, a[0].name)
	}
	return predictedLabel

}
