package taskRep

import (
	"context"
	"github.com/AtaskTracker/AtaskAPI/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type TaskRep struct {
	mongo *mongo.Client
}

const collectionName = "tasks"
const dbName = "atasktracker"

func New(mongo *mongo.Client) *TaskRep {
	return &TaskRep{mongo: mongo}
}

func (rep *TaskRep) CreateTask(task dto.Task) (dto.Task, error) {
	task.UUID = primitive.NewObjectID()
	var collection = rep.mongo.Database(dbName).Collection(collectionName)
	var result, err = collection.InsertOne(context.Background(), task)
	if err != nil {
		return task, err
	}
	task.UUID = result.InsertedID.(primitive.ObjectID)
	return task, nil
}

func (rep *TaskRep) GetAll() ([]dto.Task, error) {
	var collection = rep.mongo.Database(dbName).Collection(collectionName)
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
	var collection = rep.mongo.Database(dbName).Collection(collectionName)
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

func (rep *TaskRep) GetById(taskId string) (*dto.Task, error) {
	var collection = rep.mongo.Database(dbName).Collection(collectionName)
	var result bson.M
	objectId, _ := primitive.ObjectIDFromHex(taskId)
	if err := collection.FindOne(context.Background(), bson.M{"_id": objectId}).Decode(&result); err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return nil, nil
		default:
			return nil, err
		}
	}

	var task *dto.Task
	bsonBytes, _ := bson.Marshal(result)
	bson.Unmarshal(bsonBytes, &task)
	return task, nil
}

func (rep *TaskRep) UpdateById(newTask dto.Task) error {
	var collection = rep.mongo.Database(dbName).Collection(collectionName)
	var err = collection.FindOneAndReplace(context.Background(), bson.M{"_id": newTask.UUID}, newTask).Err()
	return err
}

func (rep *TaskRep) DeleteById(id string) error {
	var objId, err = primitive.ObjectIDFromHex(id)
	var collection = rep.mongo.Database(dbName).Collection(collectionName)
	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": objId})
	return err
}

//TODO: ?????????????? ?????? ?????????? ???????????? ???????????????? (???????? ?????? ???? ???????????????? ??\_(???)_/??)
func (rep *TaskRep) AddLabel(taskId string, label dto.Label) error {
	var collection = rep.mongo.Database(dbName).Collection(collectionName)
	var result bson.M
	if err := collection.FindOne(context.Background(), bson.M{"_id": taskId}).Decode(&result); err != nil {
		return err
	}
	result["labels"] = append(result["labels"].([]dto.Label), label)
	var err = collection.FindOneAndReplace(context.Background(), bson.M{"_id": taskId}, result).Err()
	return err
}

func (rep *TaskRep) AddParticipant(taskId string, email string) error {
	var collection = rep.mongo.Database(dbName).Collection(collectionName)
	var result bson.M
	if err := collection.FindOne(context.Background(), bson.M{"id": taskId}).Decode(&result); err != nil {
		return err
	}
	result["participants"] = append(result["labels"].([]string), email)
	var err = collection.FindOneAndReplace(context.Background(), bson.M{"id": taskId}, result).Err()
	return err
}

func (rep *TaskRep) GetWithFilter(userEmail string, dateTo time.Time, dateFrom time.Time, label string) ([]dto.Task, error) {
	filter := rep.createFilter(userEmail, dateTo, dateFrom, label)
	var collection = rep.mongo.Database(dbName).Collection(collectionName)
	var result, err = collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	var tasks []dto.Task
	if err := result.All(context.Background(), &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (rep *TaskRep) createFilter(userEmail string, dateTo time.Time, dateFrom time.Time, label string) bson.M {
	filter := []bson.M{
		{"participants": userEmail},
	}
	if !dateTo.IsZero() {
		filter = append(filter, bson.M{"date": bson.M{
			"$lt": dateTo,
		}})
	}
	if !dateFrom.IsZero() {
		filter = append(filter, bson.M{"date": bson.M{
			"$gte": dateFrom,
		}})
	}
	if label != "" {
		filter = append(filter, bson.M{"labels.summary": label})
	}
	return bson.M{"$and": filter}
}
