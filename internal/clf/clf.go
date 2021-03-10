package clf

import (
	"github.com/jaimeteb/chatto/internal/clf/dataset"
	"github.com/jaimeteb/chatto/internal/clf/pipeline"
)

// Classifier models a classifier and its classes
type Classifier struct {
	Model    Model
	Pipeline *pipeline.Config
}

// Model interface contains the basic functions for a model to have
type Model interface {
	// Learn takes a training dataset as input and trains the model
	// It returns the training accuracy for the model
	Learn(texts dataset.DataSet, pipe *pipeline.Config) float32

	// Predict makes a class prediction based on the trained model
	Predict(text string, pipe *pipeline.Config) (predictedClass string, proba float32)

	// Save persists the model to a file
	Save() error
}
