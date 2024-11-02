package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kuyjajan/kuyjajan-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var regionCollection *mongo.Collection

func init() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Ambil MONGODB_URI dari environment
	mongoURI := os.Getenv("MONGODB_URI")

	// Opsi koneksi MongoDB
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Cek koneksi MongoDB
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	log.Println("MongoDB connection established successfully!")

	// Initialize region collection
	regionCollection = client.Database("jajankuy").Collection("regions")
}

// CreateRegion handles the creation of a new region
func CreateRegion(w http.ResponseWriter, r *http.Request) {
	var region models.Region
	if err := json.NewDecoder(r.Body).Decode(&region); err != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	// Generate new ObjectID for region
	region.ID = primitive.NewObjectID()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := regionCollection.InsertOne(ctx, region)
	if err != nil {
		http.Error(w, "Failed to insert region", http.StatusInternalServerError)
		return
	}

	// Set response header and encode region to JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(region)
}

// GetRegions retrieves all regions from the database
func GetRegions(w http.ResponseWriter, r *http.Request) {
	var regions []models.Region
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := regionCollection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to get regions", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var region models.Region
		if err := cursor.Decode(&region); err != nil {
			http.Error(w, "Failed to decode region", http.StatusInternalServerError)
			return
		}
		regions = append(regions, region)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(regions)
}
