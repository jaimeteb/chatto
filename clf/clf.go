package clf

import (
	"github.com/jaimeteb/chatto/clf/nb"
	"github.com/jaimeteb/chatto/clf/pipeline"
	cmn "github.com/jaimeteb/chatto/common"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

// Model interface contains the basic functions for a model to have
type Model interface {
	// Learn takes training texts as input and trains the model
	Learn(texts []cmn.TrainingTexts, pipe *pipeline.Config)
	// Predict makes a class prediction based on the trained model
	Predict(text string, pipe *pipeline.Config) (string, float32)
}

// Config models a classification yaml file
type Config struct {
	Classification []cmn.TrainingTexts `yaml:"classification"`
	Pipeline       pipeline.Config     `yaml:"pipeline"`
	Model          string              `yaml:"model"`
}

// Classifier models a classifier and its classes
type Classifier struct {
	Model    Model
	Pipeline *pipeline.Config
}

// Load loads classification configuration from yaml
func Load(path *string) *Config {
	config := viper.New()
	config.SetConfigName("clf")
	config.AddConfigPath(*path)

	if err := config.ReadInConfig(); err != nil {
		log.Panic(err)
	}

	botClassif := new(Config)
	config.Unmarshal(botClassif)

	return botClassif
}

// Create returns a trained Classifier
func Create(path *string) *Classifier {
	config := Load(path)

	pipeline := &config.Pipeline

	var model Model
	switch config.Model {
	// case "random_forest":
	// case "knn":
	case "naive_bayes":
		fallthrough
	default:
		model = new(nb.NBClassifier)
		model.Learn(config.Classification, pipeline)
	}

	log.Info("Pipeline:")
	log.Infof("* RemoveSymbols: %v\n", pipeline.RemoveSymbols)
	log.Infof("* Lower:         %v\n", pipeline.Lower)
	log.Infof("* Threshold:     %v\n", pipeline.Threshold)

	log.Info("Loaded commands for classifier:")
	log.Infof("%2d\t%s\n", -1, "any")
	for i, c := range config.Classification {
		log.Infof("%2d\t%s\n", i, c.Command)
	}

	return &Classifier{
		Model:    model,
		Pipeline: pipeline,
	}
}
