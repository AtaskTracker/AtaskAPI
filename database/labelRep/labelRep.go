package labelRep

import (
	"context"
	"fmt"
	"github.com/AtaskTracker/AtaskAPI/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/appengine/log"
)

type LabelRep struct {
	mongo *mongo.Client
}

func New(mongo *mongo.Client) *LabelRep {
	return &LabelRep{mongo: mongo}
}

func (rep *LabelRep) CreateLabel(label dto.Label) (dto.Label, error) {
	if _, found := rep.GetBySummary(label.Summary); found {
		return label, fmt.Errorf("label already exist: %s", label.Summary)
	}
	label.UUID = primitive.NewObjectID()
	var collection = rep.mongo.Database("atasktracker").Collection("labels")
	var result, err = collection.InsertOne(context.Background(), label)
	if err != nil {
		return label, err
	}
	label.UUID = result.InsertedID.(primitive.ObjectID)
	return label, nil
}

func (rep *LabelRep) GetBySummary(summary string) (dto.Label, bool) {
	var collection = rep.mongo.Database("atasktracker").Collection("labels")
	var label bson.M
	if err := collection.FindOne(context.Background(), bson.M{"summary": summary}).Decode(&label); err != nil {
		if err != mongo.ErrNoDocuments {
			log.Warningf(context.Background(), "labelRep/GetBySummary", err)
		}
		return dto.Label{}, false
	}
	var labelStruct dto.Label
	bsonBytes, _ := bson.Marshal(label)
	bson.Unmarshal(bsonBytes, &labelStruct)
	return labelStruct, true
}

func (rep *LabelRep) GetById(id string) (dto.Label, bool) {
	var collection = rep.mongo.Database("atasktracker").Collection("labels")
	var objId, err = primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Warningf(context.Background(), "labelRep/GetBySummary", fmt.Errorf("failed to convert string %s to uuid", id))
		return dto.Label{}, false
	}
	var label bson.M
	if err := collection.FindOne(context.Background(), bson.M{"id": objId}).Decode(&label); err != nil {
		if err != mongo.ErrNoDocuments {
			log.Warningf(context.Background(), "labelRep/GetBySummary", err)
		}
		return dto.Label{}, false
	}
	var labelStruct dto.Label
	bsonBytes, _ := bson.Marshal(label)
	bson.Unmarshal(bsonBytes, &labelStruct)
	return labelStruct, true
}

func (rep *LabelRep) GetAll() ([]dto.Label, error) {
	var collection = rep.mongo.Database("atasktracker").Collection("labels")
	var result, err = collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	var tasks []dto.Label
	if err := result.All(context.Background(), &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}
