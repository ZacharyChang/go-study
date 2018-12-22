package main

import (
	"context"
	"github.com/zacharychang/go-study/grpc/example/hello/pb"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"log"
	"os"
	"time"
)

var (
	addr = "localhost:1200"
)

func main() {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Failed to connect:", err.Error())
	}
	defer func() {
		if e := conn.Close(); e != nil {
			log.Println("Failed to close connection:", e.Error())
		}
	}()
	c := pb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := c.SayHello(ctx, &pb.HelloRequest{
		Name: "Jack",
	})
	if err != nil {
		s := status.Convert(err)
		for _, v := range s.Details() {
			switch info := v.(type) {
			case *errdetails.QuotaFailure:
				log.Println("Quota failure:", info)
			default:
				log.Println("Unexpected type:", info)
			}
		}
		os.Exit(1)
	}
	log.Println("Received:", resp.Message)
}
