package models

import "fmt"

// History struct
type History struct {
	Messages []Message
}

// Print show all history
func (h History) Print() {
	for ix, mess := range h.Messages {
		fmt.Printf("%v:\t%v\n\t%v\n", ix, mess.Sender, mess.Text)
	}
}
