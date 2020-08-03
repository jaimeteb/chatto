package clf

import "strings"

// Cleaner interface has a Clean method to clean texts
// type Cleaner interface {
// 	Clean(text *string) string
// }

// Clean performs steps to clean a string
func Clean(text *string) []string {
	tokens := strings.Split(*text, " ")
	lowerTokens := make([]string, 0)
	for _, t := range tokens {
		lowerTokens = append(lowerTokens, strings.ToLower(t))
	}
	return lowerTokens
}
