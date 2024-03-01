package main

import (
	"fmt"
	"net/http"

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
	mux.PathPrefix("/crags/{key}").HandlerFunc(handlePostCrag(store)).Methods("POST")
	mux.PathPrefix("/crags/{key}").HandlerFunc(handleGetCrag(store)).Methods("GET")
	mux.HandleFunc("/", handleRoot()).Methods("GET")

}

func handlePostCrag(store cragStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		crag := vars["key"]
		store.addCrag(crag)
		w.WriteHeader(http.StatusAccepted)
	}
}

func handleGetCrag(store cragStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		cragNames := store.getNames()
		for _, name := range cragNames {
			if name == vars["key"] {
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, name)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}
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

func (i *InMemoryCragStore) getNames() []string {
	return i.Names
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
	getNames() []string
}
