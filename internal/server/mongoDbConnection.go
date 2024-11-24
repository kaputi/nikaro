package server

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

// Open new connection
func SetupMongoDB() (*mongo.Collection, *mongo.Client, context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(fmt.Sprintf("Mongo DB Connect issue %s", err))
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(fmt.Sprintf("Mongo DB ping issue %s", err))
	}
	collection := client.Database("mongo-golang-test").Collection("Users")
	return collection, client, ctx, cancel
}

func ConectMongoDb() (*mongo.Client, context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	user, ok := os.LookupEnv("MONGO_USER")
	if !ok {
		panic("MONGO_USER not found")
	}
	password, ok := os.LookupEnv("MONGO_USER_PASSWORD")
	if !ok {
		panic("MONGO_USER_PASSWORD not found")
	}
	port, ok := os.LookupEnv("MONGO_PORT")
	if !ok {
		panic("MONGO_PORT not found")
	}

	mongoUrl := fmt.Sprintf("mongodb://%s:%s@localhost:%s", user, password, port)

	client := utils.Must(mongo.Connect(ctx, options.Client().ApplyURI(mongoUrl)))

	utils.MustErr(client.Ping(ctx, readpref.Primary()))

	return client, ctx, cancel
}

// Close the connection
func CloseConnection(client *mongo.Client, context context.Context, cancel context.CancelFunc) {
	defer func() {
		cancel()
		if err := client.Disconnect(context); err != nil {
			panic(err)
		}
		fmt.Println("Close connection is called")
	}()
}
