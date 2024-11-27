package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/kaputi/nikaro/internal/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Client instance
var Client *mongo.Client

func ConnectDB() (*mongo.Client, context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	user := os.Getenv("MONGO_USER")
	password := os.Getenv("MONGO_USER_PASSWORD")
	port := os.Getenv("MONGO_PORT")

	mongoUrl := fmt.Sprintf("mongodb://%s:%s@localhost:%s", user, password, port)

	client := utils.Must(mongo.Connect(ctx, options.Client().ApplyURI(mongoUrl)))

	utils.MustErr(client.Ping(ctx, readpref.Primary()))

	Client = client

	return client, ctx, cancel
}

// getting database collections
func GetCollection(collectionName string) *mongo.Collection {
	if Client == nil {
		panic("Database connection not established")
	}
	collection := Client.Database("nikaro").Collection(collectionName)
	return collection
}

// Close the connection
func CloseConnection(client *mongo.Client, context context.Context, cancel context.CancelFunc) {
	defer func() {
		cancel()
		if err := client.Disconnect(context); err != nil {
			panic(err)
		}
		fmt.Println("MongoDB Connection Closed")
	}()
}
