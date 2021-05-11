package server

import (
	"github.com/AtaskTracker/AtaskAPI/database/taskRep"
	"github.com/AtaskTracker/AtaskAPI/handlers/taskHandler"
	"github.com/AtaskTracker/AtaskAPI/services/taskService"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)
import "github.com/gorilla/mux"

type server struct {
	router *mux.Router
	//TODO: add handlers here
	taskHandler *taskHandler.TaskHandler
}

type Config struct {
	MongoConnection string
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) Start(port string) {
	http.ListenAndServe(port, s.router)
}

func (s *server) ConfigureRouter() {
	s.router.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("hello world"))
	})

	//TODO: add endpoints here
	s.router.HandleFunc("/task", s.taskHandler.CreateTask).Methods("POST")
	s.router.HandleFunc("/task/{id}", s.taskHandler.GetTasksByUserId).Methods("GET")
	s.router.HandleFunc("/task/{id}", s.taskHandler.DeleteByUserId).Methods("DELETE")
}

func NewServer(mongoClient *mongo.Client) *server {
	server := &server{
		router:      mux.NewRouter(),
		taskHandler: taskHandler.New(taskService.New(taskRep.New(mongoClient))),
	}

	server.ConfigureRouter()
	return server
}
