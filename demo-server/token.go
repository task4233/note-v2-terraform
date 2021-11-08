package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type tokenResource struct{}

func (rs tokenResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", rs.Create)

	return r
}

func (rs tokenResource) Create(w http.ResponseWriter, r *http.Request) {

}
