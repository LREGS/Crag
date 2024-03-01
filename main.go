package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Server struct {
	Store *InMemoryCragStore
}

func NewServer(store cragStore) http.Handler {
	mux := mux.NewRouter()
	addRoutes(mux, store)

	// var handler http.Handler = mux

	//handler = middleware(handler)

	return http.Handler(mux)

}

func addRoutes(mux *mux.Router, store cragStore) {
	mux.PathPrefix("/crags/").Handler(http.HandlerFunc(handlePostCrags(store))).Methods("POST")
	mux.HandleFunc("/", handleRoot())

}

func handlePostCrags(store cragStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		crag := strings.Split(strings.TrimSuffix(r.URL.Path, "/crags/"), "/")
		store.addCrag(crag[2])
		w.WriteHeader(http.StatusAccepted)
	}
}

func handleRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Welcome to the root path!")
	}

}

func NewInMemoryCragStore() *InMemoryCragStore {
	return &InMemoryCragStore{[]string{}}
}

type InMemoryCragStore struct {
	Names []string
}

func (i *InMemoryCragStore) addCrag(name string) {
	i.Names = append(i.Names, name)
}

func main() {
	store := NewInMemoryCragStore()
	server := NewServer(store)

	err := http.ListenAndServe(":6969", server)
	if err != nil {
		fmt.Println("Error starting server")
	}
}

type cragStore interface {
	addCrag(crag string)
}
