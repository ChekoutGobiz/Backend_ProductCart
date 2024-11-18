package config

import (
	"log"
	"os"

	models "github.com/ChekoutGobiz/BackendChekout/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Mongo connection string directly from the environment
var MongoString string = os.Getenv("MONGOSTRING")

// DBInfo struct now uses MongoString directly, no SRVLookup required
var mongoinfo = models.DBInfo{
	DBString: MongoString,
	DBName:   "jajankuy",
}

// Mongo connection using the provided string
var Mongoconn *mongo.Client

// Initialize MongoDB connection
func init() {
	if MongoString == "" {
		log.Fatal("MongoDB connection string (MONGOSTRING) is required")
	}

	clientOptions := options.Client().ApplyURI(MongoString)
	client, err := mongo.Connect(nil, clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Optionally, check if the connection is successful
	err = client.Ping(nil, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	// Store the client reference in Mongoconn
	Mongoconn = client

	log.Println("MongoDB connection established successfully")
}
