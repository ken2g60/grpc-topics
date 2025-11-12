package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"

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

func (s *server) GenerateFibonacci(req *mainpb.FibonacciRequest, stream mainpb.Calculator_GenerateFibonacciServer) error {
	n := req.N
	a, b := 0, 1

	for i := 0; i < int(n); i++ {
		err := stream.Send(&mainpb.FibonacciResponse{
			Number: int32(a),
		})
		if err != nil {
			return err
		}
		log.Println("Sent Fibonacci number:", a)
		a, b = b, a+b
		time.Sleep(time.Second)
	}
	return nil
}

func (s *server) SendNumbers(stream mainpb.Calculator_SendNumbersServer) error {
	var sum int32 = 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&mainpb.NumberResponse{Sum: sum})
		}
		if err != nil {
			return err
		}
		log.Println(req.GetNumber())
		sum += req.GetNumber()
	}
}

func (s *server) Chat(stream mainpb.Calculator_ChatServer) error {
	reader := bufio.NewReader(os.Stdin)
	for {
		// receive data/message/values from the stream
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("error while receiving from client:", err)
		}
		log.Println("Received from client:", req.GetMessage())

		// read input from the terminal
		fmt.Print("Enter response:")
		input, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		input = strings.TrimSpace(input)

		// send data/message/values through the stream
		response := &mainpb.ChatMessage{
			Message: input,
		}

		err = stream.Send(response)
		if err != nil {
			return err
		}

		// response message logged on server side
	}
	fmt.Println("chat ended")
	return nil
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
