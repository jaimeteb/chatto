package pkg

import (
	"fmt"
	"log"
	"strings"

	"github.com/navossoc/bayesian"
	"github.com/spf13/viper"
)

// Classification models a classification yaml file
type Classification struct {
	Classification []TrainingTexts `yaml:"classification"`
}

// TrainingTexts models texts used for training the classifier
type TrainingTexts struct {
	Command string   `yaml:"command"`
	Texts   []string `yaml:"texts"`
}

// Classifier models a classifier and its classes
type Classifier struct {
	Model   bayesian.Classifier
	Classes []bayesian.Class
}

// Predict predict a class for a given text
func (c *Classifier) Predict(text string) (string, float64) {
	probs, likely, _ := c.Model.ProbScores(strings.Split(text, " "))
	log.Println(probs, likely)
	return string(c.Classes[likely]), probs[likely]
}

// LoadClassificationConfig loads classification configuration from yaml
func LoadClassificationConfig() Classification {
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

	fmt.Println(botClassif)
	return botClassif
}

// GetClassifier returns a trained Classifier
func GetClassifier() Classifier {
	classification := LoadClassificationConfig()

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

	// for _, test := range []string{"hi", "I am good", "bad", "oh yes", "oh no"} {
	// 	probs, likely, _ := classifier.ProbScores(
	// 		strings.Split(test, " "),
	// 	)
	// 	fmt.Println(probs, likely)
	// }

	return Classifier{*classifier, classes}
}
