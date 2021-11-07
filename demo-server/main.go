package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	PORT = 19090
)

type Server struct {
	handler http.Handler
	logs    []*Log
	mu      sync.Mutex
}

func (s *Server) init() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", s.Get)
	r.Post("/", s.Post)
	r.Put("/", s.Put)
	r.Delete("/", s.Delete)

	s.handler = r
	s.mu = sync.Mutex{}
}

type Log struct {
	Body string `json:"body"`
}

func (s *Server) Get(w http.ResponseWriter, r *http.Request) {
	resp, err := json.Marshal(s.logs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (s *Server) Post(w http.ResponseWriter, r *http.Request) {
	var log Log
	if err := json.NewDecoder(r.Body).Decode(&log); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	s.logs = append(s.logs, &log)
	s.mu.Unlock()

	resp, err := json.Marshal(s.logs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (s *Server) Put(w http.ResponseWriter, r *http.Request) {
	if len(s.logs) == 0 {
		w.WriteHeader(http.StatusOK)
		return
	}

	var log Log
	if err := json.NewDecoder(r.Body).Decode(&log); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	s.logs[0] = &log
	s.mu.Unlock()

	resp, err := json.Marshal(s.logs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (s *Server) Delete(w http.ResponseWriter, r *http.Request) {
	if len(s.logs) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	s.mu.Lock()
	s.logs = s.logs[:len(s.logs)-1]
	s.mu.Unlock()

	resp, err := json.Marshal(s.logs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func main() {
	s := &Server{}
	s.init()
	log.Fatal(http.ListenAndServe(":19090", s.handler))
}
