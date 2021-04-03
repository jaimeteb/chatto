package extensions

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/rpc"

	"github.com/jaimeteb/chatto/internal/channels/message"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/jaimeteb/chatto/extensions"
	"github.com/jaimeteb/chatto/fsm"
	log "github.com/sirupsen/logrus"
)

// RPC is an RPC Client for extension command functions
type RPC struct {
	Client *rpc.Client
}

// Execute runs the requested extension and returns the response
func (e *RPC) Execute(extension string, messageRequest message.Request, fsmDomain *fsm.Domain, machine *fsm.FSM) error {
	req := extensions.ExecuteExtensionRequest{
		FSM:       machine,
		Extension: extension,
		Domain:    fsmDomain.NoFuncs(),
		Request:   messageRequest,
	}

	err := e.Client.Call("ListenerRPC.Execute", &req, nil)
	if err != nil {
		log.Error(err)
		return fmt.Errorf(fsmDomain.DefaultMessages.Error)
	}

	return nil
}

// REST is a REST Client for extension command functions
type REST struct {
	URL   string
	http  *retryablehttp.Client
	token string
}

// Execute runs the requested extension and returns the response
func (e *REST) Execute(extension string, messageRequest message.Request, fsmDomain *fsm.Domain, machine *fsm.FSM) error {
	req := extensions.ExecuteExtensionRequest{
		FSM:       machine,
		Domain:    fsmDomain.NoFuncs(),
		Extension: extension,
		Request:   messageRequest,
	}

	jsonReq, err := json.Marshal(req)
	if err != nil {
		log.Error(err)
		return errors.New(fsmDomain.DefaultMessages.Error)
	}

	// TODO: if fail -> don't change states

	request, err := retryablehttp.NewRequest("POST", fmt.Sprintf("%s/extension", e.URL), bytes.NewBuffer(jsonReq))
	if err != nil {
		log.Error(err)
		return errors.New(fsmDomain.DefaultMessages.Error)
	}

	request.Header.Set("Content-Type", "application/json")
	if e.token != "" {
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", e.token))
	}

	resp, err := e.http.Do(request)
	if err != nil {
		log.Error(err)
		return fmt.Errorf(fsmDomain.DefaultMessages.Error)
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	return nil
}
