package database

import (
	"fmt"
	"log"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Data-Alchemist/doculex-api/config"
)

var mongoClient *mongo.Client

func ConnectDB() *mongo.Client {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(config.ConfigDB()))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("success connect to database...")

	mongoClient = client

	return client
}

func DisconnectDB() error {
	err := mongoClient.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("closed connection from database...")

	return nil
}

func GetDB() *mongo.Client {
	return mongoClient
}

func GetCollection(client *mongo.Client, name string) *mongo.Collection {
	return client.Database(config.ConfigDBname()).Collection(name)
}