package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type Label struct {
	UUID    primitive.ObjectID `bson:"_id" json:"uuid"`
	Summary string             `json:"summary"`
	Color   string             `json:"color"`
}
