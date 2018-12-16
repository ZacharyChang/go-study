//go:generate protoc -I ../pb --go_out=plugins=grpc:../pb ../pb/hello.proto
package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/zacharychang/go-study/grpc/hello-example/pb"
	"google.golang.org/grpc"
)

const (
	port = ":1200"
)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Println("Received: ", in.Name)
	return &pb.HelloResponse{
		Message: "Hello, " + in.Name,
	}, nil
}

func main() {
	svc, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
		os.Exit(1)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	// https://github.com/grpc/grpc-go/blob/master/Documentation/server-reflection-tutorial.md
	//reflection.Register(s)
	if err := s.Serve(svc); err != nil {
		log.Fatalf("Failed to serve: %v", err)
		os.Exit(2)
	}
}
