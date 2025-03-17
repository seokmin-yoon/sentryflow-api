package config

import (
	"context"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient *mongo.Client
	mongoOnce   sync.Once
)

// GetMongoClient initializes and returns a MongoDB client
func GetMongoClient() (*mongo.Client, error) {
	var err error
	mongoOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// os.Getenv("MONGODB_ADDR") => mongodb://mongodb:27017
		// mongodb://mongodb.sentryflow.svc.cluster.local:27017
		clientOptions := options.Client().ApplyURI("mongodb://mongodb.sentryflow.svc.cluster.local:27017")
		mongoClient, err = mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Fatalf("[MongoDB] Connection error: %v", err)
		}
	})

	return mongoClient, err
}
