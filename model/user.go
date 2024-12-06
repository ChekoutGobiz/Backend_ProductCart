package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
}

type BlacklistedToken struct {
    ID        primitive.ObjectID `bson:"_id,omitempty"`
    Token     string             `bson:"token"`
    CreatedAt int64              `bson:"created_at"`
}

func CreateUser(user *User, collection *mongo.Collection) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return collection.InsertOne(ctx, user)
}

func FindUserByEmail(email string, collection *mongo.Collection) *mongo.SingleResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return collection.FindOne(ctx, bson.M{"email": email})
}
