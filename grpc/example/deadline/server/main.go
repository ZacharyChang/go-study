package main

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"time"

	"io"
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

type server struct {
	client pb.EchoClient
	cc     *grpc.ClientConn
}

func (s *server) UnaryEcho(ctx context.Context, in *pb.EchoRequest) (*pb.EchoResponse, error) {
	log.Printf("UnaryEcho called with message %q\n", in.GetMessage())

	if strings.HasPrefix(in.Message, "[propagate]") {
		time.Sleep(800 * time.Millisecond)
		return s.client.UnaryEcho(ctx, &pb.EchoRequest{
			Message: strings.TrimPrefix(in.Message, "[propagate]"),
		})
	}
	if in.Message == "delay" {
		time.Sleep(1500 * time.Millisecond)
	}
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
	for {
		in, err := stream.Recv()
		if err != nil {
			log.Printf("server: error receiving from stream: %v\n", err)
			if err == io.EOF {
				return nil
			}
			return err
		}
		log.Printf("echo message: %q\n", in.Message)

		if strings.HasPrefix(in.Message, "[propagate]") {
			time.Sleep(800 * time.Millisecond)
			res, err := s.client.UnaryEcho(stream.Context(), &pb.EchoRequest{
				Message: strings.TrimPrefix(in.Message, "[propagate]"),
			})
			if err != nil {
				return err
			}
			_ = stream.Send(res)
		}
		if in.Message == "delay" {
			time.Sleep(1500 * time.Millisecond)
		}

		_ = stream.Send(&pb.EchoResponse{
			Message: in.Message,
		})
	}
}

func (s *server) Close() {
	s.cc.Close()
}

func newEchoServer() *server {
	cc, err := grpc.Dial(fmt.Sprintf(":%v", *port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to connect: %v", err)
	}
	return &server{
		client: pb.NewEchoClient(cc),
		cc:     cc,
	}
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		os.Exit(1)
	}
	log.Printf("server listening at port %v\n", lis.Addr())

	echoServer := newEchoServer()
	defer echoServer.Close()

	s := grpc.NewServer()
	pb.RegisterEchoServer(s, echoServer)

	if err = s.Serve(lis); err != nil {
		log.Printf("server failed to start: %v\n", lis.Addr())
	}
}
