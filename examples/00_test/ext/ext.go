package main

import (
	"log"
	"net"
	"net/rpc"

	"github.com/jaimeteb/chatto/fsm"
)

func greetFunc(req *fsm.Request) (res *fsm.Response) {
	return &fsm.Response{
		FSM: req.FSM,
		Res: "Hello Universe",
	}
}

var Ext = fsm.ExtensionMap{
	"ext_any": greetFunc,
}

type Listener int

func (l *Listener) GetFunc(req *fsm.Request, res *fsm.Response) error {
	res = Ext[req.Req](req)
	return nil
}

func (l *Listener) GetAllFuncs(req *fsm.Request, res *fsm.GetAllFuncsResponse) error {
	allFuncs := make([]string, 0)
	for funcName := range Ext {
		allFuncs = append(allFuncs, funcName)
	}
	res.Res = allFuncs
	log.Println(res)
	return nil
}

func main() {
	addy, err := net.ResolveTCPAddr("tcp", "0.0.0.0:42586")
	if err != nil {
		log.Fatal(err)
	}

	inbound, err := net.ListenTCP("tcp", addy)
	if err != nil {
		log.Fatal(err)
	}

	rpc.Register(new(Listener))
	rpc.Accept(inbound)
}
