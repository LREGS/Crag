package main

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type cragStore interface {
	addCrag(crag string)
}

type Server struct {
	Store *InMemoryCragStore
}

func NewServer(store cragStore) http.Handler {
	mux := mux.NewRouter()
	addRoutes(mux, store)

	var handler http.Handler = mux

	//handler = middleware(handler)

	return handler

}

func addRoutes(mux *mux.Router, store cragStore) {
	mux.Handle("https://localhost:6969/crags/", http.HandlerFunc(handlePostCrags(store)))

}

func handlePostCrags(store cragStore) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		crag := strings.TrimPrefix(r.URL.Path, "/crags/")
		store.addCrag(crag)
		w.WriteHeader(http.StatusAccepted)
	}
}
