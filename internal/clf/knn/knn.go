package knn

import (
	"encoding/gob"
	"os"
	"path"

	"github.com/jaimeteb/chatto/internal/clf/dataset"
	"github.com/jaimeteb/chatto/internal/clf/pipeline"
	"github.com/jaimeteb/chatto/internal/clf/wordvectors"
	log "github.com/sirupsen/logrus"
)

const (
	classifierFile = "clf.gob"
	modelFile      = "knn.gob"
)

// Classifier is a K-Nearest Neighbors classifier
type Classifier struct {
	KNN       *KNN
	VectorMap *wordvectors.VectorMap
	K         int
}

func (c *Classifier) SaveToFile(name string) error {
	file, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := gob.NewEncoder(file)
	return enc.Encode(&Classifier{nil, nil, c.K})
}

func NewClassifierFromFile(name string) (*Classifier, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	dec := gob.NewDecoder(file)
	classifier := new(Classifier)
	err = dec.Decode(classifier)
	return classifier, err
}

// NewClassifier creates a KNN classifier with truncate and file data
func NewClassifier(wordVecConfig wordvectors.Config, params map[string]interface{}) *Classifier {
	k := 1
	pk := params["k"]
	switch t := pk.(type) {
	case int:
		k = t
	default:
		log.Errorf("Invalid value '%v' parameter 'k'", pk)
	}

	// Generate VectorMap
	emb, err := wordvectors.NewVectorMap(&wordVecConfig)
	if err != nil {
		log.Fatal(err)
	}

	return &Classifier{
		VectorMap: emb,
		K:         k,
	}
}

// Learn takes the training texts and trains the K-Nearest Neighbors classifier
func (c *Classifier) Learn(texts dataset.DataSet, pipe *pipeline.Config) float32 {
	trainX := make([][]string, 0)
	trainY := make([]string, 0)

	// Run Pipeline
	for _, training := range texts {
		for _, trainingText := range training.Texts {
			trainX = append(trainX, pipeline.Pipeline(trainingText, pipe))
			trainY = append(trainY, training.Command)
		}
	}

	// Get embeddings from dataset
	embeddingsX := make([][]float64, len(trainX))
	for i, x := range trainX {
		embeddingsX[i] = c.VectorMap.AverageVectors(c.VectorMap.Vectors(x))
	}

	// Initialize KNN
	knn := &KNN{
		K:      c.K,
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
	embeddingsX := c.VectorMap.AverageVectors(c.VectorMap.Vectors(x))

	pred, prob := c.KNN.PredictOne(embeddingsX)

	log.Debugf("CLF | Text '%s' classified as command '%s' with a probability of %.2f", text, pred, prob)
	if prob < pipe.Threshold {
		return "", -1.0
	}

	return pred, float32(prob)
}

// Save persists the model to a file
func (c *Classifier) Save(directory string) error {
	// save Classifier
	if err := c.SaveToFile(path.Join(directory, classifierFile)); err != nil {
		return err
	}
	// save Classifier.VectorMap
	if err := c.VectorMap.SaveToFile(path.Join(directory, wordvectors.WordVectorsFile)); err != nil {
		return err
	}
	// save Classifier.KNN
	if err := c.KNN.SaveToFile(path.Join(directory, modelFile)); err != nil {
		return err
	}
	return nil
}

func Load(directory string) (classifier *Classifier, err error) {
	// load Classifier
	classifier, err = NewClassifierFromFile(path.Join(directory, classifierFile))
	if err != nil {
		return
	}

	// load Classifier.VectorMap
	vectorMap, err := wordvectors.NewVectorMapFromFile(path.Join(directory, wordvectors.WordVectorsFile))
	if err != nil {
		return
	}
	classifier.VectorMap = vectorMap

	// load Classifier.KNN
	knn, err := NewKNNClassifierFromFile(path.Join(directory, modelFile))
	if err != nil {
		return
	}
	classifier.KNN = knn

	return
}
