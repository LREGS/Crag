//store is initialised inside channels/app/server.go in mattermost
//The server has values services, which a store is part of a service

package dbStore

import "workspaces/github.com/lregs/Crag/data"

type ForecastStore interface {
	StoreForecast(forecast data.Forecast) error
	GetForecastByID(Id int) (forecast data.DBForecast, err error)
	GetMultipleForecastByID([]int) (forecasts map[int][]interface{}, err error)
	UpdateForecast(forecast data.DBForecast) error
	DryestForecast() (err error, DryestCrags []data.DBForecast)
	MarshalForecastToDB(forecast data.Forecast) (data.DBForecast, error)
}

//store needs init 	store.stores.user = newSqlUserStore(store, metrics)