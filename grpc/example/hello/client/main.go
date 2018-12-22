package main

import (
	"context"
	"github.com/zacharychang/go-study/grpc/example/hello/pb"
	"google.golang.org/grpc"
	"log"
	"os"
	"time"
)

var (
	addr = "localhost:1200"
	name = "world"
)

func main() {
	if len(os.Args) == 2 {
		name = os.Args[1]
	}

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	checkErr(err)

	defer conn.Close()
	client := pb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// Cancel function
	defer cancel()

	resp, err := client.SayHello(ctx, &pb.HelloRequest{
		Name: name,
	})
	checkErr(err)

	log.Println("Received: ", resp.Message)
}

func checkErr(err error) {
	if err != nil {
		log.Fatalf("Error: %v", err)
		os.Exit(1)
	}
}
