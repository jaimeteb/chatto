package fsm

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/rpc"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

// Request struct for extension functions
type Request struct {
	FSM *FSM
	Req string
	Txt string
	Dom *DomainNoFuncs
}

// Response struct for extension functions
type Response struct {
	FSM *FSM
	Res string
}

// GetAllFuncsResponse struct for GetAllFuncs function
type GetAllFuncsResponse struct {
	Res []string
}

// ExtensionMap maps strings to functions to be used in extensions
type ExtensionMap map[string]func(*Request) *Response

// ListenerRPC contains the ExtensionMap to be served through RPC
type ListenerRPC struct {
	ExtensionMap ExtensionMap
}

// ListenerREST contains the ExtensionMap to be served through REST
type ListenerREST struct {
	ExtensionMap ExtensionMap
}

// GetFunc returns a requested extension function
func (l *ListenerRPC) GetFunc(req *Request, res *Response) error {
	extRes := l.ExtensionMap[req.Req](req)

	res.FSM = extRes.FSM
	res.Res = extRes.Res

	log.Debugf("Request:\t%v,\t%v", req.FSM, req.Req)
	log.Debugf("Response:\t%v,\t%v", *res.FSM, res.Res)
	return nil
}

// GetAllFuncs returns all functions registered in an ExtensionMap
func (l *ListenerRPC) GetAllFuncs(req *Request, res *GetAllFuncsResponse) error {
	allFuncs := make([]string, 0)
	for funcName := range l.ExtensionMap {
		allFuncs = append(allFuncs, funcName)
	}
	res.Res = allFuncs
	log.Debug(res)
	return nil
}

// GetFunc returns a requested extension function as a REST API
func (l *ListenerREST) GetFunc(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req Request
	if err := decoder.Decode(&req); err != nil {
		// log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res := l.ExtensionMap[req.Req](&req)

	log.Debugf("Request:\t%v,\t%v", req.FSM, req.Req)
	log.Debugf("Response:\t%v,\t%v", *res.FSM, res.Res)

	js, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// GetAllFuncs returns all functions registered in an ExtensionMap as a REST API
func (l *ListenerREST) GetAllFuncs(w http.ResponseWriter, r *http.Request) {
	allFuncs := make([]string, 0)
	for funcName := range l.ExtensionMap {
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

// ServeExtensionRPC serves the registered extension functions over RPC
func ServeExtensionRPC(extMap ExtensionMap) error {
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

	log.Infof("RPC extension server started. Using port %v\n", *port)
	rpc.Register(&ListenerRPC{ExtensionMap: extMap})
	rpc.Accept(inbound)
	return nil
}

// ServeExtensionREST serves the registered extension functions as a REST API
func ServeExtensionREST(extMap ExtensionMap) error {
	port := flag.Int("port", 8770, "Port to run extension server on")

	l := ListenerREST{ExtensionMap: extMap}

	r := mux.NewRouter()
	r.HandleFunc("/ext/get_func", l.GetFunc).Methods("POST")
	r.HandleFunc("/ext/get_all_funcs", l.GetAllFuncs).Methods("GET")

	log.Infof("REST extension server started. Using port %v\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *port), r))
	return nil
}
