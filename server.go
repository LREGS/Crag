package main

import (
	"fmt"
	"net/http"
	"strings"
)

type CragStore interface {
	GetForecast(name string) string
	addForecast(name, forecast string)
}

type CragServer struct {
	store CragStore
}

func (c CragServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	crag := strings.TrimPrefix(r.URL.Path, "/crags/")
	switch r.Method {
	case http.MethodPost:
		c.processForecast(w, r, crag)
	case http.MethodGet:
		c.showForecast(w, r, crag)

	}

	if r.Method == http.MethodPost {
		w.WriteHeader(http.StatusAccepted)
	}

}

func (c *CragServer) showForecast(w http.ResponseWriter, r *http.Request, crag string) {

	forecast := c.store.GetForecast(crag)

	if forecast == "" {
		w.WriteHeader(http.StatusNotFound)
	}
	fmt.Fprint(w, forecast)

}

func (c *CragServer) processForecast(w http.ResponseWriter, r *http.Request, crag string) {
	reportedForecast := c.getForecastFromURL(r.URL.Path)
	c.store.addForecast(crag, reportedForecast)
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
