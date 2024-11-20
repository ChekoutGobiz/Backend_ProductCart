package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Client

// Koneksi ke MongoDB
func ConnectDB() {
	mongoURI := os.Getenv("MONGOSTRING")
	if mongoURI == "" {
		log.Fatal("MongoDB connection string (MONGOSTRING) tidak ditemukan di environment.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Gagal koneksi ke MongoDB: %v", err)
	}

	// Ping database untuk memastikan koneksi
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Gagal ping ke MongoDB: %v", err)
	}

	log.Println("Koneksi MongoDB berhasil.")
	DB = client
}
