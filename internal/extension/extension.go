package extension

import (
	"fmt"

	"github.com/jaimeteb/chatto/fsm"
	"github.com/jaimeteb/chatto/query"
)

// Map of extension names to extension clients
type Map map[string]Extension

// Add new extension name and client to the extension map
func (m *Map) Add(server string, client Extension) error {
	extensionMap := *m

	if _, ok := extensionMap[server]; ok {
		return fmt.Errorf("duplicate extension server found: %s", server)
	}

	extensionMap[server] = client

	return nil
}

// Extension is a service (REST or RPC) that executes commands and returns
// an answer to the Chatto bot. Extensions are written in any language and
// do whatever you want.
type Extension interface {
	GetAllExtensions() ([]string, error)
	ExecuteExtension(question *query.Question, extensionName, channel string, fsmDomain *fsm.Domain, machine *fsm.FSM) ([]query.Answer, error)
}
