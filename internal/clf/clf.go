package clf

import (
	"github.com/jaimeteb/chatto/internal/clf/dataset"
	"github.com/jaimeteb/chatto/internal/clf/knn"
	"github.com/jaimeteb/chatto/internal/clf/naivebayes"
	"github.com/jaimeteb/chatto/internal/clf/pipeline"
	log "github.com/sirupsen/logrus"
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

// New returns a trained Classifier
func New(config *Config) *Classifier {
	pipeline := config.Pipeline

	log.Info("Pipeline:")
	log.Infof("* RemoveSymbols: %v", pipeline.RemoveSymbols)
	log.Infof("* Lower:         %v", pipeline.Lower)
	log.Infof("* Threshold:     %v", pipeline.Threshold)

	log.Info("Loaded commands for classifier:")
	for i, c := range config.Classification {
		log.Infof("%2d %v", i, c.Command)
	}

	var model Model
	switch config.Model.Classifier {
	case "knn":
		model = knn.NewClassifier(
			config.Model.WordVectorsConfig,
			config.Model.ModelFile,
			config.Model.Parameters,
		)
	case "naive_bayes":
		fallthrough
	default:
		config.Model.Classifier = "naive_bayes"
		model = naivebayes.NewClassifier(
			config.Model.ModelFile,
			config.Model.Parameters,
		)
	}

	log.Infof("Using %s classifier with parameters:", config.Model.Classifier)
	for _, param := range parametersToSlice(config.Model.Parameters) {
		log.Infof(param)
	}

	log.Info("Training model...")
	acc := model.Learn(config.Classification, &pipeline)
	log.Debugf("Model training accuracy: %0.2f", acc)

	if err := model.Save(); err != nil {
		log.Error("Failed to save model:", err)
	} else {
		log.Info("Model saved successfully.")
	}

	return &Classifier{
		Model:    model,
		Pipeline: &pipeline,
	}
}
