package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/kuyjajan/kuyjajan-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var cartCollection *mongo.Collection

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
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Cek koneksi MongoDB
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	log.Println("MongoDB connection established successfully!")

	// Initialize cart collection
	cartCollection = client.Database("jajankuy").Collection("carts")
}

// AddToCart adds an item to the user's cart
func AddToCart(w http.ResponseWriter, r *http.Request) {
	// Struktur untuk membaca request body
	var requestBody struct {
		UserID    string            `json:"user_id"`
		ProductID string            `json:"product_id"`
		Quantity  int               `json:"quantity"`
	}
	
	// Dekode JSON body
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validasi userID
	if requestBody.UserID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Konversi userID dari string ke ObjectID
	userID, err := primitive.ObjectIDFromHex(requestBody.UserID)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	// Validasi productID dan konversi ke ObjectID
	if requestBody.ProductID == "" {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		return
	}
	productID, err := primitive.ObjectIDFromHex(requestBody.ProductID)
	if err != nil {
		http.Error(w, "Invalid product ID format", http.StatusBadRequest)
		return
	}

	// Buat item keranjang berdasarkan input
	cartItem := models.CartItem{
		ProductID: productID,
		Quantity:  requestBody.Quantity,
	}

	// Cari atau buat keranjang untuk pengguna
	var cart models.Cart
	err = cartCollection.FindOne(context.TODO(), bson.M{"user_id": userID}).Decode(&cart)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Buat keranjang baru jika tidak ditemukan
			cart = models.Cart{
				UserID: userID,
				Items:  []models.CartItem{cartItem},
			}
			_, err = cartCollection.InsertOne(context.TODO(), cart)
		} else {
			http.Error(w, "Failed to retrieve or create cart", http.StatusInternalServerError)
			return
		}
	} else {
		// Jika keranjang ada, tambahkan item ke dalamnya
		cart.Items = append(cart.Items, cartItem)
		_, err = cartCollection.UpdateOne(context.TODO(), bson.M{"user_id": userID}, bson.M{"$set": bson.M{"items": cart.Items}})
	}

	if err != nil {
		http.Error(w, "Failed to add item to cart", http.StatusInternalServerError)
		return
	}

	// Kirimkan respons
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}

func GetCart(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Konversi userID dari string ke ObjectID
	objectID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	var cart models.Cart
	err = cartCollection.FindOne(context.TODO(), bson.M{"user_id": objectID}).Decode(&cart)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Cart not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to retrieve cart", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}

// UpdateCartItem updates the quantity of a specific item in the cart
func UpdateCartItem(w http.ResponseWriter, r *http.Request) {
	var cartItem models.CartItem
	if err := json.NewDecoder(r.Body).Decode(&cartItem); err != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(primitive.ObjectID) // Pastikan middleware telah menambahkan userID ke context

	// Update the item in the cart
	_, err := cartCollection.UpdateOne(context.TODO(),
		bson.M{"user_id": userID, "items.product_id": cartItem.ProductID},
		bson.M{"$set": bson.M{"items.$.quantity": cartItem.Quantity}})

	if err != nil {
		http.Error(w, "Failed to update cart item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RemoveCartItem removes an item from the cart
func RemoveCartItem(w http.ResponseWriter, r *http.Request) {
	productIDStr := r.URL.Query().Get("product_id")
	productID, err := primitive.ObjectIDFromHex(productIDStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(primitive.ObjectID) // Pastikan middleware telah menambahkan userID ke context

	// Remove the item from the cart
	_, err = cartCollection.UpdateOne(context.TODO(),
		bson.M{"user_id": userID},
		bson.M{"$pull": bson.M{"items": bson.M{"product_id": productID}}})

	if err != nil {
		http.Error(w, "Failed to remove item from cart", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
