package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type Label struct {
	UUID    primitive.ObjectID `bson:"_id" json:"uuid"`
	Summary string             `bson:"summary" json:"summary"`
	Color   string             `bson:"summary" json:"color"`
}
