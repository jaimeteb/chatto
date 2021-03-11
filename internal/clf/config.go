package clf

import (
	"github.com/fsnotify/fsnotify"
	"github.com/jaimeteb/chatto/internal/clf/dataset"
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

	// Parameters holds any other model parameters
	// TODO: improve this field
	Parameters map[string]interface{} `mapstructure:"parameters"`

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
