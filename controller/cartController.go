package controllers

import (
	"context"
	"log"
	"os"

	models "github.com/ChekoutGobiz/BackendChekout/model"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
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

	// Initialize cart collection
	cartCollection = client.Database("jajankuy").Collection("carts")
}

// AddToCart adds an item to the user's cart
func AddToCart(c *fiber.Ctx) error {
	// Struktur untuk membaca request body
	var requestBody struct {
		UserID    string `json:"user_id"`
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	}

	// Dekode JSON body
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validasi userID
	if requestBody.UserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	// Konversi userID dari string ke ObjectID
	userID, err := primitive.ObjectIDFromHex(requestBody.UserID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	// Validasi productID dan konversi ke ObjectID
	if requestBody.ProductID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Product ID is required",
		})
	}
	productID, err := primitive.ObjectIDFromHex(requestBody.ProductID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID format",
		})
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
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve or create cart",
			})
		}
	} else {
		// Jika keranjang ada, tambahkan item ke dalamnya
		cart.Items = append(cart.Items, cartItem)
		_, err = cartCollection.UpdateOne(context.TODO(), bson.M{"user_id": userID}, bson.M{"$set": bson.M{"items": cart.Items}})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add item to cart",
		})
	}

	// Kirimkan respons
	return c.JSON(cart)
}

// GetCart retrieves the user's cart
func GetCart(c *fiber.Ctx) error {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	// Konversi userID dari string ke ObjectID
	objectID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	var cart models.Cart
	err = cartCollection.FindOne(context.TODO(), bson.M{"user_id": objectID}).Decode(&cart)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Cart not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve cart",
		})
	}

	return c.JSON(cart)
}

// UpdateCartItem updates the quantity of a specific item in the cart
func UpdateCartItem(c *fiber.Ctx) error {
	var cartItem models.CartItem
	if err := c.BodyParser(&cartItem); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid data",
		})
	}

	userID := c.Locals("userID").(primitive.ObjectID) // Pastikan middleware telah menambahkan userID ke context

	// Update the item in the cart
	_, err := cartCollection.UpdateOne(context.TODO(),
		bson.M{"user_id": userID, "items.product_id": cartItem.ProductID},
		bson.M{"$set": bson.M{"items.$.quantity": cartItem.Quantity}})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update cart item",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// RemoveCartItem removes an item from the cart
func RemoveCartItem(c *fiber.Ctx) error {
	productIDStr := c.Query("product_id")
	productID, err := primitive.ObjectIDFromHex(productIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	userID := c.Locals("userID").(primitive.ObjectID) // Pastikan middleware telah menambahkan userID ke context

	// Remove the item from the cart
	_, err = cartCollection.UpdateOne(context.TODO(),
		bson.M{"user_id": userID},
		bson.M{"$pull": bson.M{"items": bson.M{"product_id": productID}}})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to remove item from cart",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
