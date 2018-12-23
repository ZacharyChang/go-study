package main

import (
	"context"
	"flag"
	"github.com/zacharychang/go-study/grpc/proto/echo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding/gzip"
	"log"
	"os"
	"time"
)

var (
	addr = flag.String("addr", "localhost:1200", "the address to connect to")
	msg  = flag.String("message", "hello", "the message to send")
)

func main() {
	flag.Parse()

	creds, err := credentials.NewClientTLSFromFile("../cert/zacharychang.com.pem", "zacharychang.com")
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("fail to connect: %v", err)
		os.Exit(1)
	}
	defer conn.Close()

	c := pb.NewEchoClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := c.UnaryEcho(ctx, &pb.EchoRequest{
		Message: *msg,
	}, grpc.UseCompressor(gzip.Name))
	log.Printf("UnaryEcho call returned %q, %v\n", res.GetMessage(), err)
	if err != nil {
		log.Fatalf("Message=%q, err=%v; Want message: %q", res.GetMessage(), err, msg)
	}
}
