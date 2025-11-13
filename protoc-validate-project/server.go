package main

import (
	"context"
	"log"
	"net"

	mainapi "github.com/grpc_tutorials/protoc-validate-project/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type server struct {
	// UnimplementedYourServiceServer
	mainapi.UnimplementedGreeterServer
}

func (s *server) Greet(ctx context.Context, req *mainapi.HelloRequest) (*mainapi.HelloResponse, error) {
	err := req.Validate()
	if err != nil {
		log.Printf("Validation error: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}
	return &mainapi.HelloResponse{Message: "Hello " + req.GetName()}, nil
}

func main() {
	port := ":50051"
	log.Printf("Starting gRPC server on port %s", port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	mainapi.RegisterGreeterServer(grpcServer, &server{})
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Printf("failed to serve: %v", err)
		return
	}
}
