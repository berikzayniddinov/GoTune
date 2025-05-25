package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Username  string             `bson:"username"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	CreatedAt int64              `bson:"created_at"`
	UpdatedAt int64              `bson:"updated_at,omitempty"`
	Confirmed bool               `bson:"confirmed"` // новое поле
}
