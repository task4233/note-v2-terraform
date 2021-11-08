package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

const (
	PORT = 19090
)

type Server struct {
	handler http.Handler
}

func (s *Server) init() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Mount("/token", tokenResource{}.Routes())
	r.Mount("/logs", logResource{}.Routes())

	s.handler = r
}

func main() {
	s := &Server{}
	s.init()
	log.Printf("Server running in http://localhost:%d/", PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", PORT), s.handler))
}
