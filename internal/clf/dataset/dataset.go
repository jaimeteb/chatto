package dataset

// DataSet contains multiple dataclasses
type DataSet []DataClass

// DataClass models texts used for training the classifier
type DataClass struct {
	Command string   `yaml:"command"`
	Texts   []string `yaml:"texts"`
}
