package api

import (
	"net/http"
)

func (api *API) initForecast() {
	api.BaseRoutes.Forecast.Handle("", getForecast).Methods("GET")
}

func getForecast(w http.ResponseWriter, r *http.Request)
