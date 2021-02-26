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
	"github.com/jaimeteb/chatto/internal/logger"
	"github.com/jaimeteb/chatto/query"
	"github.com/jaimeteb/chatto/version"
	log "github.com/sirupsen/logrus"
)

var (
	extensionCommandNotFound = "extension command %s not found"
	invalidHTTPMethod        = "got method %s, expected %s"
)

// ExecuteCommandFuncRequest contains the instructions for executing a command function
type ExecuteCommandFuncRequest struct {
	FSM      *fsm.FSM        `json:"fsm"`
	Domain   *fsm.BaseDomain `json:"domain"`
	Command  string          `json:"command"`
	Question *query.Question `json:"question"`
}

// ExecuteCommandFuncResponse contains the result of executing a command function
type ExecuteCommandFuncResponse struct {
	FSM     *fsm.FSM       `json:"fsm"`
	Answers []query.Answer `json:"answers"`
}

// GetAllCommandFuncsRequest is empty for now to match RPC interface. Maybe later
// we will use it for filtering/searching commands
type GetAllCommandFuncsRequest struct {
}

// GetAllCommandFuncsResponse contains a list of all registered command functions
type GetAllCommandFuncsResponse struct {
	Commands []string
}

// GetBuildVersionRequest is empty for now to match RPC interface.
type GetBuildVersionRequest struct {
}

// RegisteredCommandFuncs maps commands to functions which are executed by extension servers
type RegisteredCommandFuncs map[string]func(*ExecuteCommandFuncRequest) *ExecuteCommandFuncResponse

// Commands returns a list of all registered function command names
func (r *RegisteredCommandFuncs) Commands() []string {
	if r == nil {
		return []string{}
	}

	commands := make([]string, 0, len(*r))
	for command := range *r {
		commands = append(commands, command)
	}

	return commands
}

// ListenerRPC contains the RegisteredCommandFuncs to be served through RPC
type ListenerRPC struct {
	RegisteredCommandFuncs RegisteredCommandFuncs
}

// ServeRPC serves the registered extension functions over RPC
func ServeRPC(registeredCommandFuncs RegisteredCommandFuncs) error {
	host := flag.String("host", "0.0.0.0", "Host to run extension server on")
	port := flag.Int("port", 8770, "Port to run extension server on")
	debug := flag.Bool("debug", false, "Enable debug logging.")

	flag.Parse()

	logger.SetLogger(*debug)

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		log.Error(err)
		return err
	}

	inbound, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Infof("RPC extension server started. Using port %d", *port)
	err = rpc.Register(&ListenerRPC{RegisteredCommandFuncs: registeredCommandFuncs})
	if err != nil {
		log.Error(err)
		return err
	}

	rpc.Accept(inbound)

	return nil
}

// ExecuteCommandFunc runs the requested command function and returns the response
func (l *ListenerRPC) ExecuteCommandFunc(req *ExecuteCommandFuncRequest, res *ExecuteCommandFuncResponse) error {
	command, ok := l.RegisteredCommandFuncs[req.Command]
	if !ok {
		return fmt.Errorf(extensionCommandNotFound, req.Command)
	}
	commandRes := command(req)

	res.FSM = commandRes.FSM
	res.Answers = commandRes.Answers

	log.Debugf("ExecuteCommandFuncRequest:    %v,    %v", req.FSM, req.Command)
	log.Debugf("ExecuteCommandFuncResponse:    %v,    %v", *res.FSM, res.Answers)

	return nil
}

// GetAllCommandFuncs returns all functions registered in the RegisteredCommandFuncs map
func (l *ListenerRPC) GetAllCommandFuncs(_ *GetAllCommandFuncsRequest, res *GetAllCommandFuncsResponse) error {
	res.Commands = l.RegisteredCommandFuncs.Commands()
	log.Debug(res)
	return nil
}

// GetBuildVersion returns the current build version of the extension
func (l *ListenerRPC) GetBuildVersion(_ *GetBuildVersionRequest, res *version.BuildResponse) error {
	buildResponse := version.Build()
	res.Version = buildResponse.Version
	res.Commit = buildResponse.Commit
	res.BuiltAt = buildResponse.BuiltAt
	res.BuiltBy = buildResponse.BuiltBy
	return nil
}

// ListenerREST contains the RegisteredCommandFuncs to be served through REST
type ListenerREST struct {
	RegisteredCommandFuncs RegisteredCommandFuncs
}

// ServeREST serves the registered extension functions as a REST API
func ServeREST(registeredCommandFuncs RegisteredCommandFuncs) error {
	port := flag.Int("port", 8770, "Port to run extension server on")
	debug := flag.Bool("debug", false, "Enable debug logging.")

	sslKey := flag.String("ssl-key", "", "SSL key file for TLS secured server.")
	sslCert := flag.String("ssl-cert", "", "SSL certificate for TLS secured server.")

	flag.Parse()

	logger.SetLogger(*debug)

	l := ListenerREST{RegisteredCommandFuncs: registeredCommandFuncs}

	r := mux.NewRouter()
	r.HandleFunc("/ext/command", l.ExecuteCommandFunc).Methods("POST")
	r.HandleFunc("/ext/commands", l.GetAllCommandFuncs).Methods("GET")
	r.HandleFunc("/ext/version", l.GetBuildVersion).Methods("GET")

	if *sslKey != "" && *sslCert != "" {
		log.Infof("REST extension server started with TLS. Using port %d", *port)
		log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", *port), *sslCert, *sslKey, r))
	} else {
		log.Infof("REST extension server started. Using port %d", *port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), r))
	}

	return nil
}

// ExecuteCommandFunc runs the requested command function and returns the response
func (l *ListenerREST) ExecuteCommandFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, fmt.Sprintf(invalidHTTPMethod, r.Method, http.MethodPost), http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)

	var req ExecuteCommandFuncRequest
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	commandFunc, ok := l.RegisteredCommandFuncs[req.Command]
	if !ok {
		http.Error(w, fmt.Sprintf(extensionCommandNotFound, req.Command), http.StatusBadRequest)
		return
	}
	res := commandFunc(&req)

	log.Debugf("ExecuteCommandFuncRequest:    %v,    %v", req.FSM, req.Command)
	log.Debugf("ExecuteCommandFuncResponse:    %v,    %v", *res.FSM, res.Answers)

	js, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetAllCommandFuncs returns all command functions in RegisteredCommandFuncs as a list of strings
func (l *ListenerREST) GetAllCommandFuncs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, fmt.Sprintf(invalidHTTPMethod, r.Method, http.MethodGet), http.StatusBadRequest)
		return
	}

	js, err := json.Marshal(l.RegisteredCommandFuncs.Commands())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetBuildVersion returns the current build version of the extension
func (l *ListenerREST) GetBuildVersion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, fmt.Sprintf(invalidHTTPMethod, r.Method, http.MethodGet), http.StatusBadRequest)
		return
	}

	buildResponse := version.Build()

	js, err := json.Marshal(&buildResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
