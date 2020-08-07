package ext

import (
	"fmt"

	"github.com/jaimeteb/chatto/fsm"
)

// PrintSlots prints all slots
func PrintSlots(m *fsm.FSM) {
	for k, v := range m.Slots {
		fmt.Printf("%v: %v\n", k, v)
	}
}
