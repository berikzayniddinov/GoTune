package repository

import (
	"context"
	"gotune/instruments/internal/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type InstrumentRepository interface {
	Create(ctx context.Context, instrument *entity.Instrument) (primitive.ObjectID, error)
	GetAll(ctx context.Context) ([]entity.Instrument, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	UpdateByID(ctx context.Context, id primitive.ObjectID, instrument *entity.Instrument) error
}

type instrumentRepository struct {
	collection *mongo.Collection
}

func NewInstrumentRepositories(db *mongo.Database) InstrumentRepository {
	return &instrumentRepository{
		collection: db.Collection("instruments"),
	}
}

func (r *instrumentRepository) Create(ctx context.Context, instrument *entity.Instrument) (primitive.ObjectID, error) {
	res, err := r.collection.InsertOne(ctx, instrument)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func (r *instrumentRepository) GetAll(ctx context.Context) ([]entity.Instrument, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var instruments []entity.Instrument
	for cursor.Next(ctx) {
		var inst entity.Instrument
		if err := cursor.Decode(&inst); err != nil {
			return nil, err
		}
		instruments = append(instruments, inst)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return instruments, nil
}

func (r *instrumentRepository) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *instrumentRepository) UpdateByID(ctx context.Context, id primitive.ObjectID, instrument *entity.Instrument) error {
	update := bson.M{
		"$set": bson.M{
			"name":        instrument.Name,
			"description": instrument.Description,
			"price":       instrument.Price,
		},
	}
	_, err := r.collection.UpdateByID(ctx, id, update)
	return err
}
