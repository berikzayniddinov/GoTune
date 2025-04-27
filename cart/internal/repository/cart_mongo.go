package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gotune/cart/internal/entity"
)

type CartRepository interface {
	AddToCart(ctx context.Context, userID, instrumentID primitive.ObjectID, quantity int32) error
	GetCart(ctx context.Context, userID primitive.ObjectID) ([]entity.CartItem, error)
	RemoveFromCart(ctx context.Context, userID, instrumentID primitive.ObjectID) error
	ClearCart(ctx context.Context, userID primitive.ObjectID) error
}

type cartRepository struct {
	collection *mongo.Collection
}

func NewCartRepository(db *mongo.Database) CartRepository {
	return &cartRepository{
		collection: db.Collection("carts"),
	}
}

func (r *cartRepository) AddToCart(ctx context.Context, userID, instrumentID primitive.ObjectID, quantity int32) error {
	filter := bson.M{"user_id": userID}
	update := bson.M{
		"$push": bson.M{
			"items": bson.M{
				"instrument_id": instrumentID,
				"quantity":      quantity,
			},
		},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	return err
}

func (r *cartRepository) GetCart(ctx context.Context, userID primitive.ObjectID) ([]entity.CartItem, error) {
	var cart entity.Cart
	err := r.collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&cart)
	if err != nil {
		return nil, err
	}
	return cart.Items, nil
}

func (r *cartRepository) RemoveFromCart(ctx context.Context, userID, instrumentID primitive.ObjectID) error {
	filter := bson.M{"user_id": userID}
	update := bson.M{
		"$pull": bson.M{
			"items": bson.M{
				"instrument_id": instrumentID,
			},
		},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *cartRepository) ClearCart(ctx context.Context, userID primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"user_id": userID})
	return err
}
