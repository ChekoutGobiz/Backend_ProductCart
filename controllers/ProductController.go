package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kuyjajan/kuyjajan-backend/config"
	"github.com/kuyjajan/kuyjajan-backend/models"
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

	// Initialize product collection
	productCollection = client.Database("jajankuy").Collection("products")
}

// GetProductsByRegion menampilkan produk berdasarkan wilayah
func GetProductsByRegion(w http.ResponseWriter, r *http.Request) {
	regionName := r.URL.Query().Get("name")
	fmt.Println("Region parameter:", regionName)
	if regionName == "" {
		http.Error(w, "Region parameter is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Cari Region berdasarkan nama
	var region models.Region
	regionCollection := config.DB.Database("jajankuy").Collection("regions")
	err := regionCollection.FindOne(ctx, bson.M{"name": regionName}).Decode(&region)
	if err != nil {
		http.Error(w, "Region not found", http.StatusNotFound)
		return
	}

	// Cari Produk berdasarkan RegionID
	var products []models.Product
	cursor, err := productCollection.Find(ctx, bson.M{"region_id": region.ID})
	if err != nil {
		http.Error(w, "Failed to retrieve products", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var product models.Product
		if err := cursor.Decode(&product); err != nil {
			http.Error(w, "Failed to decode product", http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products); err != nil {
		http.Error(w, "Failed to encode products to JSON", http.StatusInternalServerError)
	}
}

// CreateProduct menambah produk baru ke dalam database
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	// Generate ObjectID baru untuk produk
	product.ID = primitive.NewObjectID()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Masukkan produk ke dalam database
	_, err := productCollection.InsertOne(ctx, product)
	if err != nil {
		http.Error(w, "Failed to insert product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// GetProducts menampilkan semua produk dari database
func GetProducts(w http.ResponseWriter, r *http.Request) {
	var products []models.Product
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := productCollection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to get products", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var product models.Product
		if err := cursor.Decode(&product); err != nil {
			http.Error(w, "Failed to decode product", http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
