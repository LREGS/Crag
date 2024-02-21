package Services

import (
	"workspaces/github.com/lregs/Crag/Store"
)

//provides the service, to the forecast package, that provides the business logic to handle forecasts within our app/requests

type ForecastService struct {
	store Store.ForecastStore
	//api service - interace with met office api
	//config/license
	//Session store etc
}

type ServiceConfig struct {
	ForecastStore Store.ForecastStore
}

func NewForecastService(c ServiceConfig) *ForecastService {
	return &ForecastService{
		store: c.ForecastStore,
	}
}
