package fsm

import (
	"log"
	"net"
	"net/rpc"
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

	res.FSM = req.FSM
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
	Dom Domain
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

// ServeExtension serves the registered extension functions
func ServeExtension(extMap ExtensionMap) error {
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:42586")
	if err != nil {
		log.Println(err)
		return err
	}

	inbound, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Println(err)
		return err
	}

	rpc.Register(&Listener{ExtensionMap: extMap})
	rpc.Accept(inbound)
	return nil
}

// LoadExtension creates an extension
func LoadExtension(path *string) (Extension, error) {
	loadExtErr := func(err error) (Extension, error) {
		log.Println("Error while loading extensions: ", err.Error())
		log.Println("Using bot without extensions.")
		return nil, err
	}

	client, err := rpc.Dial("tcp", "localhost:42586")
	if err != nil {
		return loadExtErr(err)
	}

	res := new(GetAllFuncsResponse)
	err = client.Call("Listener.GetAllFuncs", new(Request), &res)
	if err != nil {
		return loadExtErr(err)
	}

	log.Println(res.Res)

	return client, nil
}
