package main

import (
	"context"
	"fmt"
	"github.com/zacharychang/go-study/grpc/hello-example/pb"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"os"
	"sync"
	"v2ray.com/core/common/net"
)

var (
	port = ":1200"
)

type server struct {
	mutex sync.Mutex
	count map[string]int
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.count[in.Name]++
	if s.count[in.Name] > 1 {
		st := status.New(codes.ResourceExhausted, "Request limit exceeded.")
		ds, err := st.WithDetails(&errdetails.QuotaFailure{
			Violations: []*errdetails.QuotaFailure_Violation{
				{
					Subject:     fmt.Sprintf("name: %s", in.Name),
					Description: "Limit one greeting per person",
				},
			},
		},
		)
		if err != nil {
			return nil, st.Err()
		}
		return nil, ds.Err()
	}
	return &pb.HelloResponse{
		Message: "Hello " + in.Name,
	}, nil
}

func main() {
	log.Println("Port listening on: ", port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln("Failed to listen: ", err.Error())
		os.Exit(1)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{
		count: make(map[string]int),
	})
	if err := s.Serve(lis); err != nil {
		log.Fatalln("Failed to serve: ", err.Error())
		os.Exit(2)
	}
}
