package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// I dont know if in the future we just want to pass an interface of like forecast service
// to my Server, that way if the logic for how we're storing or getting the forecasts changes,
// we will just be calling the service struct, which is agnostic to how its getting or recieving the data
// its just orchestrating the work.
// but maybe this is baking in too much complexity from the start.

// type forecastService interface {
// 	GetForecastTotals() // Get from api
// 	GetData(key string) // Get From Store
// }

// type metAPIService interface {
// 	GetForecast(url string) (Forecast, error)
// }

// type metForecastStore interface {
// 	GetData(key string) (ForecastPayload, error)
// }

// type MetForecastService struct {
// 	Store  metForecastStore
// 	MetAPI metAPIService
// }

type Server struct {
	log   *log.Logger
	mux   *http.ServeMux
	store *MetStore
}

func NewServer(
	log *log.Logger,
	forecastStore *MetStore,
) *Server {
	srv := &Server{log: log, mux: http.NewServeMux(), store: forecastStore}
	srv.Routes()
	return srv
}

func (s *Server) Routes() {

	s.mux.HandleFunc("/all", s.AllForecasts)
}

func (s *Server) AllForecasts(w http.ResponseWriter, r *http.Request) {

	s.log.Println("hit")

	d, err := s.store.GetAll()
	if err != nil {
		log.Println("Failed to get all " + err.Error())
	}

	if err := json.NewEncoder(w).Encode(d); err != nil {
		log.Println("Failed to decode " + err.Error())
	}
}
