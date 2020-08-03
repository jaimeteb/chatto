package clf

import (
	"log"
	"strings"

	"github.com/navossoc/bayesian"
	"github.com/spf13/viper"
)

// Load loads classification configuration from yaml
func Load(path *string) Classification {
	config := viper.New()
	config.SetConfigName("clf")
	config.AddConfigPath(*path)

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
func Create(path *string) Classifier {
	classification := Load(path)

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

	log.Println("Loaded commands for classifier:")
	for i, c := range classes {
		log.Printf("%v\t%v\n", i, c)
	}

	return Classifier{*classifier, classes}
}
