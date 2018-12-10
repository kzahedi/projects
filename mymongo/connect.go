package mongo

import (
	"context"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/readpref"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func Connect() *mongo.Client {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.Connect(ctx, "mongodb://localhost:27017")
	check(err)

	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	check(err)

	return client
}
