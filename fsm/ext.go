package fsm

import (
	"log"
	"net/rpc"
)

// // Extension interface
// type Extension interface {
// 	GetFunc(string) func(*FSM, *Domain, string) interface{}
// 	GetAllFuncs() []string
// }

// // FuncMap maps function names to functions
// type FuncMap map[string]func(*FSM, *Domain, string) interface{}

// // GetFunc gets a function from the function map
// func (fm FuncMap) GetFunc(action string) func(*FSM, *Domain, string) interface{} {
// 	if _, ok := fm[action]; ok {
// 		return fm[action]
// 	}
// 	return func(*FSM, *Domain, string) interface{} {
// 		return nil
// 	}
// }

// // GetAllFuncs retreives all functions in function map
// func (fm FuncMap) GetAllFuncs() []string {
// 	allFuncs := make([]string, 0)
// 	for funcName := range fm {
// 		allFuncs = append(allFuncs, funcName)
// 	}
// 	return allFuncs
// }

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

type ExtensionMap map[string]func(*Request) *Response

type Extension *rpc.Client

type Request struct {
	FSM FSM
	Req string
}

type Response struct {
	FSM FSM
	Res string
}

type GetAllFuncsResponse struct {
	Res []string
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
