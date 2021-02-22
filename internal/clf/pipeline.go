package clf

import (
	"regexp"
	"strings"
)

var removeSymbolRe = regexp.MustCompile(`\W`)

// PipelineConfig defines a Pipeline configuration
type PipelineConfig struct {
	RemoveSymbols bool    `mapstructure:"remove_symbols"`
	Lower         bool    `mapstructure:"lower"`
	Threshold     float64 `mapstructure:"threshold"`
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
	return removeSymbolRe.ReplaceAllString(text, " ")
}

// Lower converts a string to lowercase
func Lower(text string) string {
	return strings.ToLower(text)
}

// Tokenize splits the text into tokens
func Tokenize(text string) []string {
	return strings.Split(text, " ")
}
