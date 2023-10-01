package webserver

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
)

type WebServer struct {
	Router        chi.Router
	Handlers      map[string]http.HandlerFunc
	WebServerPort string
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      make(map[string]http.HandlerFunc),
		WebServerPort: serverPort,
	}
}

func (s *WebServer) AddHandler(path string, handler http.HandlerFunc) {
	s.Handlers[path] = handler
}

func (s *WebServer) Start() {
	s.Router.Use(middleware.Logger)

	for path, handler := range s.Handlers {
		s.Router.Handle(path, handler)
	}

	fmt.Println("Server is running on port", s.WebServerPort)

	err := http.ListenAndServe(":"+s.WebServerPort, s.Router)
	if err != nil {
		panic(err)
	}
}
