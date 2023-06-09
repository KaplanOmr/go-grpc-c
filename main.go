package main

import (
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
