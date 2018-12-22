package main

import (
	"context"
	"flag"
	"log"

	pb "github.com/zacharychang/go-study/grpc/example/user/pb"
	"google.golang.org/grpc"
)

var (
	address = flag.String("addr", "localhost:8972", "address")
	name    = flag.String("n", "world", "name")
)

func main() {
	flag.Parse()

	// 连接服务器
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("faild to connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewUserServiceClient(conn)

	r, err := c.ListUsers(context.Background(), &pb.ListUsersRequest{
		Query: "query test",
	})
	if err != nil {
		log.Fatalf("could not get response: %v", err)
	}
	log.Printf("Response: %s", r)
}
