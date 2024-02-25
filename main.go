package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type CragStore interface {
	GetForecast(name string) string
}

type CragServer struct {
	store CragStore
}

type InMemoryCragStore struct{}

func (i *InMemoryCragStore) GetForecast(crag string) string {
	return "dry"
}

func (c CragServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	crag := strings.TrimPrefix(r.URL.Path, "/crags/")
	fmt.Fprint(w, c.store.GetForecast(crag))

}

func GetForecast(crag string) string {
	if crag == "stanage" {
		return "cold"
	}

	if crag == "milestone" {
		return "dry"
	}
	return ""

}

func main() {

	server := &CragServer{&InMemoryCragStore{}}

	log.Fatal(http.ListenAndServe(":6969", server))
}
