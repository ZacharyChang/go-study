package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc/credentials"

	"log"
	"net"
	"os"

	"github.com/zacharychang/go-study/grpc/proto/echo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	_ "google.golang.org/grpc/encoding/gzip" // Install the gzip compressor
	"google.golang.org/grpc/status"
)

var (
	port = flag.Int("port", 1200, "the port to listen")
)

type server struct{}

func (s *server) UnaryEcho(ctx context.Context, in *pb.EchoRequest) (*pb.EchoResponse, error) {
	log.Printf("UnaryEcho called with message %q\n", in.GetMessage())
	return &pb.EchoResponse{
		Message: in.Message,
	}, nil
}

func (s *server) ServerStreamingEcho(in *pb.EchoRequest, stream pb.Echo_ServerStreamingEchoServer) error {
	return status.Error(codes.Unimplemented, "todo")
}

func (s *server) ClientStreamingEcho(stream pb.Echo_ClientStreamingEchoServer) error {
	return status.Error(codes.Unimplemented, "todo")
}

func (s *server) BidirectionalStreamingEcho(stream pb.Echo_BidirectionalStreamingEchoServer) error {
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

	creds, err := credentials.NewServerTLSFromFile("../cert/zacharychang.com.pem", "../cert/private.pem")
	if err != nil {
		log.Fatalf("failed to create credentials: %v", err)
	}

	s := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterEchoServer(s, &server{})

	if err = s.Serve(lis); err != nil {
		log.Printf("server failed to start: %v\n", lis.Addr())
	}
}
