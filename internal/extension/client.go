package extension

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/rpc"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/jaimeteb/chatto/extension"
	"github.com/jaimeteb/chatto/fsm"
	"github.com/jaimeteb/chatto/query"
	log "github.com/sirupsen/logrus"
)

// Extension is either a RPC or REST endpoint
type Extension interface {
	GetAllFuncs() ([]string, error)
	RunFunc(question *query.Question, extension string, fsmDomain *fsm.Domain, machine *fsm.FSM) ([]query.Answer, error)
}

// RPC is an RPC Client for extension functions
type RPC struct {
	Client *rpc.Client
}

// RunFunc runs an extension function over RPC
func (e *RPC) RunFunc(question *query.Question, ext string, fsmDomain *fsm.Domain, machine *fsm.FSM) ([]query.Answer, error) {
	req := extension.Request{
		FSM:       machine,
		Extension: ext,
		Question:  question,
		Domain:    fsmDomain.NoFuncs(),
	}

	res := extension.Response{}

	err := e.Client.Call("ListenerRPC.GetFunc", &req, &res)
	if err != nil {
		return nil, errors.New(fsmDomain.DefaultMessages.Error)
	}

	*machine = *res.FSM

	return res.Answers, nil
}

// GetAllFuncs retrieves all functions in extension
func (e *RPC) GetAllFuncs() ([]string, error) {
	res := new(extension.GetAllFuncsResponse)
	if err := e.Client.Call("ListenerRPC.GetAllFuncs", new(extension.Request), &res); err != nil {
		log.Error(err)
		return nil, err
	}

	return res.Funcs, nil
}

// REST is a REST API URL for extension functions
type REST struct {
	URL   string
	http  *retryablehttp.Client
	token string
}

// RunFunc runs an extension function over REST
func (e *REST) RunFunc(question *query.Question, ext string, fsmDomain *fsm.Domain, machine *fsm.FSM) ([]query.Answer, error) {
	req := extension.Request{
		FSM:       machine,
		Extension: ext,
		Question:  question,
		Domain:    fsmDomain.NoFuncs(),
	}

	jsonReq, err := json.Marshal(req)
	if err != nil {
		return nil, errors.New(fsmDomain.DefaultMessages.Error)
	}

	// TODO: if fail -> don't change states

	request, err := retryablehttp.NewRequest("POST", fmt.Sprintf("%s/ext/get_func", e.URL), bytes.NewBuffer(jsonReq))
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

	defer resp.Body.Close()
	res := extension.Response{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, errors.New(fsmDomain.DefaultMessages.Error)
	}

	*machine = *res.FSM

	return res.Answers, nil
}

// GetAllFuncs retrieves all functions in extension
func (e *REST) GetAllFuncs() ([]string, error) {
	request, err := retryablehttp.NewRequest("GET", fmt.Sprintf("%s/ext/get_all_funcs", e.URL), nil)
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

	defer resp.Body.Close()
	var res []string
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return res, nil
}
