package clf

import (
	"github.com/navossoc/bayesian"
	log "github.com/sirupsen/logrus"
)

// TrainingTexts models texts used for training the classifier
type TrainingTexts struct {
	Command string   `yaml:"command"`
	Texts   []string `yaml:"texts"`
}

// Classifier models a classifier and its classes
type Classifier struct {
	Model    *bayesian.Classifier
	Classes  []bayesian.Class
	Pipeline PipelineConfig
}

// Predict predict a class for a given text
func (c *Classifier) Predict(text string) (string, float64) {
	probabilities, likely, _ := c.Model.ProbScores(Pipeline(&text, &c.Pipeline))
	class := string(c.Classes[likely])
	probability := probabilities[likely]

	log.Debugf("CLF | Text '%s' classified as command '%s' with a probability of %0.2f%%", text, class, probability*100)

	if probability < c.Pipeline.Threshold {
		return "", -1.0
	}

	return class, probability
}
