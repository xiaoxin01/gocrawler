package model

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	// Client client to mongo db
	Client   *mongo.Client
	Database *mongo.Database
)

// InitDb initial db
func InitDb() {
	db := Db

	// Set client options
	clientOptions := options.Client().ApplyURI(db.Connection)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}
	database := client.Database(db.Database)

	Client = client
	Database = database
}
