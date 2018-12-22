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
	msg1 = flag.String("first message", "hello", "the message to send")
	msg2 = flag.String("second message", "world", "the message to send")
)

func sendMessage(stream pb.Echo_BidirectionalStreamingEchoClient, msg string) error {
	log.Printf("sending message: %q\n", msg)
	return stream.Send(&pb.EchoRequest{
		Message: msg,
	})
}

func recvMessage(stream pb.Echo_BidirectionalStreamingEchoClient, wantCode codes.Code) {
	res, err := stream.Recv()
	if err != nil {
		log.Fatalf("stream.Recv() returned error: %v\n", err)
		return
	}
	if status.Code(err) != wantCode {
		log.Fatalf("stream.Recv() = %v, %v; want code: %v", res, err, wantCode)
	}
	log.Printf("received message: %q\n", res.GetMessage())
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	stream, err := c.BidirectionalStreamingEcho(ctx)
	if err != nil {
		log.Fatalf("error createing stream: %v", err)
		os.Exit(2)
	}

	if err := sendMessage(stream, *msg1); err != nil {
		log.Fatalf("error sending to stream: %v", err)
	}
	if err := sendMessage(stream, *msg2); err != nil {
		log.Fatalf("error sending to stream: %v", err)
	}

	recvMessage(stream, codes.OK)
	recvMessage(stream, codes.OK)

	log.Println("cancelling context...")
	cancel()

	sendMessage(stream, "closed")
	recvMessage(stream, codes.Canceled)
}
