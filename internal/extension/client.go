package extension

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/rpc"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/jaimeteb/chatto/extensions"
	"github.com/jaimeteb/chatto/fsm"
	"github.com/jaimeteb/chatto/query"
	log "github.com/sirupsen/logrus"
)

// RPC is an RPC Client for extension command functions
type RPC struct {
	Client *rpc.Client
}

// ExecuteExtension runs the requested command function and returns the response
func (e *RPC) ExecuteExtension(question *query.Question, ext, chn, cmd string, fsmDomain *fsm.Domain, machine *fsm.FSM) ([]query.Answer, error) {
	req := extensions.ExecuteExtensionRequest{
		FSM:       machine,
		Extension: ext,
		Question:  question,
		Domain:    fsmDomain.NoFuncs(),
		Channel:   chn,
		Command:   cmd,
	}

	res := extensions.ExecuteExtensionResponse{}

	err := e.Client.Call("ListenerRPC.ExecuteExtension", &req, &res)
	if err != nil {
		return nil, errors.New(fsmDomain.DefaultMessages.Error)
	}

	*machine = *res.FSM

	return res.Answers, nil
}

// GetAllExtensions returns all command functions in the extension as a list of strings
func (e *RPC) GetAllExtensions() ([]string, error) {
	req := new(extensions.ExecuteExtensionRequest)
	res := new(extensions.GetAllExtensionsResponse)
	if err := e.Client.Call("ListenerRPC.GetAllExtensions", &req, &res); err != nil {
		log.Error(err)
		return nil, err
	}

	return res.Extensions, nil
}

// REST is a REST Client for extension command functions
type REST struct {
	URL     string
	http    *retryablehttp.Client
	token   string
	onError string
}

// ExecuteExtension runs the requested command function and returns the response
func (e *REST) ExecuteExtension(question *query.Question, ext, chn, cmd string, fsmDomain *fsm.Domain, machine *fsm.FSM) ([]query.Answer, error) {
	req := extensions.ExecuteExtensionRequest{
		FSM:       machine,
		Extension: ext,
		Question:  question,
		Domain:    fsmDomain.NoFuncs(),
		Channel:   chn,
		Command:   cmd,
	}

	jsonReq, err := json.Marshal(req)
	if err != nil {
		return nil, errors.New(fsmDomain.DefaultMessages.Error)
	}

	switch e.onError {
	case "continue":
	case "stay":
	case "restart":
	}

	request, err := retryablehttp.NewRequest("POST", fmt.Sprintf("%s/extension", e.URL), bytes.NewBuffer(jsonReq))
	if err != nil {
		return nil, errors.New(fsmDomain.DefaultMessages.Error)
	}
	request.Header.Set("Content-Type", "application/json")
	if e.token != "" {
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", e.token))
	}
	resp, err := e.http.Do(request)

	if err != nil {
		return nil, errors.New(fsmDomain.DefaultMessages.Error)
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	res := extensions.ExecuteExtensionResponse{}
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, errors.New(fsmDomain.DefaultMessages.Error)
	}

	*machine = *res.FSM

	return res.Answers, nil
}

// GetAllExtensions returns all command functions in the extension as a list of strings
func (e *REST) GetAllExtensions() ([]string, error) {
	request, err := retryablehttp.NewRequest("GET", fmt.Sprintf("%s/extensions", e.URL), nil)
	if err != nil {
		return nil, err
	}
	if e.token != "" {
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", e.token))
	}
	resp, err := e.http.Do(request)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	var res []string
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Error(err)
		return nil, err
	}

	return res, nil
}
