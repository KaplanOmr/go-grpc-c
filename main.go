package main

import (
	"context"
	"flag"
	"fmt"
	pb "go-grpc-c/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

type proto struct {
	pb.UnimplementedApiServer
}

var (
	port = flag.Int("port", 50051, "The server port")
)

func (p *proto) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterApiServer(s, &proto{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}