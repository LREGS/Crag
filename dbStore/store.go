//store is initialised inside channels/app/server.go in mattermost
//The server has values services, which a store is part of a service

package dbStore

type ForecastStore interface {
	StoreForecast()
	GetForecast()
	UpdateForecast()
	DeleteForecast()
	GetForecastByDate()
	GetForecastByDryest()
	GetOldestForecast()

}

//store needs init 	store.stores.user = newSqlUserStore(store, metrics)
