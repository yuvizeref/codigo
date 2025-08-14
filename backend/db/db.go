package db

import (
	"context"
	"log"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func InitMongoDB() error {

	mongoURL := os.Getenv("MONGO_URL")

	clientOptions := options.Client().ApplyURI(mongoURL)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return err
	}

	DB = client.Database(getDatabaseName(mongoURL))

	log.Println("Successfully connected to MongoDB!")
	return nil
}

func getDatabaseName(mongoURL string) string {
	parts := strings.Split(mongoURL, "/")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return "codigo"
}
