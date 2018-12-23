package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	hellopb "github.com/zacharychang/go-study/grpc/example/hello/pb"
	echopb "github.com/zacharychang/go-study/grpc/proto/echo"
	"google.golang.org/grpc"
)

var (
	addr = flag.String("addr", "localhost:1200", "the address to connect to")
	name = flag.String("name", "Jack", "the name to say hello")
	msg  = flag.String("message", "hello", "the message to send")
)

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to connect: %v", err)
		os.Exit(1)
	}
	defer conn.Close()

	helloClient := hellopb.NewGreeterClient(conn)
	echoClient := echopb.NewEchoClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// hello request
	helloRes, err := helloClient.SayHello(ctx, &hellopb.HelloRequest{
		Name: "Jack",
	})
	log.Printf("SayHello call returned %q, %v\n", helloRes.GetMessage(), err)

	// echo request
	echoRes, err := echoClient.UnaryEcho(ctx, &echopb.EchoRequest{
		Message: *msg,
	})
	log.Printf("UnaryEcho call returned %q, %v\n", echoRes.GetMessage(), err)

	if err != nil {
		log.Fatalf("request error: %v", err)
	}
}
