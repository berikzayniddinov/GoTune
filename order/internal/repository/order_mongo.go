package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gotune/order/internal/entity"
)

type OrderRepository interface {
	Create(ctx context.Context, order *entity.Order) (primitive.ObjectID, error)
	GetByUserID(ctx context.Context, userID primitive.ObjectID) ([]entity.Order, error)
	Delete(ctx context.Context, orderID primitive.ObjectID, userID primitive.ObjectID) error
}

type orderRepository struct {
	collection *mongo.Collection
}

func NewOrderRepository(db *mongo.Database) OrderRepository {
	return &orderRepository{
		collection: db.Collection("orders"),
	}
}

func (r *orderRepository) Create(ctx context.Context, order *entity.Order) (primitive.ObjectID, error) {
	res, err := r.collection.InsertOne(ctx, order)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func (r *orderRepository) GetByUserID(ctx context.Context, userID primitive.ObjectID) ([]entity.Order, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []entity.Order
	for cursor.Next(ctx) {
		var order entity.Order
		if err := cursor.Decode(&order); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *orderRepository) Delete(ctx context.Context, orderID primitive.ObjectID, userID primitive.ObjectID) error {
	filter := bson.M{
		"_id":     orderID,
		"user_id": userID,
	}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}
