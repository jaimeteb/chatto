package pipeline

import (
	"regexp"
	"strings"
)

var removeSymbolRe = regexp.MustCompile(`\W+`)

// Config defines a Pipeline configuration
type Config struct {
	RemoveSymbols bool    `mapstructure:"remove_symbols"`
	Lower         bool    `mapstructure:"lower"`
	Threshold     float64 `mapstructure:"threshold"`
}

// Pipeline performs steps to convert a string into a CLF input
func Pipeline(text string, pl *Config) []string {
	if pl.RemoveSymbols {
		text = RemoveSymbols(text)
	}
	if pl.Lower {
		text = Lower(text)
	}

	tokens := Tokenize(text)
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
