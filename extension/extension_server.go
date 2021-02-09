package extension

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/rpc"

	"github.com/gorilla/mux"
	"github.com/jaimeteb/chatto/fsm"
	"github.com/jaimeteb/chatto/logger"
	"github.com/jaimeteb/chatto/query"
	log "github.com/sirupsen/logrus"
)

// Request for an extension function
type Request struct {
	FSM       *fsm.FSM        `json:"fsm"`
	Extension string          `json:"extension"`
	Question  *query.Question `json:"question"`
	DB        *fsm.BaseDB     `json:"db"`
}

// Response from an extension function
type Response struct {
	FSM     *fsm.FSM       `json:"fsm"`
	Answers []query.Answer `json:"answers"`
}

// GetAllFuncsResponse contains a list of all registered functions
type GetAllFuncsResponse struct {
	Funcs []string
}

// RegisteredFuncs maps strings to functions to be used in extensions
type RegisteredFuncs map[string]func(*Request) *Response

// ListenerRPC contains the RegisteredFuncs to be served through RPC
type ListenerRPC struct {
	RegisteredFuncs RegisteredFuncs
}

// ListenerREST contains the RegisteredFuncs to be served through REST
type ListenerREST struct {
	RegisteredFuncs RegisteredFuncs
}

// GetFunc returns a requested extension function
func (l *ListenerRPC) GetFunc(req *Request, res *Response) error {
	extRes := l.RegisteredFuncs[req.Extension](req)

	res.FSM = extRes.FSM
	res.Answers = extRes.Answers

	log.Debugf("Request:\t%v,\t%v", req.FSM, req.Extension)
	log.Debugf("Response:\t%v,\t%v", *res.FSM, res.Answers)

	return nil
}

// GetAllFuncs returns all functions registered in an RegisteredFuncs
func (l *ListenerRPC) GetAllFuncs(req *Request, res *GetAllFuncsResponse) error {
	allFuncs := make([]string, 0)
	for funcName := range l.RegisteredFuncs {
		allFuncs = append(allFuncs, funcName)
	}
	res.Funcs = allFuncs
	log.Debug(res)
	return nil
}

// GetFunc returns a requested extension function as a REST API
func (l *ListenerREST) GetFunc(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req Request
	if err := decoder.Decode(&req); err != nil {
		// log.Info(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res := l.RegisteredFuncs[req.Extension](&req)

	log.Debugf("Request:\t%v,\t%v", req.FSM, req.Extension)
	log.Debugf("Response:\t%v,\t%v", *res.FSM, res.Answers)

	js, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// GetAllFuncs returns all functions registered in an RegisteredFuncs as a REST API
func (l *ListenerREST) GetAllFuncs(w http.ResponseWriter, r *http.Request) {
	allFuncs := make([]string, 0)
	for funcName := range l.RegisteredFuncs {
		allFuncs = append(allFuncs, funcName)
	}

	js, err := json.Marshal(allFuncs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// ServeRPC serves the registered extension functions over RPC
func ServeRPC(extMap RegisteredFuncs) error {
	logger.SetLogger()

	host := flag.String("host", "0.0.0.0", "Host to run extension server on")
	port := flag.Int("port", 8770, "Port to run extension server on")
	flag.Parse()

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%v:%v", *host, *port))
	if err != nil {
		log.Error(err)
		return err
	}

	inbound, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Infof("RPC extension server started. Using port %v", *port)
	rpc.Register(&ListenerRPC{RegisteredFuncs: extMap})
	rpc.Accept(inbound)
	return nil
}

// ServeREST serves the registered extension functions as a REST API
func ServeREST(extMap RegisteredFuncs) error {
	logger.SetLogger()

	port := flag.Int("port", 8770, "Port to run extension server on")
	flag.Parse()

	l := ListenerREST{RegisteredFuncs: extMap}

	r := mux.NewRouter()
	r.HandleFunc("/ext/get_func", l.GetFunc).Methods("POST")
	r.HandleFunc("/ext/get_all_funcs", l.GetAllFuncs).Methods("GET")

	log.Infof("REST extension server started. Using port %v", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *port), r))
	return nil
}
