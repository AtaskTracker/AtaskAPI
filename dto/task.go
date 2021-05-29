package dto

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Task struct {
	UUID         primitive.ObjectID `bson:"_id" json:"uuid"`
	Summary      string             `bson:"summary" json:"summary"`
	Description  string             `bson:"description" json:"description"`
	Photo        string             `bson:"photo" json:"photo"`
	Status       string             `bson:"status" json:"status"`
	Date         time.Time          `bson:"date" json:"date"` //"date": "2021-10-10 12:00"
	Participants []string           `bson:"participants" json:"participants"`
	Labels       []Label            `bson:"labels" json:"labels"`
}
