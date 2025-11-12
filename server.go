package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "simplegrpc/proto/gen"

	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/metadata"
)

type server struct {
	// UnimplementedYourServiceServer
	pb.UnimplementedCalculateServer
	pb.UnimplementedGreeterServer
	pb.UnimplementedBidFarewellServer
}

// Add implementation
func (s *server) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	metdata, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("No metadata found")
	}
	val, ok := metdata["authorization"]
	if !ok {
		log.Println("No authorization token found")
	}
	log.Printf("Received metadata - authorization: %v", val[0])

	// set response headers
	responseHeader := metadata.Pairs("timestamp", fmt.Sprintf("%d", 123456789))
	err := grpc.SendHeader(ctx, responseHeader)
	if err != nil {
		log.Printf("Failed to send header: %v", err)
		return nil, err
	}

	sum := req.A + req.B
	log.Printf("Received Add request: %d + %d = %d", req.A, req.B, sum)

	trailer := metadata.Pairs("processedTimestamp", fmt.Sprintf("%d", 987654321))
	err = grpc.SetTrailer(ctx, trailer)
	if err != nil {
		log.Printf("Failed to set trailer: %v", err)
		return nil, err
	}
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

	// cert := "cert.pem"
	// key := "key.pem"
	// port := ":50051"
	// log.Printf("Starting gRPC server on port %s", port)
	// lis, err := net.Listen("tcp", port)
	// if err != nil {
	// 	log.Printf("failed to listen: %v", err)
	// 	return
	// }

	// // load server certificate and key
	// creds, err := credentials.NewServerTLSFromFile(cert, key)
	// if err != nil {
	// 	log.Printf("failed to load TLS credentials: %v", err)
	// 	return
	// }

	// grpcServer := grpc.NewServer(grpc.Creds(creds))

	// pb.RegisterCalculateServer(grpcServer, &server{})
	// pb.RegisterGreeterServer(grpcServer, &server{})
	// pb.RegisterBidFarewellServer(grpcServer, &server{})
	// err = grpcServer.Serve(lis)
	// if err != nil {
	// 	log.Printf("failed to serve: %v", err)
	// 	return
	// }

	port := ":50051"
	log.Printf("Starting gRPC server on port %s", port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}

	grpcServer := grpc.NewServer()

	pb.RegisterCalculateServer(grpcServer, &server{})
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Printf("failed to serve: %v", err)
		return
	}
}
