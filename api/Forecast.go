package api

import "workspaces/github.com/lregs/Crag/web"

func (api *API) InitForecast() {
	//at some point when we implement aut, this should be like
	//.handle("" APISessionTokenRequired(postForecast)) - this will mean only those with
	//a valid auth token will be able to post data to the server
	//because at the moment anyone could delete my whole database
	api.BaseRoutes.Forecast.Handle("", web.APIHandler(getForecasts)).Methods("GET") //return all forecasts
	api.BaseRoutes.Forecast.Handle("", postForecast).Methods("POST")
	api.BaseRoutes.Forecast.Handle("/id", getForecastByID).Methods("GET")
	api.BaseRoutes.Forecast.Handle("/id", updateForecastByID).Methods("PUT")
	api.BaseRoutes.Forecast.Handle("/id", deleteForecastByID).Methods("DELETE")
	api.BaseRoutes.Forecast.Handle("/id/date", getForecastIDByDate).Methods("GET")
	api.BaseRoutes.Forecast.Handle("/id/date", deleteForcastIDByDate).Methods("DELETE")

}

// use this pattern so that we can use DI on the handler and then return a Handler func
// func GetForecast(logger *Logger) http.Handler {
// 	thing := prepareThing()
// 	return http.HandlerFunc(
// 		func(w http.ResponseWriter, r *http.Request) {
// 			// use thing to handle request
// 			logger.Info(r.Context(), "msg", "handleSomething")
// 		}
// 	)
// }
