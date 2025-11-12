package main

import (
	"context"
	"fmt"
	"log"
	mainpipb "simplegrpc/grpc_client/proto/gen"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/metadata"
)

func main() {

	cert := "cert.pem"

	creds, err := credentials.NewClientTLSFromFile(cert, "")
	if err != nil {
		fmt.Printf("failed to load TLS credentials: %v\n", err)
		return
	}

	port := "50051"
	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%s", port), grpc.WithTransportCredentials(creds), grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)))
	if err != nil {
		fmt.Printf("failed to connect: %v\n", err)
		return
	}

	defer conn.Close()

	client := mainpipb.NewCalculateClient(conn)

	client2 := mainpipb.NewGreeterClient(conn)
	fwclient := mainpipb.NewBidFarewellClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// add client
	// add encryption on individual RPC
	/* grpc.UseCompressor(gzip.Name) */
	// create metadata
	md := metadata.Pairs("authorization", "Bearer=jwt-token")
	ctx = metadata.NewOutgoingContext(ctx, md)

	// process income metadata from server
	var resHeader metadata.MD
	var resTrailer metadata.MD

	req := mainpipb.AddRequest{
		A: 10,
		B: 20,
	}
	res, err := client.Add(ctx, &req, grpc.Header(&resHeader))
	if err != nil {
		fmt.Printf("error while calling Add RPC: %v\n", err)
		return
	}
	log.Println("Response Header from server: ", resHeader)
	log.Printf("resHeader : %s", resHeader["timestamp"][0])
	log.Println("Response Trailer from server: ", resTrailer)
	log.Println("Response Trailer from server: ", resTrailer["processedTimestamp"])

	// hello client
	helloReq := mainpipb.HelloRequest{
		Name: "Kenneth",
	}
	helloRes, err := client2.Greet(ctx, &helloReq)
	if err != nil {
		fmt.Printf("error while calling Greet RPC: %v\n", err)
		return
	}

	// farewell client
	farewellRequest := mainpipb.GoodbyeRequest{
		Name: "Kenneth",
	}
	farewell, err := fwclient.BidGoodbye(ctx, &farewellRequest)
	if err != nil {
		fmt.Printf("error while calling BidGoodbye RPC: %v\n", err)
		return
	}
	// farewell
	log.Println("Farewell: ", farewell.Message)

	// greet
	log.Println("Greeting: ", helloRes.Message)
	log.Println("Sum: ", res.Sum)
	state := conn.GetState()
	log.Println("Connection state: ", state)

}
