package extension

import (
	"encoding/json"
	"errors"
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
	log "github.com/sirupsen/logrus"
)

var (
	// ErrExtensionNotFound happens when an extension is requested but
	// is not found in the extension server
	ErrExtensionNotFound = errors.New("extension not found")
	// ErrExtensionUnauthorized happens when the server requires a token
	// but it is missing or is incorrect in the request
	ErrExtensionUnauthorized = errors.New("missing or incorrect token")
)

// Request for an extension function
type Request struct {
	FSM       *fsm.FSM        `json:"fsm"`
	Extension string          `json:"extension"`
	Question  *query.Question `json:"question"`
	Domain    *fsm.BaseDomain `json:"domain"`
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

// RegisteredCommandFuncs maps bot commands to functions to be used in extensions
type RegisteredCommandFuncs map[string]func(*Request) *Response

// ListenerRPC contains the RegisteredCommandFuncs to be served through RPC
type ListenerRPC struct {
	RegisteredCommandFuncs RegisteredCommandFuncs
}

// ServeRPC serves the registered extension functions over RPC
func ServeRPC(RegisteredCommandFuncs RegisteredCommandFuncs) error {
	host := flag.String("host", "0.0.0.0", "Host to run extension server on")
	port := flag.Int("port", 8770, "Port to run extension server on")
	debug := flag.Bool("debug", false, "Enable debug logging.")
	flag.Parse()

	logger.SetLogger(*debug)

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
	err = rpc.Register(&ListenerRPC{RegisteredCommandFuncs: RegisteredCommandFuncs})
	if err != nil {
		log.Error(err)
		return err
	}

	rpc.Accept(inbound)

	return nil
}

// GetFunc returns a requested extension function
func (l *ListenerRPC) GetFunc(req *Request, res *Response) error {
	extFunc, ok := l.RegisteredCommandFuncs[req.Extension]
	if !ok {
		return ErrExtensionNotFound
	}
	extRes := extFunc(req)

	res.FSM = extRes.FSM
	res.Answers = extRes.Answers

	log.Debugf("Request:    %v,    %v", req.FSM, req.Extension)
	log.Debugf("Response:    %v,    %v", *res.FSM, res.Answers)

	return nil
}

// GetAllFuncs returns all functions registered in an RegisteredCommandFuncs
func (l *ListenerRPC) GetAllFuncs(req *Request, res *GetAllFuncsResponse) error {
	allFuncs := make([]string, 0)
	for funcName := range l.RegisteredCommandFuncs {
		allFuncs = append(allFuncs, funcName)
	}
	res.Funcs = allFuncs
	log.Debug(res)
	return nil
}

// ListenerREST contains the RegisteredCommandFuncs to be served through REST
type ListenerREST struct {
	RegisteredCommandFuncs RegisteredCommandFuncs
	token                  string
}

// ServeREST serves the registered extension functions as a REST API
func ServeREST(RegisteredCommandFuncs RegisteredCommandFuncs) error {
	port := flag.Int("port", 8770, "Port to run extension server on")
	debug := flag.Bool("debug", false, "Enable debug logging.")

	sslKey := flag.String("ssl-key", "", "SSL key file for TLS secured server.")
	sslCert := flag.String("ssl-cert", "", "SSL certificate for TLS secured server.")

	token := flag.String("token", "", "Authorization token to be required by Chatto bot.")

	flag.Parse()

	logger.SetLogger(*debug)

	l := ListenerREST{RegisteredCommandFuncs: RegisteredCommandFuncs, token: *token}

	r := mux.NewRouter()
	r.HandleFunc("/ext/get_func", l.GetFunc).Methods("POST")
	r.HandleFunc("/ext/get_all_funcs", l.GetAllFuncs).Methods("GET")

	if *sslKey != "" && *sslCert != "" {
		log.Infof("REST extension server started with TLS. Using port %v", *port)
		log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", *port), *sslCert, *sslKey, r))
	} else {
		log.Infof("REST extension server started. Using port %v", *port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), r))
	}

	return nil
}

// GetFunc returns a requested extension function as a REST API
func (l *ListenerREST) GetFunc(w http.ResponseWriter, r *http.Request) {
	if l.token != "" {
		reqToken := r.Header.Get("Authorization")
		reqToken = strings.TrimPrefix(reqToken, "Bearer ")

		if l.token != reqToken {
			http.Error(w, ErrExtensionUnauthorized.Error(), http.StatusUnauthorized)
			return
		}
	}

	decoder := json.NewDecoder(r.Body)

	var req Request
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	extFunc, ok := l.RegisteredCommandFuncs[req.Extension]
	if !ok {
		http.Error(w, ErrExtensionNotFound.Error(), http.StatusBadRequest)
		return
	}
	res := extFunc(&req)

	log.Debugf("Request:    %v,    %v", req.FSM, req.Extension)
	log.Debugf("Response:    %v,    %v", *res.FSM, res.Answers)

	js, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

// GetAllFuncs returns all functions registered in an RegisteredCommandFuncs as a REST API
func (l *ListenerREST) GetAllFuncs(w http.ResponseWriter, r *http.Request) {
	if l.token != "" {
		reqToken := r.Header.Get("Authorization")
		reqToken = strings.TrimPrefix(reqToken, "Bearer ")

		if l.token != reqToken {
			http.Error(w, ErrExtensionUnauthorized.Error(), http.StatusUnauthorized)
			return
		}
	}

	allFuncs := make([]string, 0)
	for funcName := range l.RegisteredCommandFuncs {
		allFuncs = append(allFuncs, funcName)
	}

	js, err := json.Marshal(allFuncs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
