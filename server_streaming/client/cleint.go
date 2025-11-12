package main

import (
	"context"
	"fmt"
	"log"
	"time"

	mainpb "github.com/grpc_tutorials/server_streaming/proto/gen"
	"google.golang.org/grpc"
)

func main() {

	port := "50052"
	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%s", port), grpc.WithInsecure())
	if err != nil {
		fmt.Printf("failed to connect: %v\n", err)
		return
	}

	defer conn.Close()

	client := mainpb.NewCalculatorClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	req := mainpb.AddRequest{
		A: 20,
		B: 20,
	}

	res, err := client.Add(ctx, &req)
	if err != nil {
		log.Printf("error while calling add rpc:%v\n", err)
		return
	}

	log.Printf("adding implementation: %d", res.Sum)
	log.Printf("getState :%s", conn.GetState())

}
