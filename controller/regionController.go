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
func CreateRegion(c *fiber.Ctx) error {
	var region models.Region
	if err := c.BodyParser(&region); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid data",
		})
	}

	// Generate new ObjectID for region
	region.ID = primitive.NewObjectID()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := regionCollection.InsertOne(ctx, region)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to insert region",
		})
	}

	// Mengembalikan respons dalam format JSON
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"region": region,
	})
}

// GetRegions retrieves all regions from the database
func GetRegions(c *fiber.Ctx) error {
	var regions []models.Region
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := regionCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get regions",
		})
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var region models.Region
		if err := cursor.Decode(&region); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to decode region",
			})
		}
		regions = append(regions, region)
	}

	// Mengembalikan hasil dalam format JSON
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"regions": regions,
	})
}
