package mongodb

import (
	"context"
	"log"

	"github.com/grpc_tutorials/my_project/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateMongoClient() (*mongo.Client, error) {
	ctx := context.Background()
	// client, err := mongo.Connect(ctx, options.Client().ApplyURI("username:password@mongodb://localhost:27017"))
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Printf("Failed to connect %v", err)
		return nil, utils.ErrorHandler(err, "unable to connect to mongodb")
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("failed ping %v", err)
		return nil, nil
	}

	log.Println("connected to Mongodb")
	return client, nil
}
