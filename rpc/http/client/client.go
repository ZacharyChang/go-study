package main

import (
	"fmt"
	"log"
	"net/rpc"

	"github.com/zacharychang/go-study/rpc/http/common"
)

var (
	addr = "localhost:1234"
)

func callMultiply(a int, b int) {
	client, err := rpc.DialHTTP("tcp", addr)
	defer client.Close()
	if err != nil {
		log.Fatal("dialing:", err)
	}
	args := &common.Args{a, b}
	var reply int
	err = client.Call("Arith.Multiply", args, &reply)

	if err != nil {
		log.Fatal("arith error: ", err)
	}
	fmt.Printf("Arith.Multiply: %d*%d=%d\n", args.A, args.B, reply)
}

func callDivide(a int, b int) {
	client, err := rpc.DialHTTP("tcp", addr)
	defer client.Close()
	if err != nil {
		log.Fatal("dialing:", err)
	}
	args := common.Args{a, b}
	quo := &common.Quotient{}
	err = client.Call("Arith.Divide", args, quo)

	if err != nil {
		log.Fatal("arith error: ", err)
	}

	fmt.Printf("Arith.Divide: %d/%d=%d...%d\n", args.A, args.B, quo.Quo, quo.Rem)

}

func main() {
	callMultiply(5, 4)
	callDivide(5, 4)
}
