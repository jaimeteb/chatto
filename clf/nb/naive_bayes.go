package nb

import (
	"github.com/jaimeteb/chatto/clf/pipeline"
	cmn "github.com/jaimeteb/chatto/common"
	"github.com/navossoc/bayesian"
	log "github.com/sirupsen/logrus"
)

type NBClassifier struct {
	Model   *bayesian.Classifier
	Classes []bayesian.Class
}

// Learn takes the training texts and trains the Naive-Bayes classifier
func (c *NBClassifier) Learn(texts []cmn.TrainingTexts, pipe *pipeline.Config) {
	var classes []bayesian.Class

	for _, class := range texts {
		classes = append(classes, bayesian.Class(class.Command))
	}
	classifier := bayesian.NewClassifier(classes...)

	for _, cls := range texts {
		for _, txt := range cls.Texts {
			classifier.Learn(pipeline.Pipeline(txt, pipe), bayesian.Class(cls.Command))
		}
	}

	// classes = append(classes, bayesian.Class("any"))
	c.Model = classifier
	c.Classes = classes
}

// Predict predict a class for a given text
func (c *NBClassifier) Predict(text string, pipe *pipeline.Config) (string, float32) {
	probs, likely, _ := c.Model.ProbScores(pipeline.Pipeline(text, pipe))
	class := string(c.Classes[likely])
	prob := probs[likely]

	log.Debugf("CLF | \"%v\" classified as %v (%0.2f%%)", text, class, prob*100)
	if prob < pipe.Threshold {
		return "", -1.0
	}

	return class, float32(prob)
}
