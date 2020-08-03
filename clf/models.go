package clf

import (
	"github.com/navossoc/bayesian"
)

// Classification models a classification yaml file
type Classification struct {
	Classification []TrainingTexts `yaml:"classification"`
}

// TrainingTexts models texts used for training the classifier
type TrainingTexts struct {
	Command string   `yaml:"command"`
	Texts   []string `yaml:"texts"`
}

// Classifier models a classifier and its classes
type Classifier struct {
	Model   bayesian.Classifier
	Classes []bayesian.Class
}

// Predict predict a class for a given text
func (c *Classifier) Predict(text string) (string, float64) {
	probs, likely, _ := c.Model.ProbScores(Clean(&text))
	class := string(c.Classes[likely])
	prob := probs[likely]
	// log.Printf("CLF\t| \"%v\" classified as %v (%0.2f%%)", text, class, prob*100)
	return class, prob
}
