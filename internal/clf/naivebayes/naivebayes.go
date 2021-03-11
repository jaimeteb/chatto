package naivebayes

import (
	"github.com/jaimeteb/chatto/internal/clf/dataset"
	"github.com/jaimeteb/chatto/internal/clf/pipeline"
	"github.com/navossoc/bayesian"
	log "github.com/sirupsen/logrus"
)

// Classifier is a Naïve-Bayes classifier
type Classifier struct {
	Model     *bayesian.Classifier
	Classes   []bayesian.Class
	modelFile string
	tfidf     bool
}

// NewClassifier creates a KNN classifier with file data
func NewClassifier(modelFile string, params map[string]interface{}) *Classifier {
	var tfidf bool
	ptfidf := params["tfidf"]
	switch ptfidf.(type) {
	case bool:
		tfidf = ptfidf.(bool)
	default:
		log.Errorf("Invalid value '%v' parameter 'tfidf'", ptfidf)
	}

	return &Classifier{
		modelFile: modelFile,
		tfidf:     tfidf,
	}
}

// Learn takes the training texts and trains the Naive-Bayes classifier
func (c *Classifier) Learn(texts dataset.DataSet, pipe *pipeline.Config) float32 {
	// Create model with classes
	classes := make([]bayesian.Class, 0, len(texts))
	for _, class := range texts {
		classes = append(classes, bayesian.Class(class.Command))
	}

	var classifier *bayesian.Classifier
	if c.tfidf {
		classifier = bayesian.NewClassifierTfIdf(classes...)
	} else {
		classifier = bayesian.NewClassifier(classes...)
	}

	// Train the model with clean text
	testX := make([][]string, 0)
	testY := make([]string, 0)

	for n := range texts {
		for i := range texts[n].Texts {
			cleanText := pipeline.Pipeline(texts[n].Texts[i], pipe)
			testX = append(testX, cleanText)
			classifier.Learn(cleanText, bayesian.Class(texts[n].Command))
			testY = append(testY, texts[n].Command)
		}
	}

	if c.tfidf {
		classifier.ConvertTermsFreqToTfIdf()
	}

	c.Model = classifier
	c.Classes = classes

	// Compute train accuracy
	correct := 0
	for i := range testX {
		_, likely, _ := c.Model.ProbScores(testX[i])
		pred := string(c.Classes[likely])
		if pred == testY[i] {
			correct++
		}
	}
	return float32(correct) / float32(len(testX))
}

// Predict predict a class for a given text
func (c *Classifier) Predict(text string, pipe *pipeline.Config) (predictedClass string, proba float32) {
	probs, likely, _ := c.Model.ProbScores(pipeline.Pipeline(text, pipe))
	class := string(c.Classes[likely])
	prob := probs[likely]

	log.Debugf("CLF | Text '%s' classified as command '%s' with a probability of %.2f", text, class, prob)
	if prob < pipe.Threshold {
		return "", -1.0
	}

	return class, float32(prob)
}

// Save persists the model to a file
func (c *Classifier) Save() error {
	return c.Model.WriteToFile(c.modelFile)
}
