package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const connectionString = "mongodb://localhost:27017"
const db = "go-grpc-c-db"

func close(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {
	defer cancel()

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func connect() (*mongo.Client, context.Context, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))

	return client, ctx, cancel, err
}

func Ping() {
	client, ctx, cancel, err := connect()
	if err != nil {
		fmt.Println("connected failed")
	}
	defer close(client, ctx, cancel)

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		fmt.Println("connected failed")
	}

	fmt.Println("connected successfully")
}

func Insert(collection string, data interface{}) (*mongo.InsertOneResult, error) {
	client, ctx, cancel, err := connect()
	if err != nil {
		return nil, errors.New("CANNOT_CONNECTED_TO_DB")
	}
	defer close(client, ctx, cancel)

	coll := client.Database(db).Collection(collection)
	result, err := coll.InsertOne(ctx, data)
	if err != nil {
		return nil, errors.New("CANNOT_CREATED_USER")
	}

	return result, nil
}

func Find(collection string, data interface{}) (*mongo.SingleResult, error) {
	client, ctx, cancel, err := connect()
	if err != nil {
		return nil, errors.New("CANNOT_CONNECTED_TO_DB")
	}
	defer close(client, ctx, cancel)

	coll := client.Database(db).Collection(collection)
	result := coll.FindOne(ctx, data)
	if result.Err() != nil {
		return nil, errors.New("USER_CANNOT_FOUND")
	}

	return result, nil
}
