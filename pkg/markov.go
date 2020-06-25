package pkg

import (
	"fmt"
	"math/rand"
	"strings"
)

// Chain maps states as strings to messages as strings.
// The strings are formed with the name of the message sender ("usr" or "bot").
type Chain struct {
	chain     map[string][]string
	stateSize int
}

// NewChain returns a new Chain with states of size stateSize.
func NewChain(stateSize int) *Chain {
	return &Chain{make(map[string][]string), stateSize}
}

// State is a Markov chain state of one or more messages.
type State []string

func (p State) shift(token string) {
	copy(p, p[1:])
	p[len(p)-1] = token
}

func (p State) string() string {
	return strings.Join(p, " ")
}

// Build builds a Markov Chain based on the conversations
func (c *Chain) Build(convs []Conversation) {
	// p := make(State, c.stateSize)
	for _, conv := range convs {
		fmt.Printf("## %v\n", conv.Name)

		p := make(State, c.stateSize)
		for _, mess := range conv.Path {
			messCode := mess.Text // messCode := fmt.Sprintf("%v:%v", mess.Sender, mess.Text)
			fmt.Printf("%v -> %v\n", p, messCode)

			key := p.string()
			c.chain[key] = append(c.chain[key], messCode)
			p.shift(messCode)
		}
	}
}

// Generate returns up to n messages from a chain.
func (c *Chain) Generate(n int) string {
	p := make(State, c.stateSize)
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

// Predict takes the current state of the conversation to predict the next step.
func (c *Chain) Predict(curr *[]Message) {
	p := make(State, c.stateSize)

	l := len(*curr)
	states := (*curr)[l-c.stateSize : l]

	for _, mess := range states {
		messCode := mess.Text

		p.shift(messCode)
	}

	fmt.Println("prefix ", p)
}
