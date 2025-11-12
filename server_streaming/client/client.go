package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	mainpb "github.com/grpc_tutorials/server_streaming/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	port := "50052"
	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%s", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("failed to connect: %v\n", err)
		return
	}

	defer conn.Close()

	client := mainpb.NewCalculatorClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

	// Server Streaming RPC
	streamReq := mainpb.FibonacciRequest{
		N: 10,
	}

	stream, err := client.GenerateFibonacci(ctx, &streamReq)
	if err != nil {
		log.Printf("error while calling GenerateFibonacci rpc:%v\n", err)
		return
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			log.Printf("no more data in stream:%v\n", err)
			return
		}
		if err != nil {
			log.Printf("error while receiving stream:%v\n", err)
			break
		}
		log.Printf("Fibonacci number: %d", msg.GetNumber())
	}

	log.Printf("adding implementation: %d", res.Sum)
	log.Printf("getState :%s", conn.GetState())

	// client side streaming rpc
	streamCtx, streamCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer streamCancel()
	stream1, err := client.SendNumbers(streamCtx)
	if err != nil {
		log.Printf("error while calling SendNumbers rpc:%v\n", err)
		return
	}

	for num := range 9 {
		err := stream1.Send(&mainpb.NumberRequest{Number: int32(num)})
		if err != nil {
			log.Printf("error while sending data to stream:%v\n", err)
			return
		}
		time.Sleep(time.Second)
	}

	clientres, err := stream1.CloseAndRecv()
	if err != nil {
		log.Printf("error while receiving response from SendNumbers:%v\n", err)
		return
	}

	log.Println("sum of numbers sent:", clientres.Sum)

	// bidirectional streaming rpc
	bidirectionCtx, bidirectionalCancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer bidirectionalCancel()
	chatStream, err := client.Chat(bidirectionCtx)
	if err != nil {
		log.Printf("error while creating chat stream:%v\n", err)
		return
	}

	waitc := make(chan struct{})
	// send message in a goroutine
	go func() {
		messages := []string{"Hello", "How are you?", "I am fine.", "Goodbye!"}
		for _, message := range messages {
			err := chatStream.Send(&mainpb.ChatMessage{Message: message})
			if err != nil {
				log.Printf("error while sending message to chat stream:%v\n", err)
				return
			}
			time.Sleep(time.Second)
		}
		chatStream.CloseSend()
	}()

	// receive messages
	go func() {
		for {
			res, err := chatStream.Recv()
			if err == io.EOF {
				log.Println("chat ended by server")
				return
			}
			if err != nil {
				log.Printf("error while receiving message from chat stream:%v\n", err)
				return
			}
			log.Printf("Received from server: %s", res.GetMessage())
		}
		close(waitc)
	}()
	<-waitc

}
