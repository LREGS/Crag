package api

import (
	"net/http"
)

func (api *API) InitForecast() {
	//at some point when we implement aut, this should be like
	//.handle("" APISessionTokenRequired(postForecast)) - this will mean only those with
	//a valid auth token will be able to post data to the server
	//because at the moment anyone could delete my whole database
	api.BaseRoutes.Forecast.Handle("", getForecasts).Methods("GET") //return all forecasts
	api.BaseRoutes.Forecast.Handle("", postForecast).Methods("POST")
	api.BaseRoutes.Forecast.Handle("/id", getForecastByID).Methods("GET")
	api.BaseRoutes.Forecast.Handle("/id", updateForecastByID).Methods("PUT")
	api.BaseRoutes.Forecast.Handle("/id", deleteForecastByID).Methods("DELETE")
	api.BaseRoutes.Forecast.Handle("/id/date", getForecastIDByDate).Methods("GET")
	api.BaseRoutes.Forecast.Handle("/id/date", deleteForcastIDByDate).Methods("DELETE")

}

func getForecasts(w http.ResponseWriter, r *http.Request) {
	return
}
