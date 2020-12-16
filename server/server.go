package main

import (
	"flag"
	"fmt"
	"net"
	"net/rpc"
)

type Params struct {
	ImageWidth  int
	ImageHeight int
}

var (
	globalWorld [][]byte
	p           Params
)

type Broker struct{}

func (b *Broker) WorldTransfer(req WorldRequest, res *WorldResponse) (err error) {
	globalWorld = req.World
	p.ImageHeight = req.Height
	p.ImageWidth = req.Width
	res.Flag = true
	return
}

func (b *Broker) Subscribe(req SubRequest, res *SubResponse) (err error) {
	client, _ := rpc.Dial("tcp", req.FactoryAddress[0])

	calRequest := CalRequest{0, p.ImageHeight, 0, p.ImageWidth, globalWorld}
	calResponse := new(CalResponse)
	err2 := client.Call(req.FuncName, calRequest, calResponse)
	if err2 != nil {
		fmt.Println("Func of Subserver Return Error")
		fmt.Println(err2)
		fmt.Println("Closing subscriber thread.")
	}

	globalWorld = calResponse.World
	res.World = calResponse.World
	res.FlipCell = calResponse.FlipCell
	return
}

func main() {
	pAddr := flag.String("port", "8030", "Port to listen on")
	flag.Parse()
	rpc.Register(&Broker{})
	listener, _ := net.Listen("tcp", ":"+*pAddr)
	defer listener.Close()
	rpc.Accept(listener)
}
