package extensions

import (
	"encoding/json"
	"net/http"

	"github.com/jaimeteb/chatto/internal/channels/message"

	"github.com/gorilla/websocket"
	"github.com/jaimeteb/chatto/extensions"
	"github.com/jaimeteb/chatto/fsm"
	log "github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{} // use default options

// WebSocketServer is an extension server interface that sends execute request events
type WebSocketServer struct {
	ExecuteRequestQueue chan extensions.ExecuteExtensionRequest
}

// NewWebSocket returns an initialized WebSocketServer
func NewWebSocket() *WebSocketServer {
	return &WebSocketServer{ExecuteRequestQueue: make(chan extensions.ExecuteExtensionRequest)}
}

// ExtensionWebsocketHandler listens for execute extension requests and publishes them to the websocket
func (e *WebSocketServer) ExtensionWebsocketHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("upgrade", err)
		return
	}

	defer func() {
		err = c.Close()
		if err != nil {
			log.Error("close", err)
		}
	}()

	for req := range e.ExecuteRequestQueue {
		data, err := json.Marshal(req)
		if err != nil {
			log.Error("marshal", err)
			continue
		}

		err = c.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Error("write_message", err)
			continue
		}
	}
}

// Execute runs the requested extension and returns the response
func (e *WebSocketServer) Execute(extension string, msgRequest message.Request, fsmDomain *fsm.Domain, machine *fsm.FSM) error {
	req := extensions.ExecuteExtensionRequest{
		FSM:       machine,
		Domain:    fsmDomain.NoFuncs(),
		Extension: extension,
		Request:   msgRequest,
	}

	e.ExecuteRequestQueue <- req

	return nil
}
