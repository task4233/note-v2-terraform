package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/task4233/note-v2-terraform/client"
)

const (
	PORT = 19090
)

type Server struct {
	handler http.Handler
	logs    []*client.Log
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

func (s *Server) Get(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(os.Stderr, "[Read] begins!\n")
	// ひとまず、今の実装ではユーザを管理しないのでorderIDを使う実装は省略する
	/*orderID, err := strconv.Atoi(chi.URLParam(r, "orderID"))
	if err != nil {
		log.Printf("failed in Get: %s\n", err.Error())
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// index out of range
	if orderID > len(s.logs)-1 {
		log.Println("invalid index in Get")
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	*/

	logs := make([]client.OrderItem, len(s.logs))
	for idx := range s.logs {
		logs[idx] = client.OrderItem{
			Log: client.Log{
				Body: s.logs[idx].Body,
			},
		}
	}

	result := client.Order{
		Items: logs,
	}

	resp, err := json.Marshal(result)
	if err != nil {
		log.Printf("failed in Get: %s\n", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(os.Stderr, "[Read] resp: %s\n", string(resp))

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (s *Server) Post(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(os.Stderr, "[Create] begins!\n")
	var log client.Order
	if err := json.NewDecoder(r.Body).Decode(&log); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	for idx := range log.Items {
		s.logs = append(s.logs, &log.Items[idx].Log)
	}

	logs := make([]client.OrderItem, len(s.logs))
	for idx := range s.logs {
		logs[idx] = client.OrderItem{
			Log: client.Log{
				Body: s.logs[idx].Body,
			},
		}
	}
	s.mu.Unlock()

	result := client.Order{
		Items: logs,
	}

	resp, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(os.Stderr, "[Create] resp: %s\n", string(resp))

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (s *Server) Put(w http.ResponseWriter, r *http.Request) {
	if len(s.logs) == 0 {
		w.WriteHeader(http.StatusOK)
		return
	}

	var log client.Log
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
	log.Printf("Server running in http://localhost:%d/", PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", PORT), s.handler))
}
