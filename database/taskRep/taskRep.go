package taskRep

import (
	"context"
	"github.com/AtaskTracker/AtaskAPI/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskRep struct {
	mongo *mongo.Client
}

func New(mongo *mongo.Client) *TaskRep {
	return &TaskRep{mongo: mongo}
}

func (rep *TaskRep) CreateTask(task dto.Task) (dto.Task, error) {
	task.UUID = primitive.NewObjectID()
	var collection = rep.mongo.Database("atasktracker").Collection("tasks")
	var result, err = collection.InsertOne(context.Background(), task)
	if err != nil {
		return task, err
	}
	task.UUID = result.InsertedID.(primitive.ObjectID)
	return task, nil
}

func (rep *TaskRep) GetAll() ([]dto.Task, error) {
	var collection = rep.mongo.Database("atasktracker").Collection("tasks")
	var result, err = collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	var tasks []dto.Task
	if err := result.All(context.Background(), &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (rep *TaskRep) GetByUserId(id string) ([]dto.Task, error) {
	var collection = rep.mongo.Database("atasktracker").Collection("tasks")
	var result, err = collection.Find(context.Background(), bson.M{"participants": id})
	if err != nil {
		return nil, err
	}
	defer result.Close(context.Background())
	var tasks []dto.Task
	if err := result.All(context.Background(), &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (rep *TaskRep) UpdateById(newTask dto.Task) error {
	var collection = rep.mongo.Database("atasktracker").Collection("tasks")
	var err = collection.FindOneAndReplace(context.Background(), bson.M{"_id": newTask.UUID}, newTask).Err()
	return err
}

func (rep *TaskRep) DeleteById(id string) error {
	var objId, err = primitive.ObjectIDFromHex(id)
	var collection = rep.mongo.Database("atasktracker").Collection("tasks")
	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": objId})
	return err
}

//TODO: сделать так чтобы лейблы работали (пока они не работают ¯\_(ツ)_/¯)
func (rep *TaskRep) AddLabel(taskId string, labelId string) error {
	var collection = rep.mongo.Database("atasktracker").Collection("tasks")
	var result bson.M
	if err := collection.FindOne(context.Background(), bson.M{"_id": taskId}).Decode(&result); err != nil {
		return err
	}
	result["labels"] = append(result["labels"].([]string), labelId)
	var err = collection.FindOneAndReplace(context.Background(), bson.M{"_id": taskId}, result).Err()
	return err
}

func (rep *TaskRep) AddParticipant(taskId string, email string) error {
	var collection = rep.mongo.Database("atasktracker").Collection("tasks")
	var result bson.M
	if err := collection.FindOne(context.Background(), bson.M{"id": taskId}).Decode(&result); err != nil {
		return err
	}
	result["participants"] = append(result["labels"].([]string), email)
	var err = collection.FindOneAndReplace(context.Background(), bson.M{"id": taskId}, result).Err()
	return err
}
