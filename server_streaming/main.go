package main

import (
	"context"
	"log"
	"net"

	mainpb "github.com/grpc_tutorials/server_streaming/proto/gen"
	"google.golang.org/grpc"
)

type server struct {
	mainpb.UnimplementedCalculatorServer
}

func (s *server) Add(ctx context.Context, in *mainpb.AddRequest) (*mainpb.AddResponse, error) {
	sum := in.A + in.B
	log.Printf("Received Add request: %d + %d = %d", in.A, in.B, sum)
	return &mainpb.AddResponse{
		Sum: sum,
	}, nil
}

func main() {

	port := ":50052"
	log.Printf("start gRPC server on port %s:", port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}

	grpcServer := grpc.NewServer()
	// register implementations
	mainpb.RegisterCalculatorServer(grpcServer, &server{})
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Printf("failed to serve:%v", err)
		return
	}
}
