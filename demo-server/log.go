package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/go-chi/chi/v5"

	"github.com/task4233/note-v2-terraform/client"
)

type logResource struct {
}

var l = Log{}

type Log struct {
	logs []*client.Log
	mu   sync.Mutex
}

func (rs logResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", l.Get)
	r.Post("/", l.Post)
	r.Put("/", l.Put)
	r.Delete("/", l.Delete)

	l.mu = sync.Mutex{}

	return r
}

func (rs *Log) Get(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(os.Stderr, "[Read] begins!\n")
	// ひとまず、今の実装ではユーザを管理しないのでorderIDを使う実装は省略する
	/*orderID, err := strconv.Atoi(chi.URLParam(r, "orderID"))
	if err != nil {
		log.Printf("failed in Get: %s\n", err.Error())
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// index out of range
	if orderID > len(l.logs)-1 {
		log.Println("invalid index in Get")
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	*/

	logs := make([]client.OrderItem, len(l.logs))
	for idx := range l.logs {
		logs[idx] = client.OrderItem{
			Log: client.Log{
				Body: l.logs[idx].Body,
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

func (rs *Log) Post(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(os.Stderr, "[Create] begins!\n")
	var log client.Order
	if err := json.NewDecoder(r.Body).Decode(&log); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	l.mu.Lock()
	l.logs = []*client.Log{}
	for idx := range log.Items {
		l.logs = append(l.logs, &client.Log{
			Body: log.Items[idx].Log.Body},
		)
	}
	l.mu.Unlock()

	logs := make([]client.OrderItem, len(l.logs))
	for idx := range l.logs {
		logs[idx] = client.OrderItem{
			Log: client.Log{
				Body: l.logs[idx].Body,
			},
		}
	}

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

func (rs *Log) Put(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(os.Stderr, "[Update] begins!\n")
	// orderID, err := strconv.Atoi(chi.URLParam(r, "orderID"))
	// if err != nil {
	// 	log.Printf("failed in Update: %s\n", err.Error())
	// 	http.Error(w, "bad request", http.StatusBadRequest)
	// 	return
	// }

	if len(l.logs) == 0 {
		w.WriteHeader(http.StatusOK)
		return
	}
	// index out of range
	// if orderID > len(l.logs)-1 {
	// 	log.Println("invalid index in Update")
	// 	http.Error(w, "bad request", http.StatusBadRequest)
	// 	return
	// }

	var log client.Order
	if err := json.NewDecoder(r.Body).Decode(&log); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	l.mu.Lock()
	l.logs = []*client.Log{}
	for idx := range log.Items {
		l.logs = append(l.logs, &client.Log{
			Body: log.Items[idx].Log.Body},
		)
	}
	l.mu.Unlock()

	logs := make([]client.OrderItem, len(l.logs))
	for idx := range l.logs {
		logs[idx] = client.OrderItem{
			Log: client.Log{
				Body: l.logs[idx].Body,
			},
		}
	}

	result := client.Order{
		Items: logs,
	}

	resp, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (rs *Log) Delete(w http.ResponseWriter, r *http.Request) {
	if len(l.logs) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	l.mu.Lock()
	l.logs = l.logs[:len(l.logs)-1]

	logs := make([]client.OrderItem, len(l.logs))
	for idx := range l.logs {
		logs[idx] = client.OrderItem{
			Log: client.Log{
				Body: l.logs[idx].Body,
			},
		}
	}
	l.mu.Unlock()

	result := client.Order{
		Items: logs,
	}

	resp, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}
