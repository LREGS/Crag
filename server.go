package main

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Server struct {
	Store *InMemoryCragStore
}

func NewServer(cragStore *InMemoryCragStore) http.Handler {
	mux := mux.NewRouter()
	addRoutes(mux, cragStore)

	var handler http.Handler = mux

	//handler = middleware(handler)

	return handler

}

func addRoutes(mux *mux.Router, cragStore *InMemoryCragStore) {
	mux.Handle("/crags", http.HandlerFunc(handlePostCrags(cragStore)))

}

func handlePostCrags(s *InMemoryCragStore) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		crag := strings.TrimPrefix(r.URL.Path, "/crags/")
		s.addCrag(crag)
		w.WriteHeader(http.StatusAccepted)
	}
}
