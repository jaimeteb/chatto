package clf

import (
	"github.com/fsnotify/fsnotify"
	"github.com/navossoc/bayesian"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config models a classification yaml file
type Config struct {
	Classification []TrainingTexts `yaml:"classification"`
	Pipeline       PipelineConfig  `yaml:"pipeline"`
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
	classes := make([]bayesian.Class, 0, len(classifConfig.Classification))
	for _, class := range classifConfig.Classification {
		classes = append(classes, bayesian.Class(class.Command))
	}

	classifier := bayesian.NewClassifier(classes...)
	pipeline := classifConfig.Pipeline

	log.Info("Pipeline:")
	log.Infof("* RemoveSymbols: %v", pipeline.RemoveSymbols)
	log.Infof("* Lower:         %v", pipeline.Lower)
	log.Infof("* Threshold:     %v", pipeline.Threshold)

	classif := classifConfig.Classification

	for n := range classif {
		for i := range classif[n].Texts {
			classifier.Learn(Pipeline(&classif[n].Texts[i], &pipeline), bayesian.Class(classif[n].Command))
		}
	}

	classes = append(classes, bayesian.Class("any"))

	log.Info("Loaded commands for classifier:")
	for i, c := range classes {
		log.Infof("%2d %v", i, c)
	}

	return &Classifier{classifier, classes, pipeline}
}
