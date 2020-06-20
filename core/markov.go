package core

import (
	"fmt"

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

// Build builds a Markov Chain
// based on the conversations
// func (c *Chain) Build(convs []models.Conversation) {
func (c *Chain) Build(convs []models.Conversation) {
	// p := make(Prefix, c.prefixLen)
	for _, conv := range convs {
		fmt.Printf("## %v\n", conv.Name)

		p := make(Prefix, c.prefixLen)
		for _, mess := range conv.Path {
			messKey := fmt.Sprintf("%v:%v", mess.Sender, mess.Text)
			fmt.Printf("%v -> %v\n", p, messKey)

			p.shift(messKey)

			// c.chain[key] = append(c.chain[key])
		}
	}
}
