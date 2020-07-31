package clf

import (
	"strings"

	"github.com/navossoc/bayesian"
	"github.com/spf13/viper"
)

// Load loads classification configuration from yaml
func Load() Classification {
	config := viper.New()
	config.SetConfigName("classification")
	config.AddConfigPath(".")
	config.AddConfigPath("config")

	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}

	var botClassif Classification
	if err := config.Unmarshal(&botClassif); err != nil {
		panic(err)
	}

	return botClassif
}

// Create returns a trained Classifier
func Create() Classifier {
	classification := Load()

	// classes := make([]bayesian.Class, 0)
	var classes []bayesian.Class
	for _, class := range classification.Classification {
		classes = append(classes, bayesian.Class(class.Command))
	}

	classifier := bayesian.NewClassifier(classes...)

	for _, cls := range classification.Classification {
		for _, txt := range cls.Texts {
			tokens := strings.Split(txt, " ")
			classifier.Learn(tokens, bayesian.Class(cls.Command))
		}
	}

	return Classifier{*classifier, classes}
}
