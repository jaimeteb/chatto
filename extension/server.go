package extension

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jaimeteb/chatto/fsm"
	"github.com/jaimeteb/chatto/internal/logger"
	"github.com/jaimeteb/chatto/query"
	"github.com/jaimeteb/chatto/version"
	log "github.com/sirupsen/logrus"
)

var (
	invalidExtensionCommand = "extension command '%s' not found"
	invalidHTTPMethod       = "got method '%s', expected '%s'"
	invalidAuthToken        = "missing or incorrect authorization token"
)

// ExecuteCommandFuncRequest contains the instructions for executing a command function
type ExecuteCommandFuncRequest struct {
	FSM      *fsm.FSM        `json:"fsm"`
	Domain   *fsm.BaseDomain `json:"domain"`
	Command  string          `json:"command"`
	Question *query.Question `json:"question"`
	Channel  string          `json:"channel"`
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

// ErrorResponse is used when an error occurred processing a request
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
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

func httpError(w http.ResponseWriter, errorResponse ErrorResponse) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(errorResponse.Code)
	_ = json.NewEncoder(w).Encode(errorResponse)
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
		return fmt.Errorf(invalidExtensionCommand, req.Command)
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
	token                  string
}

// NewListenerREST creates a ListenerREST with command functions and a token
func NewListenerREST(registeredCommandFuncs RegisteredCommandFuncs, token string) *ListenerREST {
	return &ListenerREST{RegisteredCommandFuncs: registeredCommandFuncs, token: token}
}

// ServeREST serves the registered extension functions as a REST API
func ServeREST(registeredCommandFuncs RegisteredCommandFuncs) error {
	port := flag.Int("port", 8770, "Port to run extension server on")
	debug := flag.Bool("debug", false, "Enable debug logging.")

	sslKey := flag.String("ssl-key", "", "SSL key file for TLS secured server.")
	sslCert := flag.String("ssl-cert", "", "SSL certificate for TLS secured server.")

	token := flag.String("token", "", "Authorization token to be required by Chatto bot.")

	flag.Parse()

	logger.SetLogger(*debug)

	l := NewListenerREST(registeredCommandFuncs, *token)

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
	if l.token != "" {
		reqToken := r.Header.Get("Authorization")
		reqToken = strings.TrimPrefix(reqToken, "Bearer ")

		if l.token != reqToken {
			httpError(w, ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: invalidAuthToken,
			})
			return
		}
	}

	if r.Method != http.MethodPost {
		httpError(w, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf(invalidHTTPMethod, r.Method, http.MethodPost),
		})
		return
	}

	decoder := json.NewDecoder(r.Body)

	var req ExecuteCommandFuncRequest
	if err := decoder.Decode(&req); err != nil {
		httpError(w, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	commandFunc, ok := l.RegisteredCommandFuncs[req.Command]
	if !ok {
		httpError(w, ErrorResponse{
			Code:    http.StatusNotFound,
			Message: fmt.Sprintf(invalidExtensionCommand, req.Command),
		})
		return
	}
	res := commandFunc(&req)

	log.Debugf("ExecuteCommandFuncRequest:    %v,    %v", req.FSM, req.Command)
	log.Debugf("ExecuteCommandFuncResponse:    %v,    %v", *res.FSM, res.Answers)

	js, err := json.Marshal(res)
	if err != nil {
		httpError(w, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		httpError(w, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
}

// GetAllCommandFuncs returns all command functions in RegisteredCommandFuncs as a list of strings
func (l *ListenerREST) GetAllCommandFuncs(w http.ResponseWriter, r *http.Request) {
	if l.token != "" {
		reqToken := r.Header.Get("Authorization")
		reqToken = strings.TrimPrefix(reqToken, "Bearer ")

		if l.token != reqToken {
			httpError(w, ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: invalidAuthToken,
			})
			return
		}
	}

	if r.Method != http.MethodGet {
		httpError(w, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf(invalidHTTPMethod, r.Method, http.MethodGet),
		})
		return
	}

	js, err := json.Marshal(l.RegisteredCommandFuncs.Commands())
	if err != nil {
		httpError(w, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		httpError(w, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
}

// GetBuildVersion returns the current build version of the extension
func (l *ListenerREST) GetBuildVersion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpError(w, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf(invalidHTTPMethod, r.Method, http.MethodGet),
		})
		return
	}

	buildResponse := version.Build()

	js, err := json.Marshal(&buildResponse)
	if err != nil {
		httpError(w, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		httpError(w, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
}

// ExecuteCommandFuncResponseOption functions represent an option to build a new ExecuteCommandFuncResponse
type ExecuteCommandFuncResponseOption func(*ExecuteCommandFuncResponse)

// NewExecuteCommandFuncResponse creates a new ExecuteCommandFuncResponse based on the data from
// the ExecuteCommandFuncRequest that was sent to the extension command function
func (r *ExecuteCommandFuncRequest) NewExecuteCommandFuncResponse(opts ...ExecuteCommandFuncResponseOption) *ExecuteCommandFuncResponse {
	response := ExecuteCommandFuncResponse{
		FSM:     r.FSM,
		Answers: make([]query.Answer, 0),
	}
	for _, o := range opts {
		o(&response)
	}
	return &response
}

// WithAnswer appends a text and image answer to the response
func WithAnswer(text, image string) ExecuteCommandFuncResponseOption {
	return func(r *ExecuteCommandFuncResponse) {
		r.Answers = append(r.Answers, query.Answer{Text: text, Image: image})
	}
}

// WithTextAnswer appends a text answer to the response
func WithTextAnswer(text string) ExecuteCommandFuncResponseOption {
	return func(r *ExecuteCommandFuncResponse) {
		r.Answers = append(r.Answers, query.Answer{Text: text})
	}
}

// WithState sets a different state to the response's FSM
func WithState(state int) ExecuteCommandFuncResponseOption {
	return func(r *ExecuteCommandFuncResponse) {
		r.FSM.State = state
	}
}

// WithSlot sets a slot value to the response's FSM
func WithSlot(key, value string) ExecuteCommandFuncResponseOption {
	return func(r *ExecuteCommandFuncResponse) {
		r.FSM.Slots[key] = value
	}
}
