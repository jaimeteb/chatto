package knn

import (
	"encoding/gob"
	"os"

	"github.com/jaimeteb/chatto/internal/clf/dataset"
	"github.com/jaimeteb/chatto/internal/clf/embeddings"
	"github.com/jaimeteb/chatto/internal/clf/pipeline"
	log "github.com/sirupsen/logrus"
)

// Classifier is a K-Nearest Neighbors classifier
type Classifier struct {
	KNN        *KNN
	Embeddings *embeddings.VectorMap
	modelFile  string
	k          int
}

// NewClassifier creates a KNN classifier with truncate and file data
func NewClassifier(wordVecConfig embeddings.WordVectorsConfig, modelFile string, params map[string]interface{}) *Classifier {
	k := 1
	pk := params["k"]
	switch pk.(type) {
	case int:
		k = pk.(int)
	default:
		log.Errorf("Invalid value '%v' parameter 'k'", pk)
	}

	// Generate VectorMap
	emb, err := embeddings.NewVectorMap(&wordVecConfig)
	if err != nil {
		log.Fatal(err)
	}

	return &Classifier{
		Embeddings: emb,
		modelFile:  modelFile,
		k:          k,
	}
}

// Learn takes the training texts and trains the K-Nearest Neighbors classifier
func (c *Classifier) Learn(texts dataset.DataSet, pipe *pipeline.Config) float32 {
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

	// Get embeddings from dataset
	embeddingsX := make([][]float64, len(trainX))
	for i, x := range trainX {
		embeddingsX[i] = c.Embeddings.AverageEmbeddings(c.Embeddings.Embeddings(x))
	}

	// Initialize KNN
	knn := &KNN{
		K:      c.k,
		Data:   embeddingsX,
		Labels: trainY,
	}
	c.KNN = knn

	// Compute train accuracy
	preds, _ := c.KNN.PredictMany(embeddingsX)
	correct := 0
	for i, pred := range preds {
		if pred == trainY[i] {
			correct++
		}
	}
	return float32(correct) / float32(len(preds))
}

// Predict predict a class for a given text
func (c *Classifier) Predict(text string, pipe *pipeline.Config) (predictedClass string, proba float32) {
	x := pipeline.Pipeline(text, pipe)
	embeddingsX := c.Embeddings.AverageEmbeddings(c.Embeddings.Embeddings(x))

	pred, prob := c.KNN.PredictOne(embeddingsX)

	log.Debugf("CLF | Text '%s' classified as command '%s' with a probability of %.2f", text, pred, prob)
	if prob < pipe.Threshold {
		return "", -1.0
	}

	return pred, float32(prob)
}

// Save persists the model to a file
func (c *Classifier) Save() error {
	file, err := os.OpenFile(c.modelFile, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := gob.NewEncoder(file)
	return enc.Encode(c.KNN)
}
