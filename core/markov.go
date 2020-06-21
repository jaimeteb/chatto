package core

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/jaimeteb/chatto/models"
)

// MarkovDecision takes a decision
// based on a Markov Chain.
func MarkovDecision(bot *models.Bot) {
	return
}

// Chain contains a map ("chain") of prefixes to a list of suffixes.
// A prefix is a slice of prefixLen messages.
// A suffix is a message.
// A prefix can have multiple suffixes.
type Chain struct {
	chain     map[string][]string
	prefixLen int
}

// NewChain returns a new Chain with prefixes of prefixLen words.
func NewChain(prefixLen int) *Chain {
	return &Chain{make(map[string][]string), prefixLen}
}

// Prefix is a Markov chain prefix of one or more messages.
type Prefix []string

func (p Prefix) shift(token string) {
	copy(p, p[1:])
	p[len(p)-1] = token
}

func (p Prefix) string() string {
	return strings.Join(p, " ")
}

// Build builds a Markov Chain
// based on the conversations
func (c *Chain) Build(convs []models.Conversation) {
	// p := make(Prefix, c.prefixLen)
	for _, conv := range convs {
		fmt.Printf("## %v\n", conv.Name)

		p := make(Prefix, c.prefixLen)
		for _, mess := range conv.Path {
			messCode := fmt.Sprintf("%v:%v", mess.Sender, mess.Text)
			fmt.Printf("%v -> %v\n", p, messCode)

			key := p.string()
			c.chain[key] = append(c.chain[key], messCode)
			p.shift(messCode)
		}
	}
}

// Generate returns a string of at most n words generated from Chain.
func (c *Chain) Generate(n int) string {
	p := make(Prefix, c.prefixLen)
	var words []string
	for i := 0; i < n; i++ {
		choices := c.chain[p.string()]
		if len(choices) == 0 {
			break
		}
		next := choices[rand.Intn(len(choices))]
		words = append(words, next)
		p.shift(next)
	}
	return strings.Join(words, " ")
}
