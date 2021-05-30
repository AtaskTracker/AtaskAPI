package labelRep

import (
	"context"
	"github.com/AtaskTracker/AtaskAPI/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type LabelRep struct {
	mongo *mongo.Client
}

const collectionName = "tasks"
const dbName = "atasktracker"

func New(mongo *mongo.Client) *LabelRep {
	return &LabelRep{mongo: mongo}
}

func (rep *LabelRep) GetLabels(userId string) ([]dto.Label, error) {
	var collection = rep.mongo.Database(dbName).Collection(collectionName)
	var result, err = collection.Find(context.Background(), bson.M{"participants": userId})
	if err != nil {
		return nil, err
	}
	var tasks []dto.Task
	if err := result.All(context.Background(), &tasks); err != nil {
		return nil, err
	}
	var labelSet []dto.Label
	for _, task := range tasks {
		for _, label := range task.Labels {
			if !isLabelInSet(labelSet, label) {
				labelSet = append(labelSet, label)
			}
		}
	}
	return labelSet, nil
}

func isLabelInSet(labelSet []dto.Label, label dto.Label) bool {
	for _, existingLabel := range labelSet {
		if existingLabel.Summary == label.Summary {
			return true
		}
	}
	return false
}
