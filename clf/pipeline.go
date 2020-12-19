package clf

import (
	"log"
	"regexp"
	"strings"

	"github.com/spf13/viper"
)

// PipelineConfig defines a Pipeline configuration
type PipelineConfig struct {
	RemoveSymbols bool `mapstructure:"remove_symbols"`
	Lower         bool `mapstructure:"lower"`
}

// LoadPipeline returns a PipelineConfig
func LoadPipeline(path *string) PipelineConfig {
	config := viper.New()
	config.SetConfigName("pl")
	config.AddConfigPath(*path)

	if err := config.ReadInConfig(); err != nil {
		log.Println(err)
		return PipelineConfig{true, true}
	}

	var plc PipelineConfig
	if err := config.Unmarshal(&plc); err != nil {
		log.Println(err)
		return PipelineConfig{true, true}
	}

	return plc
}

// Pipeline performs steps to convert a string into a CLF input
func Pipeline(text *string, pl *PipelineConfig) []string {
	newText := *text
	if pl.RemoveSymbols {
		newText = RemoveSymbols(newText)
	}
	if pl.Lower {
		newText = Lower(newText)
	}

	tokens := Tokenize(newText)
	return tokens
}

// RemoveSymbols removes symbols from string
func RemoveSymbols(text string) string {
	re, err := regexp.Compile(`[^\w]`)
	if err != nil {
		return text
	}
	return re.ReplaceAllString(text, " ")
}

// Lower converts a string to lowercase
func Lower(text string) string {
	return strings.ToLower(text)
}

// Tokenize splits the text into tokens
func Tokenize(text string) []string {
	return strings.Split(text, " ")
}
