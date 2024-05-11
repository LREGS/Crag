package server

import (
	"net/http"

	"github.com/gorilla/mux"
	store "github.com/lregs/Crag/SqlStore"
	"github.com/lregs/Crag/services/climb"
	"github.com/lregs/Crag/services/crag"
	"github.com/lregs/Crag/services/forecast"
)

func NewServer(store *store.SqlStore) http.Handler {
	mux := mux.NewRouter()
	subrouter := mux.PathPrefix("/api/v1").Subrouter()
	cragRouter := subrouter.PathPrefix("/crags").Subrouter()
	climbRouter := subrouter.PathPrefix("/climb").Subrouter()
	forecastRouter := subrouter.PathPrefix("/forecast").Subrouter()

	cragHandler := crag.NewHandler(store.Stores.CragStore)
	cragHandler.RegisterRoutes(cragRouter)

	climbHandler := climb.NewHandler(store.Stores.ClimbStore)
	climbHandler.RegisterRoutes(climbRouter)
	// addRoutes(mux, store)

	forecastHanlder := forecast.NewHandler(store.Stores.ForecastStore)
	forecastHanlder.RegisterRoutes(forecastRouter)

	return mux
}
