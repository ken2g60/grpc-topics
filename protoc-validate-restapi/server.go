package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	mainapi "github.com/grpc_tutorials/protoc-validate-project/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
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

func runGRPCServer() {
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

// RESTAPI
func runGatewayServer() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := mainapi.RegisterGreeterHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
	if err != nil {
		log.Fatalf("Failed to register gateway: %v", err)
		return
	}

	log.Println("Starting REST API server on port :8080")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("failed to start REST API server: %v", err)
	}
}

func main() {
	go runGRPCServer()
	runGatewayServer()
}
