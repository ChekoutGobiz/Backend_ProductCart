package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/kuyjajan/kuyjajan-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var client *mongo.Client

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

	fmt.Println("MongoDB connection established successfully!")
}

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	user.Password = string(hashedPassword)

	// Insert user to database
	collection := client.Database("jajankuy").Collection("users")
	_, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Println("Error inserting user:", err) // Log the error
		http.Error(w, "Error saving user", http.StatusInternalServerError)
		return
	}

	// Return the created user in response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User successfully registered",
		"user":    user,
	})
}

// Login function
// Login function
func Login(w http.ResponseWriter, r *http.Request) {
	var loginData models.User
	var storedUser models.User

	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	collection := client.Database("jajankuy").Collection("users")
	err := collection.FindOne(context.TODO(), bson.M{"email": loginData.Email}).Decode(&storedUser)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Cek password
	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(loginData.Password)); err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// Generate JWT Token
	token, err := generateJWT(storedUser.Email)
	if err != nil {
		log.Println("Error generating JWT token:", err)
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Log the generated token
	log.Println("Generated JWT token:", token)

	// Kirimkan token bersama user data
	response := map[string]interface{}{
		"user":  storedUser,
		"token": token,
	}
	json.NewEncoder(w).Encode(response)
}

// generateJWT generates a JWT token for the given email
func generateJWT(email string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	// Set claims untuk token
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Expired dalam 24 jam

	// Ambil JWT_SECRET dari .env
	jwtSecret := os.Getenv("JWT_SECRET")

	// Buat token
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		log.Println("Error saat membuat JWT:", err)
		return "", err
	}

	return tokenString, nil
}
