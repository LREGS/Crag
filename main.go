package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type CragStore interface {
	GetForecast(name string) string
	addForecast(name string)
}

type CragServer struct {
	store CragStore
}

type InMemoryCragStore struct{}

func (i *InMemoryCragStore) GetForecast(crag string) string {
	return "dry"
}

func (i *InMemoryCragStore) addForecast(crag string) {
}

func (c CragServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		c.processForecast(w, r)
	case http.MethodGet:
		c.showForecast(w, r)

	}

	if r.Method == http.MethodPost {
		w.WriteHeader(http.StatusAccepted)
	}

}

func (c *CragServer) showForecast(w http.ResponseWriter, r *http.Request) {
	crag := strings.TrimPrefix(r.URL.Path, "/crags/")
	forecast := c.store.GetForecast(crag)

	if forecast == "" {
		w.WriteHeader(http.StatusNotFound)
	}
	fmt.Fprint(w, forecast)

}

func (c *CragServer) processForecast(w http.ResponseWriter, r *http.Request) {
	reportedForecast := c.getForecastFromURL(r.URL.Path)
	c.store.addForecast(reportedForecast)
	w.WriteHeader(http.StatusAccepted)
}

func (c *CragServer) getForecastFromURL(url string) string {
	crag := strings.TrimPrefix(url, "/crags/")
	comps := strings.Split(crag, "/")
	return comps[2]

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
