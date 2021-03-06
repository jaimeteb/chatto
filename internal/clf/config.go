package clf

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/jaimeteb/chatto/internal/clf/dataset"
	"github.com/jaimeteb/chatto/internal/clf/pipeline"
	"github.com/jaimeteb/chatto/internal/clf/wordvectors"
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

	// Directory is the path to the saved model files
	Directory string `mapstructure:"directory"`

	// Whether to save the trained model or not
	Save bool `mapstructure:"save"`

	// Whether to load the saved model or not
	Load bool `mapstructure:"load"`

	// WordVectorsConfig contains configuration for fasttext word vectors
	WordVectorsConfig wordvectors.Config `mapstructure:"word_vectors"`
}

// LoadConfig loads classification configuration from yaml
func LoadConfig(path string, reloadChan chan Config) (*Config, error) {
	config := viper.New()
	config.SetConfigName("clf")
	config.AddConfigPath(path)
	config.SetDefault("pipeline.remove_symbols", true)
	config.SetDefault("pipeline.lower", true)
	config.SetDefault("pipeline.threshold", 0.1)
	config.SetDefault("model.directory", "./model")
	config.SetDefault("model.word_vectors.truncate", 1.0)

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

func parametersToSlice(params map[string]interface{}) []string {
	s := make([]string, len(params))
	i := 0
	for param, value := range params {
		s[i] = fmt.Sprintf("* %s: %v", param, value)
		i++
	}
	return s
}
