package fsm

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
)

// BuildPlugin builds the extension code as a plugin
// func BuildPlugin(path *string) error {
// 	buildGo := "go"
// 	buildArgs := []string{
// 		"build",
// 		"-buildmode=plugin",
// 		"-o",
// 		filepath.Join(*path, "ext/ext.so"),
// 		filepath.Join(*path, "ext/ext.go"),
// 	}

// 	cmd := exec.Command(buildGo, buildArgs...)
// 	_, err := cmd.Output()

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// ExtensionMap maps strings to functions to be
// used in extensions
type ExtensionMap map[string]func(*Request) *Response

// Listener contains the ExtensionMap to be served through RPC
type Listener struct {
	ExtensionMap ExtensionMap
}

// GetFunc returns a requested extension function
func (l *Listener) GetFunc(req *Request, res *Response) error {
	extRes := l.ExtensionMap[req.Req](req)

	res.FSM = extRes.FSM
	res.Res = extRes.Res

	log.Printf("Request:\t%v,\t%v", req.FSM, req.Req)
	log.Printf("Response:\t%v,\t%v", *res.FSM, res.Res)
	return nil
}

// GetAllFuncs returns all functions registered in an ExtensionMap
func (l *Listener) GetAllFuncs(req *Request, res *GetAllFuncsResponse) error {
	allFuncs := make([]string, 0)
	for funcName := range l.ExtensionMap {
		allFuncs = append(allFuncs, funcName)
	}
	res.Res = allFuncs
	log.Println(res)
	return nil
}

// Extension is an RPC Client that serves extension functions
type Extension *rpc.Client

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

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

var extensionHost = getEnv("EXTENSION_HOST", "localhost")
var extensionPort = getEnv("EXTENSION_PORT", "42586")
var extensionAddr = fmt.Sprintf("0.0.0.0:%v", extensionPort)
var extensionDial = fmt.Sprintf("%v:%v", extensionHost, extensionPort)

// ServeExtension serves the registered extension functions
func ServeExtension(extMap ExtensionMap) error {
	addr, err := net.ResolveTCPAddr("tcp", extensionAddr)
	if err != nil {
		log.Println(err)
		return err
	}

	inbound, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Printf("Extension server started. Using port %v\n", extensionPort)
	rpc.Register(&Listener{ExtensionMap: extMap})
	rpc.Accept(inbound)
	return nil
}

// LoadExtension creates an extension
func LoadExtension(path *string) (Extension, error) {
	loadExtErr := func(err error) (Extension, error) {
		log.Println("Error while loading extensions: ", err.Error())
		return nil, err
	}

	client, err := rpc.Dial("tcp", extensionDial)
	if err != nil {
		return loadExtErr(err)
	}

	res := new(GetAllFuncsResponse)
	err = client.Call("Listener.GetAllFuncs", new(Request), &res)
	if err != nil {
		return loadExtErr(err)
	}

	log.Println("Loaded extensions:")
	for i, fun := range res.Res {
		log.Printf("%v\t%v\n", i, fun)
	}

	return client, nil
}
