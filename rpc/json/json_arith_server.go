package main

import (
	"errors"
	"github.com/zacharychang/go-study/rpc/json/common"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type Arith int

func (t *Arith) Multiply(args *common.Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *Arith) Divide(args *common.Args, quo *common.Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}

func main() {
	arith := new(Arith)
	err := rpc.Register(arith)
	common.CheckErr(err)

	tcpAddr, err := net.ResolveTCPAddr("tcp", "localhost:1200")
	common.CheckErr(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	common.CheckErr(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		jsonrpc.ServeConn(conn)
	}
}
