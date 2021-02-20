package clf

import (
	log "github.com/sirupsen/logrus"

	"github.com/navossoc/bayesian"
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
	probs, likely, _ := c.Model.ProbScores(Pipeline(&text, &c.Pipeline))
	class := string(c.Classes[likely])
	prob := probs[likely]

	log.Debugf("CLF | \"%v\" classified as %v (%0.2f%%)", text, class, prob*100)
	if prob < c.Pipeline.Threshold {
		return "", -1.0
	}

	return class, prob
}
