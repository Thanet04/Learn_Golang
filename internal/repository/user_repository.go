package repository

import (
	"context"
	"learn_golang/internal/model"

	"go.mongodb.org/mongo-driver/bson"
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
