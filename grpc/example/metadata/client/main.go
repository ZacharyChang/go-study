package main

import (
	"context"
	"flag"
	"io"
	"log"
	"time"

	"github.com/zacharychang/go-study/grpc/proto/echo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	addr = flag.String("addr", "localhost:1200", "the address to connect")
	msg  = flag.String("msg", "hello world", "message to send")
)

const (
	timestampFormat = time.StampNano
	streamingCount  = 10
)

func unaryCallWithMetadata(c pb.EchoClient, message string) {
	log.Printf("unaryCallWithMetadata called")
	md := metadata.Pairs("timestamp", time.Now().Format(timestampFormat))
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	var header, trailer metadata.MD
	r, err := c.UnaryEcho(ctx, &pb.EchoRequest{
		Message: message,
	}, grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		log.Fatalf("fail to call UnaryEcho: %v", err)
	}

	if t, ok := header["timestamp"]; ok {
		log.Printf("timestamp from header:\n")
		for i, v := range t {
			log.Printf(" %d. %s\n", i, v)
		}
	} else {
		log.Fatalf("timestamp expected but not exist in header")
	}

	if l, ok := header["location"]; ok {
		log.Printf("location from header:\n")
		for i, v := range l {
			log.Printf(" %d. %s\n", i, v)
		}
	} else {
		log.Fatalf("location expected but not exist in header")
	}

	log.Println("response:")
	log.Printf(" - %s\n", r.Message)

	if t, ok := trailer["timestamp"]; ok {
		log.Printf("timestamp from tailer:\n")
		for i, v := range t {
			log.Printf(" %d. %s\n", i, v)
		}
	} else {
		log.Fatalf("timestamp expected but not exist in trailer")
	}
}

func serverStreamingWithMetadata(c pb.EchoClient, message string) {
	log.Printf("serverStreamingWithMetadata called")
	md := metadata.Pairs("timestamp", time.Now().Format(timestampFormat))
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	stream, err := c.ServerStreamingEcho(ctx, &pb.EchoRequest{
		Message: message,
	})
	if err != nil {
		log.Fatalf("fail to call ServerStreamingEcho: %v", err)
	}

	header, err := stream.Header()
	if err != nil {
		log.Fatalf("fail to get header from stream: %v", err)
	}

	// read timestamp from header
	if t, ok := header["timestamp"]; ok {
		log.Printf("timestamp from header:\n")
		for i, v := range t {
			log.Printf(" %d. %s\n", i, v)
		}
	} else {
		log.Fatalf("timestamp expected but not exist int header")
	}

	// read location from header
	if l, ok := header["location"]; ok {
		log.Printf("location from header:\n")
		for i, v := range l {
			log.Printf(" %d. %s\n", i, v)
		}
	} else {
		log.Fatalf("location expected but not exist in header")
	}

	// read the response
	var rpcStatus error
	for {
		r, err := stream.Recv()
		if err != nil {
			rpcStatus = err
			break
		}
		log.Printf(" - %s\n", r.Message)
	}
	if rpcStatus != io.EOF {
		log.Fatalf("fail to finish server streaming: %v", rpcStatus)
	}

	//  read after RPC finished
	trailer := stream.Trailer()
	// read timestamp from trailer
	if t, ok := trailer["timestamp"]; ok {
		log.Printf("timestamp from tailer:\n")
		for i, v := range t {
			log.Printf(" %d. %s\n", i, v)
		}
	} else {
		log.Fatalf("timestamp expected but not exist in trailer")
	}
}

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("connect fail: %v", err)
	}
	defer conn.Close()

	c := pb.NewEchoClient(conn)

	log.Println("------------")
	unaryCallWithMetadata(c, *msg)
	time.Sleep(1 * time.Second)

	log.Println("------------")
	serverStreamingWithMetadata(c, *msg)

}