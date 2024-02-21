package Services

import "workspaces/github.com/lregs/Crag/Store"

type Services struct {
	Forecast *ForecastService
}

func InitServices(store *Store.SqlStore) (*Services, error) {

	Forecast := NewForecastService(ServiceConfig{ForecastStore: store.Forecast()})

	return &Services{
		Forecast: Forecast,
	}, nil
}
