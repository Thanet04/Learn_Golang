package repository

import (
	"context"
	"learn_golang/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(client *mongo.Client, dbName string) *UserRepository {
	collection := client.Database(dbName).Collection("users")
	return &UserRepository{collection: collection}
}

func (repo *UserRepository) CreateUser(ctx context.Context, user model.User) (*mongo.InsertOneResult, error) {
	return repo.collection.InsertOne(ctx, user)
}

func (repo *UserRepository) GetUser(ctx context.Context, id string) (*model.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var user model.User
	if err := repo.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) UpdateUser(ctx context.Context, id string, user model.User) (*mongo.UpdateResult, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objectID}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "name", Value: user.Name}, {Key: "email", Value: user.Email}, {Key: "password", Value: user.Password}, {Key: "age", Value: user.Age}}}}
	return repo.collection.UpdateOne(ctx, filter, update)
}

func (repo *UserRepository) DeleteUser(ctx context.Context, id string) (*mongo.DeleteResult, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objectID}
	return repo.collection.DeleteOne(ctx, filter)
}

func (repo *UserRepository) GetAllUsers(ctx context.Context) ([]model.User, error) {
	var users []model.User
	cursor, err := repo.collection.Find(ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var user model.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (repo *UserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	filter := bson.M{"email": email}

	var user model.User
	err := repo.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}

	return &user, nil
}
