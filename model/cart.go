package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CartItem represents an item in the cart
type CartItem struct {
	ProductID primitive.ObjectID `json:"product_id,omitempty" bson:"product_id,omitempty"`
	Quantity  int                `json:"quantity,omitempty" bson:"quantity,omitempty"`
}

// Cart represents a shopping cart
type Cart struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID    primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Items     []CartItem         `json:"items,omitempty" bson:"items,omitempty"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

// AddItem adds an item to the cart
func (c *Cart) AddItem(item CartItem) {
	for i, cartItem := range c.Items {
		if cartItem.ProductID == item.ProductID {
			c.Items[i].Quantity += item.Quantity
			return
		}
	}
	c.Items = append(c.Items, item)
}

// RemoveItem removes an item from the cart by its ProductID
func (c *Cart) RemoveItem(productID primitive.ObjectID) {
	for i, cartItem := range c.Items {
		if cartItem.ProductID == productID {
			c.Items = append(c.Items[:i], c.Items[i+1:]...) // Remove item
			return
		}
	}
}
