package models

import "fmt"

// History struct
type History struct {
	Messages map[string][]Message
}

// Print show all history
func (h History) Print(id string) {
	for ix, mess := range h.Messages[id] {
		fmt.Printf("%v:\t%v\n\t%v\n", ix, mess.Sender, mess.Text)
	}
}
