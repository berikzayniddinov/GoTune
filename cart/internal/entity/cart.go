package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type CartItem struct {
	InstrumentID primitive.ObjectID `bson:"instrument_id"`
	Quantity     int32              `bson:"quantity"`
}

type Cart struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	UserID primitive.ObjectID `bson:"user_id"`
	Items  []CartItem         `bson:"items"`
}
