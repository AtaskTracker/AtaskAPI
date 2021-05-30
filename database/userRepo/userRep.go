package userRepo

import (
	"context"
	"github.com/AtaskTracker/AtaskAPI/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo struct {
	mongo *mongo.Client
}

const collectionName = "users"
const dbName = "atasktracker"

func New(mongo *mongo.Client) *UserRepo {
	return &UserRepo{mongo: mongo}
}

func (rep *UserRepo) CreateUser(user dto.User) (dto.User, error) {
	user.UUID = primitive.NewObjectID()
	var result, err = rep.mongo.
		Database(dbName).
		Collection(collectionName).
		InsertOne(context.Background(), user)
	if err != nil {
		return user, err
	}
	user.UUID = result.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (rep *UserRepo) GetUserByEmail(email string) (dto.User, error) {
	user := dto.User{}
	filter := bson.D{{"email", email}}
	err := rep.mongo.
		Database(dbName).
		Collection(collectionName).
		FindOne(context.Background(), filter).
		Decode(&user)
	if err != nil && err != mongo.ErrNoDocuments {
		return user, err
	}
	return user, nil
}

func (rep *UserRepo) GetUserById(userId string) (dto.User, error) {
	user := dto.User{}
	objectId, _ := primitive.ObjectIDFromHex(userId)
	filter := bson.M{"_id": objectId}
	err := rep.mongo.
		Database(dbName).
		Collection(collectionName).
		FindOne(context.Background(), filter).
		Decode(&user)
	if err != nil && err != mongo.ErrNoDocuments {
		return user, err
	}
	return user, nil
}

func (rep *UserRepo) UpdateUser(user dto.User) (dto.User, error) {
	filter := bson.D{{"_id", user.UUID}}
	_, err := rep.mongo.Database(dbName).Collection(collectionName).
		ReplaceOne(
			context.Background(),
			filter,
			bson.M{
				"name":        user.Name,
				"picture_url": user.PictureUrl,
			})
	if err != nil {
		return dto.User{}, err
	}
	return user, nil
}
