package naivebayes

import (
	"github.com/jaimeteb/chatto/internal/clf/dataset"
	"github.com/jaimeteb/chatto/internal/clf/pipeline"
	"github.com/navossoc/bayesian"
	log "github.com/sirupsen/logrus"
)

// Classifier is a Naïve-Bayes classifier
type Classifier struct {
	Model   *bayesian.Classifier
	Classes []bayesian.Class
}

// Learn takes the training texts and trains the Naive-Bayes classifier
func (c *Classifier) Learn(texts dataset.DataSet, pipe *pipeline.Config) {
	classes := make([]bayesian.Class, 0, len(texts))

	for _, class := range texts {
		classes = append(classes, bayesian.Class(class.Command))
	}
	classifier := bayesian.NewClassifier(classes...)

	for n := range texts {
		for i := range texts[n].Texts {
			classifier.Learn(pipeline.Pipeline(texts[n].Texts[i], pipe), bayesian.Class(texts[n].Command))
		}
	}

	classes = append(classes, bayesian.Class("any"))

	c.Model = classifier
	c.Classes = classes
}

// Predict predict a class for a given text
func (c *Classifier) Predict(text string, pipe *pipeline.Config) (predictedClass string, proba float32) {
	probs, likely, _ := c.Model.ProbScores(pipeline.Pipeline(text, pipe))
	class := string(c.Classes[likely])
	prob := probs[likely]

	if prob < pipe.Threshold {
		return "", -1.0
	}

	log.Debugf("CLF | Text '%s' classified as command '%s' with a probability of %.2f", text, class, prob)
	return class, float32(prob)
}
