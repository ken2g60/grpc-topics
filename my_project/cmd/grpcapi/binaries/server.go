package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/grpc_tutorials/my_project/internal/api/handlers"
	"github.com/grpc_tutorials/my_project/internal/repositories/mongodb"
	mainapi "github.com/grpc_tutorials/my_project/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// err := godotenv.Load("../.env")
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// 	return
	// }

	// database configuration
	mongodb.CreateMongoClient()

	s := grpc.NewServer()
	mainapi.RegisterExecsServiceServer(s, &handlers.Server{})
	mainapi.RegisterStudentsServiceServer(s, &handlers.Server{})
	mainapi.RegisterTeachersServiceServer(s, &handlers.Server{})

	reflection.Register(s)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = ":50021"
	}
	fmt.Println("gRPC Server is running on port", port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Printf("Error listen to port %v", err)
		return
	}

	// listen to configuration
	err = s.Serve(lis)
	if err != nil {
		log.Printf("Error serve %s", err)
		return
	}
}
