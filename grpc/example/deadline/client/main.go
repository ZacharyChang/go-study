package main

import (
	"context"
	"flag"
	"github.com/zacharychang/go-study/grpc/proto/echo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"os"
	"time"
)

var (
	addr = flag.String("addr", "localhost:1200", "the address to connect to")
	msg  = flag.String("message", "hello", "the message to send")
)

func unaryCall(c pb.EchoClient, requestID int, message string, want codes.Code) {
	// the timeout value will affect the got response
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.EchoRequest{
		Message: message,
	}

	_, err := c.UnaryEcho(ctx, req)
	got := status.Code(err)

	log.Printf("[%v] wanted = %v, got = %v", requestID, want, got)
}

func streamingCall(c pb.EchoClient, requestID int, message string, want codes.Code) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stream, err := c.BidirectionalStreamingEcho(ctx)
	if err != nil {
		log.Printf("Stream error: %v", err)
		return
	}

	err = stream.Send(&pb.EchoRequest{
		Message: message,
	})
	if err != nil {
		log.Printf("Send error: %v", err)
		return
	}

	_, err = stream.Recv()

	got := status.Code(err)
	log.Printf("[%v] wanted = %v, got = %v\n", requestID, want, got)
}

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to connect: %v", err)
		os.Exit(1)
	}
	defer conn.Close()

	c := pb.NewEchoClient(conn)

	unaryCall(c, 1, "world", codes.OK)
	unaryCall(c, 2, "delay", codes.DeadlineExceeded)
	unaryCall(c, 3, "[propagate]world", codes.OK)
	unaryCall(c, 4, "[propagate][propagate]world", codes.DeadlineExceeded)
	streamingCall(c, 5, "[propagate]world", codes.OK)
	streamingCall(c, 6, "[propagate][propagate]world", codes.DeadlineExceeded)
}
