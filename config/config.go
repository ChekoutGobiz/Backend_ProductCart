package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Global variable to hold the database client instance
var DB *mongo.Client

// ConnectDB initializes the MongoDB connection and assigns it to DB
func ConnectDB() *mongo.Client {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get MONGODB_URI from environment variables
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI is not set in .env")
	}

	// Set MongoDB client options
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}

	// Context with timeout for the connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Release resources if connection setup takes longer

	// Connect to MongoDB
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ping MongoDB to ensure the connection is established
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	log.Println("MongoDB connection established successfully!")
	DB = client
	return client
}

// GetCollection retrieves a specific collection from the database
func GetCollection(collectionName string) *mongo.Collection {
	if DB == nil {
		log.Fatal("Database connection has not been initialized")
	}
	return DB.Database("jajankuy").Collection(collectionName)
}
