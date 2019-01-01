package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/zacharychang/go-study/grpc/proto/echo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	port = flag.Int("port", 1200, "the port to serve")
)

const (
	timestampFormat = time.StampNano
	streamingCount  = 10
)

type server struct{}

func (s *server) UnaryEcho(ctx context.Context, in *pb.EchoRequest) (*pb.EchoResponse, error) {
	log.Println("UnaryEcho called")
	defer func() {
		trailer := metadata.Pairs("timestamp", time.Now().Format(timestampFormat))
		grpc.SetTrailer(ctx, trailer)
	}()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.DataLoss, "UnaryEcho: failed to get metadata")
	}

	readTimestampFromHeader(md)

	grpc.SendHeader(ctx, initHeader())

	log.Printf("request received: %v\n", in)

	return &pb.EchoResponse{
		Message: in.Message,
	}, nil
}

func (s *server) ServerStreamingEcho(in *pb.EchoRequest, stream pb.Echo_ServerStreamingEchoServer) error {
	log.Printf("ServeStreamingEcho called")
	defer func() {
		trailer := metadata.Pairs("timestamp", time.Now().Format(timestampFormat))
		stream.SetTrailer(trailer)
	}()

	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return status.Errorf(codes.DataLoss, "ServeStreamingEcho: failed to get metadata")
	}

	readTimestampFromHeader(md)

	stream.SendHeader(initHeader())

	log.Printf("request received: %v\n", in)

	for i := 0; i < streamingCount; i++ {
		log.Printf("echo message: %v\n", in.Message)
		err := stream.Send(&pb.EchoResponse{
			Message: strconv.Itoa(i) + " " + in.Message,
		})
		time.Sleep(1 * time.Second)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *server) ClientStreamingEcho(stream pb.Echo_ClientStreamingEchoServer) error {
	log.Printf("ClientStreamingEcho called")
	defer func() {
		trailer := metadata.Pairs("timestamp", time.Now().Format(timestampFormat))
		stream.SetTrailer(trailer)
	}()

	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return status.Errorf(codes.DataLoss, "ClientStreamingEcho: failed to get metadata")
	}

	readTimestampFromHeader(md)

	stream.SendHeader(initHeader())

	var message string
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			log.Printf("echo last received message\n")
			return stream.SendAndClose(&pb.EchoResponse{
				Message: message,
			})
		}
		message = in.Message
		log.Printf("request received: %v, building echo\n", in)
		if err != nil {
			return nil
		}
	}
	return nil
}

func (s *server) BidirectionalStreamingEcho(stream pb.Echo_BidirectionalStreamingEchoServer) error {
	log.Printf("BidirectionalStreamingEcho called")
	defer func() {
		trailer := metadata.Pairs("timestamp", time.Now().Format(timestampFormat))
		stream.SetTrailer(trailer)
	}()

	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return status.Errorf(codes.DataLoss, "BidirectionalStreamingEcho: failed to get metadata")
	}

	readTimestampFromHeader(md)

	stream.SendHeader(initHeader())

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return nil
		}
		log.Printf("request received: %v, sending echo\n", in)
		err = stream.Send(&pb.EchoResponse{
			Message: in.Message,
		})
		if err != nil {
			return err
		}
		time.Sleep(2 * time.Second)
	}
}

func readTimestampFromHeader(header metadata.MD) {
	if t, ok := header["timestamp"]; ok {
		log.Printf("timestamp from header:\n")
		for i, v := range t {
			log.Printf(" %d. %s\n", i, v)
		}
	} else {
		log.Fatal("timestamp expected but not exist in header")
	}
}

func initHeader() metadata.MD {
	header := metadata.New(map[string]string{
		"location":  "MTV",
		"timestamp": time.Now().Format(timestampFormat),
	})
	return header
}

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}
	log.Printf("server listening at %v\n", lis.Addr())

	s := grpc.NewServer()
	pb.RegisterEchoServer(s, &server{})
	s.Serve(lis)
}
