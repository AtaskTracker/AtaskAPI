package dto

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	UUID       primitive.ObjectID `bson:"_id" json:"uuid"`
	Name       string             `bson:"name" json:"name"`
	Email      string             `bson:"email" json:"email"`
	PictureUrl string             `bson:"picture_url" json:"picture_url"`
}
