package knn

import (
	"github.com/jaimeteb/chatto/internal/clf/dataset"
	"github.com/jaimeteb/chatto/internal/clf/pipeline"
)

// Classifier is a K-Nearest Neighbors classifier
type Classifier struct {
}

// Learn takes the training texts and trains the K-Nearest Neighbors classifier
func (c *Classifier) Learn(texts dataset.DataSet, pipe *pipeline.Config) {
}

// Predict predict a class for a given text
func (c *Classifier) Predict(text string, pipe *pipeline.Config) (predictedClass string, proba float32) {
	return
}
