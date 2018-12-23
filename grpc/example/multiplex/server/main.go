package main

import (
	"context"
	"flag"
	"fmt"

	"log"
	"net"
	"os"

	hellopb "github.com/zacharychang/go-study/grpc/example/hello/pb"
	echopb "github.com/zacharychang/go-study/grpc/proto/echo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	port = flag.Int("port", 1200, "the port to listen")
)

type helloServer struct{}

func (s *helloServer) SayHello(ctx context.Context, in *hellopb.HelloRequest) (*hellopb.HelloResponse, error) {
	return &hellopb.HelloResponse{
		Message: "Hello " + in.Name,
	}, nil
}

type echoServer struct{}

func (s *echoServer) UnaryEcho(ctx context.Context, in *echopb.EchoRequest) (*echopb.EchoResponse, error) {
	log.Printf("UnaryEcho called with message %q\n", in.GetMessage())
	return &echopb.EchoResponse{
		Message: in.Message,
	}, nil
}

func (s *echoServer) ServerStreamingEcho(in *echopb.EchoRequest, stream echopb.Echo_ServerStreamingEchoServer) error {
	return status.Error(codes.Unimplemented, "todo")
}

func (s *echoServer) ClientStreamingEcho(stream echopb.Echo_ClientStreamingEchoServer) error {
	return status.Error(codes.Unimplemented, "todo")
}

func (s *echoServer) BidirectionalStreamingEcho(stream echopb.Echo_BidirectionalStreamingEchoServer) error {
	return status.Error(codes.Unimplemented, "todo")
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		os.Exit(1)
	}
	log.Printf("server listening at port %v\n", lis.Addr())

	s := grpc.NewServer()
	hellopb.RegisterGreeterServer(s, &helloServer{})
	echopb.RegisterEchoServer(s, &echoServer{})

	if err = s.Serve(lis); err != nil {
		log.Printf("server failed to start: %v\n", lis.Addr())
	}
}
