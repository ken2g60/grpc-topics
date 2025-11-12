package main

import (
	"context"
	"fmt"
	"log"
	mainpipb "simplegrpc/grpc_client/proto/gen"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {

	cert := "cert.pem"

	creds, err := credentials.NewClientTLSFromFile(cert, "")
	if err != nil {
		fmt.Printf("failed to load TLS credentials: %v\n", err)
		return
	}

	port := "50051"
	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%s", port), grpc.WithTransportCredentials(creds))
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
	req := mainpipb.AddRequest{
		A: 10,
		B: 20,
	}
	res, err := client.Add(ctx, &req)
	if err != nil {
		fmt.Printf("error while calling Add RPC: %v\n", err)
		return
	}

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
