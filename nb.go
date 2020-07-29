package main

import (
	"fmt"

	"github.com/navossoc/bayesian"
)

const (
	// Good class
	Good bayesian.Class = "Good"
	// Bad class
	Bad bayesian.Class = "Bad"
)

func main() {
	// Create a classifier with TF-IDF support.
	classifier := bayesian.NewClassifierTfIdf(Good, Bad)

	good1 := []string{"feel", "good"}
	good2 := []string{"i", "am", "great"}
	good3 := []string{"feel", "very", "happy"}

	bad1 := []string{"feel", "very", "sad"}
	bad2 := []string{"bad", "day"}
	bad3 := []string{"not", "good"}

	classifier.Learn(good1, Good)
	classifier.Learn(good2, Good)
	classifier.Learn(good3, Good)

	classifier.Learn(bad1, Bad)
	classifier.Learn(bad2, Bad)
	classifier.Learn(bad3, Bad)

	// Required
	classifier.ConvertTermsFreqToTfIdf()

	probs, likely, _ := classifier.ProbScores(
		[]string{"sad"},
	)
	fmt.Println(probs, likely)
}
