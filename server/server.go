package server

import "net/http"
import "github.com/gorilla/mux"

type server struct {
	router *mux.Router
	//TODO: add services here
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
}

func NewServer() *server {
	server := &server{
		router: mux.NewRouter(),
	}

	server.ConfigureRouter()
	return server
}
