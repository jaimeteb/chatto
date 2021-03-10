package clf

import (
	"github.com/fsnotify/fsnotify"
	"github.com/jaimeteb/chatto/internal/clf/dataset"
	"github.com/jaimeteb/chatto/internal/clf/knn"
	"github.com/jaimeteb/chatto/internal/clf/naivebayes"
	"github.com/jaimeteb/chatto/internal/clf/pipeline"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config models a classification yaml file
type Config struct {
	Classification dataset.DataSet `yaml:"classification"`
	Pipeline       pipeline.Config `yaml:"pipeline"`
	Model          ModelConfig     `yaml:"model"`
}

// ModelConfig defines the data for the model's operation
type ModelConfig struct {
	// Classifier is the type of classifier to be used
	Classifier string `mapstructure:"classifier"`

	// Truncate is a number between 0 and 1, which represents how many
	// words will be used from the word embeddings
	Truncate float32 `mapstructure:"truncate"`

	// VectorsFile is the path to the word embeddings or vectors file
	VectorsFile string `mapstructure:"vectors_file"`

	// ModelFile is the path to the saved model
	ModelFile string `mapstructure:"model_file"`
}

// LoadConfig loads classification configuration from yaml
func LoadConfig(path string, reloadChan chan Config) (*Config, error) {
	config := viper.New()
	config.SetConfigName("clf")
	config.AddConfigPath(path)

	config.WatchConfig()
	config.OnConfigChange(func(in fsnotify.Event) {
		if in.Op == fsnotify.Create || in.Op == fsnotify.Write {
			log.Info("Reloading CLF configuration.")

			if err := config.ReadInConfig(); err != nil {
				log.Error(err)
				return
			}

			var classifConfig Config
			if err := config.Unmarshal(&classifConfig); err != nil {
				log.Error(err)
				return
			}

			reloadChan <- classifConfig
		}
	})

	if err := config.ReadInConfig(); err != nil {
		return nil, err
	}

	var classifConfig Config
	err := config.Unmarshal(&classifConfig)

	return &classifConfig, err
}

// New returns a trained Classifier
func New(classifConfig *Config) *Classifier {
	pipeline := classifConfig.Pipeline

	log.Info("Pipeline:")
	log.Infof("* RemoveSymbols: %v", pipeline.RemoveSymbols)
	log.Infof("* Lower:         %v", pipeline.Lower)
	log.Infof("* Threshold:     %v", pipeline.Threshold)

	log.Info("Loaded commands for classifier:")
	for i, c := range classifConfig.Classification {
		log.Infof("%2d %v", i, c.Command)
	}

	var model Model
	switch classifConfig.Model.Classifier {
	case "knn":
		model = knn.NewClassifier(classifConfig.Model.Truncate, classifConfig.Model.VectorsFile, classifConfig.Model.ModelFile)
	case "naive_bayes":
		fallthrough
	default:
		model = naivebayes.NewClassifier(classifConfig.Model.ModelFile)
	}

	log.Info("Training model...")
	model.Learn(classifConfig.Classification, &pipeline)
	log.Debugf("Model training accuracy: %0.2f", model.Accuracy(classifConfig.Classification, &pipeline))

	if err := model.Save(); err != nil {
		log.Error("Failed to save model:", err)
	} else {
		log.Info("Model saved successfully.")
	}

	return &Classifier{
		Model:    model,
		Pipeline: &pipeline,
	}
}
