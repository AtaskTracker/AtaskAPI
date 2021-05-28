package main

import (
	"context"
	"github.com/AtaskTracker/AtaskAPI/server"
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
	server := server.NewServer(client)
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
