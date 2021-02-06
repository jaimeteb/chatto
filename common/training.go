package common

// TrainingTexts models texts used for training the classifier
type TrainingTexts struct {
	Command string   `yaml:"command"`
	Texts   []string `yaml:"texts"`
}
