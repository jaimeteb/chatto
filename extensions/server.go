package extensions

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
	invalidExtension  = "extension '%s' not found"
	invalidHTTPMethod = "got method '%s', expected '%s'"
	invalidAuthToken  = "missing or incorrect authorization token"
)

// ExecuteExtensionRequest contains the instructions for executing a command function
type ExecuteExtensionRequest struct {
	FSM       *fsm.FSM        `json:"fsm"`
	Domain    *fsm.BaseDomain `json:"domain"`
	Extension string          `json:"extension"`
	Question  *query.Question `json:"question"`
	Channel   string          `json:"channel"`
	Command   string          `json:"command"`
}

// ExecuteExtensionResponse contains the result of executing a command function
type ExecuteExtensionResponse struct {
	FSM     *fsm.FSM       `json:"fsm"`
	Answers []query.Answer `json:"answers"`
}

// GetAllExtensionsRequest is empty for now to match RPC interface. Maybe later
// we will use it for filtering/searching commands
type GetAllExtensionsRequest struct {
}

// GetAllExtensionsResponse contains a list of all registered command functions
type GetAllExtensionsResponse struct {
	Extensions []string
}

// GetBuildVersionRequest is empty for now to match RPC interface.
type GetBuildVersionRequest struct {
}

// ErrorResponse is used when an error occurred processing a request
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// RegisteredExtensions maps commands to functions which are executed by extension servers
type RegisteredExtensions map[string]func(*ExecuteExtensionRequest) *ExecuteExtensionResponse

// Get returns a list of all registered function command names
func (r *RegisteredExtensions) Get() []string {
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

// ListenerRPC contains the RegisteredExtensions to be served through RPC
type ListenerRPC struct {
	RegisteredExtensions RegisteredExtensions
}

// ServeRPC serves the registered extension functions over RPC
func ServeRPC(registeredExtensions RegisteredExtensions) error {
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
	err = rpc.Register(&ListenerRPC{RegisteredExtensions: registeredExtensions})
	if err != nil {
		log.Error(err)
		return err
	}

	rpc.Accept(inbound)

	return nil
}

// ExecuteExtension runs the requested command function and returns the response
func (l *ListenerRPC) ExecuteExtension(req *ExecuteExtensionRequest, res *ExecuteExtensionResponse) error {
	command, ok := l.RegisteredExtensions[req.Extension]
	if !ok {
		return fmt.Errorf(invalidExtension, req.Extension)
	}
	commandRes := command(req)

	res.FSM = commandRes.FSM
	res.Answers = commandRes.Answers

	log.Debugf("ExecuteExtensionRequest:    %v,    %v", req.FSM, req.Extension)
	log.Debugf("ExecuteExtensionResponse:    %v,    %v", *res.FSM, res.Answers)

	return nil
}

// GetAllExtensions returns all functions registered in the RegisteredExtensions map
func (l *ListenerRPC) GetAllExtensions(_ *GetAllExtensionsRequest, res *GetAllExtensionsResponse) error {
	res.Extensions = l.RegisteredExtensions.Get()
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

// ListenerREST contains the RegisteredExtensions to be served through REST
type ListenerREST struct {
	RegisteredExtensions RegisteredExtensions
	token                string
}

// NewListenerREST creates a ListenerREST with command functions and a token
func NewListenerREST(registeredExtensions RegisteredExtensions, token string) *ListenerREST {
	return &ListenerREST{RegisteredExtensions: registeredExtensions, token: token}
}

// ServeREST serves the registered extension functions as a REST API
func ServeREST(registeredExtensions RegisteredExtensions) error {
	port := flag.Int("port", 8770, "Port to run extension server on")
	debug := flag.Bool("debug", false, "Enable debug logging.")

	sslKey := flag.String("ssl-key", "", "SSL key file for TLS secured server.")
	sslCert := flag.String("ssl-cert", "", "SSL certificate for TLS secured server.")

	token := flag.String("token", "", "Authorization token to be required by Chatto bot.")

	flag.Parse()

	logger.SetLogger(*debug)

	l := NewListenerREST(registeredExtensions, *token)

	r := mux.NewRouter()
	r.HandleFunc("/extension", l.ExecuteExtension).Methods("POST")
	r.HandleFunc("/extensions", l.GetAllExtensions).Methods("GET")
	r.HandleFunc("/version", l.GetBuildVersion).Methods("GET")

	if *sslKey != "" && *sslCert != "" {
		log.Infof("REST extension server started with TLS. Using port %d", *port)
		log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", *port), *sslCert, *sslKey, r))
	} else {
		log.Infof("REST extension server started. Using port %d", *port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), r))
	}

	return nil
}

// ExecuteExtension runs the requested command function and returns the response
func (l *ListenerREST) ExecuteExtension(w http.ResponseWriter, r *http.Request) {
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

	var req ExecuteExtensionRequest
	if err := decoder.Decode(&req); err != nil {
		httpError(w, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	commandFunc, ok := l.RegisteredExtensions[req.Extension]
	if !ok {
		httpError(w, ErrorResponse{
			Code:    http.StatusNotFound,
			Message: fmt.Sprintf(invalidExtension, req.Extension),
		})
		return
	}
	res := commandFunc(&req)

	log.Debugf("ExecuteExtensionRequest:    %v,    %v", req.FSM, req.Extension)
	log.Debugf("ExecuteExtensionResponse:    %v,    %v", *res.FSM, res.Answers)

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

// GetAllExtensions returns all command functions in RegisteredExtensions as a list of strings
func (l *ListenerREST) GetAllExtensions(w http.ResponseWriter, r *http.Request) {
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

	js, err := json.Marshal(l.RegisteredExtensions.Get())
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
