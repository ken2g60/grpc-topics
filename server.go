package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "simplegrpc/proto/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type server struct {
	// UnimplementedYourServiceServer
	pb.UnimplementedCalculateServer
	pb.UnimplementedGreeterServer
	pb.UnimplementedBidFarewellServer
}

// Add implementation
func (s *server) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	sum := req.A + req.B
	log.Printf("Received Add request: %d + %d = %d", req.A, req.B, sum)
	return &pb.AddResponse{
		Sum: sum,
	}, nil

}

// Greet implementation
func (s *server) Greet(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	name := req.Name
	log.Printf("Received Greet request: %s", name)
	return &pb.HelloResponse{
		Message: fmt.Sprintf("Hello from %s", name),
	}, nil

}

// BidGoodbye implementation
func (s *server) BidGoodbye(ctx context.Context, req *pb.GoodbyeRequest) (*pb.GoodbyeResponse, error) {
	name := req.Name
	log.Printf("Recieved GoodbyeName request: %s", name)
	return &pb.GoodbyeResponse{
		Message: name,
	}, nil
}

func main() {

	cert := "cert.pem"
	key := "key.pem"
	port := ":50051"
	log.Printf("Starting gRPC server on port %s", port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}

	// load server certificate and key
	creds, err := credentials.NewServerTLSFromFile(cert, key)
	if err != nil {
		log.Printf("failed to load TLS credentials: %v", err)
		return
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds))

	pb.RegisterCalculateServer(grpcServer, &server{})
	pb.RegisterGreeterServer(grpcServer, &server{})
	pb.RegisterBidFarewellServer(grpcServer, &server{})
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Printf("failed to serve: %v", err)
		return
	}

}
