package main

import (
	"fmt"
	"github.com/zacharychang/go-study/rpc/json/common"
	"net/rpc/jsonrpc"
)

func main() {
	svc := "localhost:1200"

	client, err := jsonrpc.Dial("tcp", svc)
	common.CheckErr(err)

	args := common.Args{13, 4}

	var reply int
	err = client.Call("Arith.Multiply", args, &reply)
	common.CheckErr(err)
	fmt.Printf("Arith: %d*%d=%d\n", args.A, args.B, reply)

	var quot common.Quotient
	err = client.Call("Arith.Divide", args, &quot)
	common.CheckErr(err)
	fmt.Printf("Arith: %d/%d=%d with %d left\n", args.A, args.B, quot.Quo, quot.Rem)
}
