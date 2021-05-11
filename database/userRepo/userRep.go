package userRepo

import (
	"context"
	"github.com/AtaskTracker/AtaskAPI/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo struct {
	mongo *mongo.Client
}

func New(mongo *mongo.Client) *UserRepo {
	return &UserRepo{mongo: mongo}
}

func (rep *UserRepo) CreateUser(user dto.User) (dto.User, error) {
	user.UUID = primitive.NewObjectID()
	var collection = rep.mongo.Database("atasktracker").Collection("users")
	var result, err = collection.InsertOne(context.Background(), user)
	if err != nil {
		return user, err
	}
	user.UUID = result.InsertedID.(primitive.ObjectID)
	return user, nil
}
