package extension

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/rpc"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/jaimeteb/chatto/fsm"
	"github.com/jaimeteb/chatto/query"
	log "github.com/sirupsen/logrus"
)

// Extension is either a RPC or REST endpoint
type Extension interface {
	GetAllFuncs() ([]string, error)
	RunExtFunc(question *query.Question, extension string, fsmDomain *fsm.Domain, machine *fsm.FSM) ([]query.Answer, error)
}

// RPC is an RPC Client for extension functions
type RPC struct {
	Client *rpc.Client
}

// RunExtFunc runs an extension function over RPC
func (e *RPC) RunExtFunc(question *query.Question, extension string, fsmDomain *fsm.Domain, machine *fsm.FSM) ([]query.Answer, error) {
	req := Request{
		FSM:       machine,
		Extension: extension,
		Question:  question,
		Domain:    fsmDomain.NoFuncs(),
	}

	res := Response{}

	err := e.Client.Call("ListenerRPC.GetFunc", &req, &res)
	if err != nil {
		return nil, errors.New(fsmDomain.DefaultMessages.Error)
	}

	*machine = *res.FSM

	return res.Answers, nil
}

// GetAllFuncs retrieves all functions in extension
func (e *RPC) GetAllFuncs() ([]string, error) {
	res := new(GetAllFuncsResponse)
	if err := e.Client.Call("ListenerRPC.GetAllFuncs", new(Request), &res); err != nil {
		log.Error(err)
		return nil, err
	}

	return res.Funcs, nil
}

// REST is a REST API URL for extension functions
type REST struct {
	URL  string
	http *retryablehttp.Client
}

// RunExtFunc runs an extension function over REST
func (e *REST) RunExtFunc(question *query.Question, extension string, fsmDomain *fsm.Domain, machine *fsm.FSM) ([]query.Answer, error) {
	req := Request{
		FSM:       machine,
		Extension: extension,
		Question:  question,
		Domain:    fsmDomain.NoFuncs(),
	}

	jsonReq, err := json.Marshal(req)
	if err != nil {
		return nil, errors.New(fsmDomain.DefaultMessages.Error)
	}

	// TODO: if fail -> don't change states

	resp, err := e.http.Post(fmt.Sprintf("%s/ext/get_func", e.URL), "application/json", bytes.NewBuffer(jsonReq))
	if err != nil {
		return nil, errors.New(fsmDomain.DefaultMessages.Error)
	}

	defer resp.Body.Close()
	res := Response{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, errors.New(fsmDomain.DefaultMessages.Error)
	}

	*machine = *res.FSM

	return res.Answers, nil
}

// GetAllFuncs retrieves all functions in extension
func (e *REST) GetAllFuncs() ([]string, error) {
	resp, err := e.http.Get(fmt.Sprintf("%s/ext/get_all_funcs", e.URL))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	defer resp.Body.Close()
	var res []string
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Error(err)
		return nil, err
	}

	return res, nil
}
