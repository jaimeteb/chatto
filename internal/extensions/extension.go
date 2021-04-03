package extensions

import (
	"fmt"

	"github.com/jaimeteb/chatto/internal/channels/message"

	"github.com/jaimeteb/chatto/fsm"
)

// ServerMap of extension server names to their clients
type ServerMap map[string]Extension

// Add new extension name and client to the extension map
func (m *ServerMap) Add(server string, client Extension) error {
	extensionMap := *m

	if _, ok := extensionMap[server]; ok {
		return fmt.Errorf("duplicate extension server found: %s", server)
	}

	extensionMap[server] = client

	return nil
}

// Extension is a service (REST or RPC or WebSocketServer) that executes commands and returns
// an answer to the Chatto bot. Extensions are written in any language and
// do whatever you want.
type Extension interface {
	Execute(extension string, msgRequest message.Request, fsmDomain *fsm.Domain, machine *fsm.FSM) error
}
