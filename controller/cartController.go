package controllers

import (
	"context"
	"log"
	"os"
	"time"

	models "github.com/ChekoutGobiz/BackendChekout/model"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	cartCollection     *mongo.Collection
	productsCollection *mongo.Collection // Tambahkan ini untuk akses koleksi produk
)

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

	// Initialize products collection
	productsCollection = client.Database("jajankuy").Collection("products")
}

// AddItemToCart adds an item to the user's cart
func AddItemToCart(c *fiber.Ctx) error {
	// Parse request body
	var cartItem models.CartItem
	if err := c.BodyParser(&cartItem); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request data",
		})
	}

	// Get user_id from query parameter
	userID := c.Query("user_id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	// Convert user_id to ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	// Find the cart for the user
	var cart models.Cart
	err = cartCollection.FindOne(context.Background(), bson.M{"user_id": userObjectID}).Decode(&cart)
	if err != nil {
		// If cart doesn't exist, create a new cart
		cart = models.Cart{
			UserID:    userObjectID,
			Items:     []models.CartItem{},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	}

	// Add the item to the cart
	cart.AddItem(cartItem)

	// Save the updated cart to the database
	_, err = cartCollection.UpdateOne(
		context.Background(),
		bson.M{"user_id": userObjectID},
		bson.M{"$set": bson.M{"items": cart.Items, "updated_at": time.Now()}})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update cart",
		})
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Item added to cart successfully",
		"cart":    cart,
	})
}

// GetCart retrieves the user's cart with total price
func GetCart(c *fiber.Ctx) error {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	// Convert userID from string to ObjectID
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

	// Calculate the total price of the cart
	totalPrice, err := cart.TotalPrice(productsCollection)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to calculate total price",
		})
	}

	// Send the cart details along with the total price
	return c.JSON(fiber.Map{
		"cart":        cart,
		"total_price": totalPrice,
	})
}

// UpdateCartItem updates the quantity of a specific item in the cart
func UpdateCartItem(c *fiber.Ctx) error {
	var cartItem models.CartItem
	if err := c.BodyParser(&cartItem); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid data",
		})
	}

	// Get userID from context (middleware should add userID to context)
	userID := c.Locals("userID").(primitive.ObjectID)

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

// RemoveItemFromCart removes an item from the cart
func RemoveItemFromCart(c *fiber.Ctx) error {
	// Ambil product_id dari parameter path
	productIDStr := c.Params("product_id") // Mengambil product_id dari path URL
	productID, err := primitive.ObjectIDFromHex(productIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID format",
		})
	}

	// Ambil userID dari context (pastikan middleware telah menambahkan userID ke context)
	userID := c.Locals("userID").(primitive.ObjectID)

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
// RemoveCartItem removes an item from the cart
func RemoveCartItem(c *fiber.Ctx) error {
	productIDStr := c.Query("product_id")
	productID, err := primitive.ObjectIDFromHex(productIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	// Ambil userID dari context (pastikan middleware telah menambahkan userID ke context)
	userID := c.Locals("userID").(primitive.ObjectID)

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
