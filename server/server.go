package server

import (
	"github.com/AtaskTracker/AtaskAPI/database/taskRep"
	"github.com/AtaskTracker/AtaskAPI/database/userRepo"
	"github.com/AtaskTracker/AtaskAPI/handlers/taskHandler"
	"github.com/AtaskTracker/AtaskAPI/handlers/userHandler"
	"github.com/AtaskTracker/AtaskAPI/services/taskService"
	"github.com/AtaskTracker/AtaskAPI/services/userService"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)
import "github.com/gorilla/mux"

type server struct {
	router *mux.Router
	//TODO: add handlers here
	taskHandler *taskHandler.TaskHandler
	userHandler *userHandler.UserHandler
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
	s.router.HandleFunc("/task", s.userHandler.AuthorizationMW(s.taskHandler.CreateTask)).Methods("POST")
	s.router.HandleFunc("/task", s.userHandler.AuthorizationMW(s.taskHandler.GetUserTasks)).Methods("GET")
	s.router.HandleFunc("/task/{id}", s.userHandler.AuthorizationMW(s.taskHandler.GetTasksByUserId)).Methods("GET")
	s.router.HandleFunc("/task/{id}", s.userHandler.AuthorizationMW(s.taskHandler.DeleteByUserId)).Methods("DELETE")
	s.router.HandleFunc("/task/{taskId}", s.userHandler.AuthorizationMW(s.taskHandler.UpdateTask)).Methods("PUT")

	s.router.HandleFunc("/auth/logout", s.userHandler.Logout).Methods("POST")
	s.router.HandleFunc("/auth/google", s.userHandler.Login).Methods("POST")

	s.router.HandleFunc("/user", s.userHandler.GetUserByEmail).Methods("GET")

}

func NewServer(mongoClient *mongo.Client, redis *redis.Client) *server {
	server := &server{
		router:      mux.NewRouter(),
		taskHandler: taskHandler.New(taskService.New(taskRep.New(mongoClient))),
		userHandler: userHandler.New(userService.New(userRepo.New(mongoClient), redis)),
	}

	server.ConfigureRouter()
	return server
}
