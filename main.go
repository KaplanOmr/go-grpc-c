package main

import (
	"context"
	"flag"
	"fmt"
	"go-grpc-c/db"
	pb "go-grpc-c/proto"
	"log"
	"net"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type proto struct {
	pb.UnimplementedApiServer
}

var (
	port = flag.Int("port", 50051, "The server port")
)

func (p *proto) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	log.Printf("Token: %s", strings.Split(md["authorization"][0], " ")[1])

	dbClient, dbCtx, dbCancelFunc, dbErr := db.Connect()
	if dbErr != nil {
		panic(dbErr)
	}
	defer db.Close(dbClient, dbCtx, dbCancelFunc)

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
