package main

import (
	"context"
	"log"
	"net/http"

	"github.com/a-h/templ"
)

// I dont know if in the future we just want to pass an interface of like forecast service
// to my server, that way if the logic for how we're storing or getting the forecasts changes,
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

func NewServer(
	log *log.Logger,
	forecastStore *MetStore,
) http.Handler {
	mux := http.NewServeMux()
	addRoutes(log, mux, forecastStore)
	return mux
}

func addRoutes(log *log.Logger, mux *http.ServeMux, store *MetStore) {
	mux.Handle("/home", HomePage(log, store))
}

func HomePage(log *log.Logger, store *MetStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Print("handling home page request")

		data := make(chan ForecastTotals)
		go GetStoreData(data, r.Context(), store)

		comp := Page(data)
		templ.Handler(comp, templ.WithStreaming()).ServeHTTP(w, r)

	},
	)
}

// Can/should this be in some kind of service struct? - I know im mentioning this above
// but ye maybe be better if inside some kind of service struct ot encapsulate
// the functionality we want for our handlers inside the service struct,
// and the servive struct worries about the store, api etc, and just serving that to our
// handler rather than our handlers having use/access of the store itself.
func GetStoreData(data chan<- ForecastTotals, ctx context.Context, store *MetStore) {
	defer close(data)

	select {
	case <-ctx.Done():
		return
	default:
		d, err := store.Get("orme")
		if err != nil {
			return
		}
		f := d.Totals["18"]
		data <- *f
	}
}
