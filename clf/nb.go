package clf

import (
	log "github.com/sirupsen/logrus"

	"github.com/navossoc/bayesian"
	"github.com/spf13/viper"
)

// Classification models a classification yaml file
type Classification struct {
	Classification []TrainingTexts `yaml:"classification"`
	Pipeline       PipelineConfig  `yaml:"pipeline"`
}

// TrainingTexts models texts used for training the classifier
type TrainingTexts struct {
	Command string   `yaml:"command"`
	Texts   []string `yaml:"texts"`
}

// Classifier models a classifier and its classes
type Classifier struct {
	Model    bayesian.Classifier
	Classes  []bayesian.Class
	Pipeline PipelineConfig
}

// Predict predict a class for a given text
func (c *Classifier) Predict(text string) (string, float64) {
	probs, likely, _ := c.Model.ProbScores(Pipeline(&text, &c.Pipeline))
	class := string(c.Classes[likely])
	prob := probs[likely]

	log.Debugf("CLF | \"%v\" classified as %v (%0.2f%%)", text, class, prob*100)
	if prob < c.Pipeline.Threshold {
		return "", -1.0
	}

	return class, prob
}

// Load loads classification configuration from yaml
func Load(path *string) Classification {
	config := viper.New()
	config.SetConfigName("clf")
	config.AddConfigPath(*path)

	if err := config.ReadInConfig(); err != nil {
		log.Panic(err)
	}

	var botClassif Classification
	config.Unmarshal(&botClassif)

	return botClassif
}

// Create returns a trained Classifier
func Create(path *string) Classifier {
	classification := Load(path)

	// classes := make([]bayesian.Class, 0)
	var classes []bayesian.Class
	for _, class := range classification.Classification {
		classes = append(classes, bayesian.Class(class.Command))
	}

	classifier := bayesian.NewClassifier(classes...)
	pipeline := classification.Pipeline

	log.Info("Pipeline:")
	log.Infof("* RemoveSymbols: \t%v\n", pipeline.RemoveSymbols)
	log.Infof("* Lower: \t\t%v\n", pipeline.Lower)
	log.Infof("* Threshold: \t%v\n", pipeline.Threshold)

	for _, cls := range classification.Classification {
		for _, txt := range cls.Texts {
			classifier.Learn(Pipeline(&txt, &pipeline), bayesian.Class(cls.Command))
		}
	}

	classes = append(classes, bayesian.Class("any"))

	log.Info("Loaded commands for classifier:")
	for i, c := range classes {
		log.Infof("%v\t%v\n", i, c)
	}

	return Classifier{*classifier, classes, pipeline}
}
