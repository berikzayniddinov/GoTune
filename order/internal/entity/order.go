package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type OrderItem struct {
	InstrumentID primitive.ObjectID `bson:"instrument_id"`
	Quantity     int32              `bson:"quantity"`
}

type Order struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id"`
	Items     []OrderItem        `bson:"items"`
	CreatedAt time.Time          `bson:"created_at"`
}
