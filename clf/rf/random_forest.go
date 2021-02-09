package rf

import (
	"github.com/jaimeteb/chatto/clf/pipeline"
	cmn "github.com/jaimeteb/chatto/common"
	log "github.com/sirupsen/logrus"
	"github.com/sjwhitworth/golearn/ensemble"
)

// RandomForestClassifier contains a golearn RandomForest and necessary metadata
type RandomForestClassifier struct {
	Model       *ensemble.RandomForest
	VectorMap   *VectorMap
	truncate    float32
	vectorsFile string
	modelFile   string
	classes     []string
}

// NewRandomForestClassifier returns a RandomForestClassifier object with the necessary metadata
// truncate will default to 1.0 (will use the full vector file)
// modelFile will default to "model/rf.clf"
func NewRandomForestClassifier(truncate float32, vectorsFile, modelFile string) *RandomForestClassifier {
	if truncate == 0.0 {
		truncate = 1.0
	}
	if modelFile == "" {
		modelFile = "model/rf.clf"
	}

	return &RandomForestClassifier{
		truncate:    truncate,
		vectorsFile: vectorsFile,
		modelFile:   modelFile,
	}
}

// Learn will create a new VectorMap and generate the training data from the training texts
// A RandomForestClassifier will be trained and saved
func (c *RandomForestClassifier) Learn(texts []cmn.TrainingTexts, pipe *pipeline.Config) {
	trainX := make([][]string, 0)
	trainY := make([]string, 0)
	classes := make([]string, 0)

	for _, training := range texts {
		for _, trainingText := range training.Texts {
			trainX = append(trainX, pipeline.Pipeline(trainingText, pipe))
			trainY = append(trainY, training.Command)
		}
		classes = append(classes, training.Command)
	}

	c.VectorMap = NewVectorMapFromFile(c.vectorsFile, c.truncate)

	trainingData := c.VectorMap.ConvertDataToInstances(trainX, trainY)

	cls := ensemble.NewRandomForest(20, 10)
	cls.Fit(trainingData)
	cls.Save(c.modelFile)

	c.classes = classes
	c.Model = cls
}

// Predict creates a new golearn instance and predicts the class of the passed text
// If successful, confidence will always be 1.0
//
// TODO: replace golearn with other library, in order to make predictions more
// accurate and intuitive
//
func (c *RandomForestClassifier) Predict(text string, pipe *pipeline.Config) (string, float32) {
	processedText := pipeline.Pipeline(text, pipe)
	inputTexts := make([][]string, len(c.classes))
	for i := 0; i < len(c.classes); i++ {
		inputTexts[i] = processedText
	}

	inst := c.VectorMap.ConvertDataToInstances(inputTexts, c.classes)

	pred, err := c.Model.Predict(inst)
	if err != nil {
		log.Error(err)
		return "", -1.0
	}
	classPred := pred.RowString(0)
	log.Debugf("CLF | \"%v\" classified as %v", text, classPred)

	return classPred, 1.0
}
