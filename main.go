package main

import (
	"context"
	"github.com/AtaskTracker/AtaskAPI/server"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		"mongodb+srv://newUser:newUserPassword@cluster0.6t1sr.mongodb.net/atasktracker?retryWrites=true&w=majority",
	))
	if err != nil {
		log.Fatal(err)
	}
	redis := redis.NewClient(&redis.Options{
		Addr:     "ec2-52-213-88-26.eu-west-1.compute.amazonaws.com:21929",
		Password: "p421db2c77453e5864bcb5421e8323598a2679a179e529fc9ddcccde150cf8bf1",
		DB:       0,
	})
	if res := redis.Ping(context.Background()); res.Err() != nil {
		log.Fatal("redis failed: ", res.Err())
	}
	server := server.NewServer(client, redis)
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	server.Start(":" + port)

	//var rep = taskRep.New(client)
	//var uuid, _ = primitive.ObjectIDFromHex("609985f86ee56ee9c5b68542")
	//var task = dto.Task{
	//	UUID:         uuid,
	//	Summary:      "UPDATED",
	//	Description:  "desc",
	//	Photo:        "phphph",
	//	Status:       "in progress",
	//	Date:         time.Now(),
	//	Participants: nil,
	//	Labels:       nil,
	//}
	////taskHandler, _ = rep.CreateTask(taskHandler)
	////var jsonbytes, _ = json.Marshal(taskHandler)
	////fmt.Println(string(jsonbytes))
	//err = rep.DeleteById("609985f86ee56ee9c5b68542")
	//
	//fmt.Print(task)
	//fmt.Print(err)
}
