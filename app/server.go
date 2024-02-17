package app

import (
	"net/http"
	"workspaces/github.com/lregs/Crag/Store"
	"workspaces/github.com/lregs/Crag/app/forecast"

	"github.com/gorilla/mux"
)

type Server struct {
	RootRouter *mux.Router

	//Router for API
	Router *mux.Router

	Server *http.Server

	//Services for the app - atm forecast service just holds all the services for forecast,
	//mayb this wants to be user, with auth and which parts of the db(forecast being one of them)
	//the user can access but atm we make forecast service for application logic across the
	//back end
	forecastService *forecast.ForecastService
	//user services

	//need to use a logger - whether I implement now I dont know
	// logger *log.Logger

	store *Store.SqlStore
}

func (s *Server) Store() *Store.SqlStore {
	return s.store
}

func NewServer() *Server {

	rootRouter := mux.NewRouter()
	Router := mux.NewRouter()

	HTTPServer := http.Server{Addr: "127.0.0.1:8080"}

	s := &Server{
		RootRouter: rootRouter,
		Router:     Router,
		Server:     &HTTPServer,
	}

	store, sErr := Store.New()
	if sErr != nil {
		return nil
	}
	s.store = store

	//im sure there should be some kind of error handling here
	s.forecastService = forecast.New(forecast.ServiceConfig{
		ForecastStore: s.Store().Forecast(),
		//external api service
	})
	return s

}
