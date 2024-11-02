package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
    ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
    Name      string             `json:"name,omitempty" bson:"name,omitempty"`
    RegionID  primitive.ObjectID `json:"region_id,omitempty" bson:"region_id,omitempty"`
    Price     float64            `json:"price,omitempty" bson:"price,omitempty"`
    // Tambahkan field lain sesuai kebutuhan
}

