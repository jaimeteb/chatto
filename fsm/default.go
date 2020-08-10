package fsm

import (
	"fmt"
)

// PrintSlots prints all slots
func PrintSlots(m *FSM) {
	for k, v := range m.Slots {
		fmt.Printf("%v: %v\n", k, v)
	}
}
