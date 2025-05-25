package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"gotune/users/intern/entity"
)

type UserRepository interface {
	Create(ctx mongo.SessionContext, user *entity.User) error
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	GetAll(ctx context.Context) ([]entity.User, error)
	FindByID(ctx context.Context, id string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id string) error
	DeleteAll(ctx context.Context) error // Новый метод
	ConfirmUser(ctx context.Context, id string) error
	GetDatabase() *mongo.Database
}

type userRepository struct {
	collection *mongo.Collection
	db         *mongo.Database
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &userRepository{
		collection: db.Collection("users"),
		db:         db,
	}
}
func (r *userRepository) GetDatabase() *mongo.Database {
	return r.db
}
func (r *userRepository) Create(ctx mongo.SessionContext, user *entity.User) error {
	user.CreatedAt = time.Now().Unix()
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetAll(ctx context.Context) ([]entity.User, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []entity.User
	for cursor.Next(ctx) {
		var user entity.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user entity.User
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	objID, err := primitive.ObjectIDFromHex(user.ID.Hex())
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"username":   user.Username,
			"email":      user.Email,
			"password":   user.Password,
			"updated_at": time.Now().Unix(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}
func (r *userRepository) DeleteAll(ctx context.Context) error {
	_, err := r.collection.DeleteMany(ctx, bson.M{})
	return err
}
func (r *userRepository) ConfirmUser(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"confirmed":  true,
			"updated_at": time.Now().Unix(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}
