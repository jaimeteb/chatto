package pkg

import "fmt"

// Conversation struct
type Conversation struct {
	Name string
	Path []Message
}

// History struct
type History struct {
	Messages map[string][]Message
	MaxHist  int
}

// Print show all history
func (h *History) Print(id string) {
	for ix, mess := range h.Messages[id] {
		fmt.Printf("%v:\t%v\n\t%v\n", ix, mess.Sender, mess.Text)
	}
}

// Append adds a new message to the history of id
func (h *History) Append(id string, mess Message) {
	if len(h.Messages[id]) >= h.MaxHist {
		copy(h.Messages[id], h.Messages[id][1:])
		h.Messages[id][len(h.Messages[id])-1] = mess
	} else {
		h.Messages[id] = append(h.Messages[id], mess)
	}
}
