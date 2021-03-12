package clf

import (
	"os"

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
	Save(directory string) error

	// Load loads the model from a file
	// Load(directory string) (Model, error)
}

// New returns a trained Classifier
func New(config *Config) *Classifier {
	pipe := config.Pipeline

	log.Info("Pipeline:")
	log.Infof("* RemoveSymbols: %v", pipe.RemoveSymbols)
	log.Infof("* Lower:         %v", pipe.Lower)
	log.Infof("* Threshold:     %v", pipe.Threshold)

	log.Info("Loaded commands for classifier:")
	for i, c := range config.Classification {
		log.Infof("%2d %v", i, c.Command)
	}

	log.Infof("Using %s classifier with parameters:", config.Model.Classifier)
	for _, param := range parametersToSlice(config.Model.Parameters) {
		log.Infof(param)
	}

	var model Model
	switch config.Model.Classifier {
	case "knn":
		model = knn.NewClassifier(
			config.Model.WordVectorsConfig,
			config.Model.Parameters,
		)
	default:
		config.Model.Classifier = "naive_bayes"
		model = naivebayes.NewClassifier(
			config.Model.Parameters,
		)
	}

	if config.Model.Load {
		// Check for directory to load
		if _, err := os.Stat(config.Model.Directory); os.IsNotExist(err) {
			log.Error("Failed to load model:", err)
			// TODO: exit?
		}

		// Load model
		log.Info("Loading model...")
		var err error
		switch config.Model.Classifier {
		case "knn":
			model, err = knn.Load(config.Model.Directory)
		case "naive_bayes":
			model, err = naivebayes.Load(config.Model.Directory)
		}
		if err != nil {
			log.Error("Failed to load model:", err)
			// TODO: exit?
		} else {
			log.Info("Model loaded successfully.")
		}
	} else {
		// Train model
		log.Info("Training model...")
		acc := model.Learn(config.Classification, &pipe)
		log.Debugf("Model training accuracy: %0.2f", acc)

		// Check for directory to save
		if _, err := os.Stat(config.Model.Directory); os.IsNotExist(err) {
			if err := os.MkdirAll(config.Model.Directory, 0755); err != nil {
				log.Error("Couldn't create directory:", err)
			}
		}

		// Save model
		if err := model.Save(config.Model.Directory); err != nil {
			log.Error("Failed to save model:", err)
		} else {
			log.Info("Model saved successfully.")
		}
	}

	return &Classifier{
		Model:    model,
		Pipeline: &pipe,
	}
}
