package knn

import (
	"github.com/jaimeteb/chatto/internal/clf/dataset"
	"github.com/jaimeteb/chatto/internal/clf/embeddings"
	"github.com/jaimeteb/chatto/internal/clf/pipeline"
	log "github.com/sirupsen/logrus"
)

// Classifier is a K-Nearest Neighbors classifier
type Classifier struct {
	KNN         *KNN
	Embeddings  *embeddings.VectorMap
	truncate    float32
	vectorsFile string
	modelFile   string
}

// NewClassifier creates a KNN classifier with truncate and file data
func NewClassifier(truncate float32, vectorsFile, modelFile string) *Classifier {
	return &Classifier{
		truncate:    truncate,
		vectorsFile: vectorsFile,
		modelFile:   modelFile,
	}
}

// Learn takes the training texts and trains the K-Nearest Neighbors classifier
func (c *Classifier) Learn(texts dataset.DataSet, pipe *pipeline.Config) {
	trainX := make([][]string, 0)
	trainY := make([]string, 0)
	classes := make([]string, 0)

	// Run Pipeline
	for _, training := range texts {
		for _, trainingText := range training.Texts {
			trainX = append(trainX, pipeline.Pipeline(trainingText, pipe))
			trainY = append(trainY, training.Command)
		}
		classes = append(classes, training.Command)
	}

	// Generate VectorMap
	emb, err := embeddings.NewVectorMapFromFile(c.vectorsFile, c.truncate)
	if err != nil {
		log.Fatal(err)
	}
	c.Embeddings = emb

	// Get embeddings from dataset
	embeddingsX := make([][]float64, len(trainX))
	for i, x := range trainX {
		embeddingsX[i] = embeddings.AverageEmbeddings(c.Embeddings.Embeddings(x))
	}

	// Initialize KNN
	knn := &KNN{
		k:      3,
		data:   embeddingsX,
		labels: trainY,
	}
	c.KNN = knn

	// check accuracy
	predicted, _ := knn.PredictMany(embeddingsX)
	correct := 0
	for i := range predicted {
		if predicted[i] == trainY[i] {
			correct++
		}
	}
	log.Debugf("Train accuracy: %f\n", float64(correct)/float64(len(predicted)))
}

// Predict predict a class for a given text
func (c *Classifier) Predict(text string, pipe *pipeline.Config) (predictedClass string, proba float32) {
	x := pipeline.Pipeline(text, pipe)
	embeddingsX := embeddings.AverageEmbeddings(c.Embeddings.Embeddings(x))

	pred, prob := c.KNN.PredictOne(embeddingsX)

	if prob < pipe.Threshold {
		return "", -1.0
	}

	log.Debugf("CLF | Text '%s' classified as command '%s' with a probability of %.2f", text, pred, prob)
	return pred, float32(prob)
}
