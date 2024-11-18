package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ChekoutGobiz/BackendChekout/config"
	models "github.com/ChekoutGobiz/BackendChekout/model"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var productCollection *mongo.Collection

func init() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get MONGODB_URI from environment
	mongoURI := os.Getenv("MONGODB_URI")

	// Set MongoDB client options
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Test MongoDB connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Initialize product collection
	productCollection = client.Database("jajankuy").Collection("products")
}

// GetProductsByRegion retrieves products based on region
func GetProductsByRegion(w http.ResponseWriter, r *http.Request) {
	regionName := r.URL.Query().Get("name")
	if regionName == "" {
		http.Error(w, "Region parameter is required", http.StatusBadRequest)
		return
	}

	// Set a timeout context for the database query
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find region by name
	var region models.Region
	regionCollection := config.DB.Database("jajankuy").Collection("regions")
	err := regionCollection.FindOne(ctx, bson.M{"name": regionName}).Decode(&region)
	if err != nil {
		http.Error(w, "Region not found", http.StatusNotFound)
		return
	}

	// Find products by region ID
	var products []models.Product
	cursor, err := productCollection.Find(ctx, bson.M{"region_id": region.ID})
	if err != nil {
		http.Error(w, "Failed to retrieve products", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	// Decode the products from the cursor
	for cursor.Next(ctx) {
		var product models.Product
		if err := cursor.Decode(&product); err != nil {
			http.Error(w, "Failed to decode product", http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}

	// Return the products as a JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products); err != nil {
		http.Error(w, "Failed to encode products to JSON", http.StatusInternalServerError)
	}
}

// CreateProduct adds a new product to the database
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	// Generate a new ObjectID for the product
	product.ID = primitive.NewObjectID()

	// Set a timeout context for the database operation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Insert the new product into the database
	_, err := productCollection.InsertOne(ctx, product)
	if err != nil {
		http.Error(w, "Failed to insert product", http.StatusInternalServerError)
		return
	}

	// Return the newly created product as a JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// GetProducts retrieves all products from the database
func GetProducts(w http.ResponseWriter, r *http.Request) {
	var products []models.Product

	// Set a timeout context for the database query
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Retrieve all products from the database
	cursor, err := productCollection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to get products", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	// Decode the products from the cursor
	for cursor.Next(ctx) {
		var product models.Product
		if err := cursor.Decode(&product); err != nil {
			http.Error(w, "Failed to decode product", http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}

	// Return the list of products as a JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products); err != nil {
		http.Error(w, "Failed to encode products to JSON", http.StatusInternalServerError)
	}
}
