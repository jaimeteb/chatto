package extension

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/rpc"

	"github.com/jaimeteb/chatto/fsm"
	"github.com/jaimeteb/chatto/query"
	log "github.com/sirupsen/logrus"
)

// Config options for an extension function
type Config struct {
	Type string `mapstructure:"type"`
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	URL  string `mapstructure:"url"`
}

// Extension is either a RPC or REST endpoint
type Extension interface {
	GetAllFuncs() ([]string, error)
	RunExtFunc(question *query.Question, extension string, db *fsm.DB, machine *fsm.FSM) ([]query.Answer, error)
}

// RPC is an RPC Client for extension functions
type RPC struct {
	Client *rpc.Client
}

// RunExtFunc runs an extension function over RPC
func (e *RPC) RunExtFunc(question *query.Question, extension string, db *fsm.DB, machine *fsm.FSM) ([]query.Answer, error) {
	req := Request{
		FSM:       machine,
		Extension: extension,
		Question:  question,
		DB:        db.NoFuncs(),
	}

	res := Response{}

	err := e.Client.Call("ListenerRPC.GetFunc", &req, &res)
	if err != nil {
		log.Error(err)
		return query.Answers(db.DefaultMessages.Error), err
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
	URL string
}

// RunExtFunc runs an extension function over REST
func (e *REST) RunExtFunc(question *query.Question, extension string, db *fsm.DB, machine *fsm.FSM) ([]query.Answer, error) {
	req := Request{
		FSM:       machine,
		Extension: extension,
		Question:  question,
		DB:        db.NoFuncs(),
	}

	jsonReq, err := json.Marshal(req)
	if err != nil {
		log.Error(err)
		return query.Answers(db.DefaultMessages.Error), err
	}

	// TODO: if fail -> don't change states

	resp, err := http.Post(fmt.Sprintf("%v/ext/get_func", e.URL), "application/json", bytes.NewBuffer(jsonReq))
	if err != nil {
		log.Error(err)
		return query.Answers(db.DefaultMessages.Error), err
	}

	defer resp.Body.Close()
	res := Response{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Error(err)
		return query.Answers(db.DefaultMessages.Error), err
	}

	*machine = *res.FSM

	return res.Answers, nil
}

// GetAllFuncs retrieves all functions in extension
func (e *REST) GetAllFuncs() ([]string, error) {
	resp, err := http.Get(fmt.Sprintf("%v/ext/get_all_funcs", e.URL))
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

// LoadExtensions loads the extensions configuration and connects to the server
func LoadExtensions(extCfg Config) (Extension, error) {
	var extension Extension

	switch extCfg.Type {
	case "RPC":
		client, err := rpc.Dial("tcp", fmt.Sprintf("%v:%v", extCfg.Host, extCfg.Port))
		if err != nil {
			return nil, err
		}

		rpcExtension := &RPC{client}

		allFuncs, err := rpcExtension.GetAllFuncs()
		if err != nil {
			return nil, err
		}

		log.Info("Loaded extensions (RPC):")
		for i, fun := range allFuncs {
			log.Infof("%v\t%v", i, fun)
		}

		extension = rpcExtension
	case "REST":
		restExtention := &REST{extCfg.URL}

		allFuncs, err := restExtention.GetAllFuncs()
		if err != nil {
			return nil, err
		}

		log.Info("Loaded extensions (REST):")
		for i, fun := range allFuncs {
			log.Infof("%v\t%v", i, fun)
		}

		extension = restExtention
	}

	if extension == nil {
		log.Info("Using bot without extensions.")
	}

	return extension, nil
}
