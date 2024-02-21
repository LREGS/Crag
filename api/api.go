package api

import (
	"workspaces/github.com/lregs/Crag/Services"
	"workspaces/github.com/lregs/Crag/Store"
	"workspaces/github.com/lregs/Crag/app"

	"github.com/gorilla/mux"
)

type Routes struct {
	Root *mux.Router

	Forecast *mux.Router
	Crag     *mux.Router
	Climbs   *mux.Router
	//Users
	//Comments
	//Posts
	//Auth
}

type Dependecnies struct {
	Store    *Store.SqlStore
	Services *Services.Services
	//auth
	//logging
	//config
	//messaging/tasks
	//external api?? - actually getting data from the met api
}

type API struct {
	BaseRoutes *Routes
	Deps       *Dependecnies
	Server     *app.Server
}

func Init(srv *app.Server, deps *Dependecnies) (*API, error) {
	api := &API{
		BaseRoutes: &Routes{},
		Deps:       deps,
	}

	api.BaseRoutes.Root = srv.Router

	api.BaseRoutes.Forecast = api.BaseRoutes.Root.PathPrefix("/forecast").Subrouter()

	api.BaseRoutes.Crag = api.BaseRoutes.Root.PathPrefix("/crag").Subrouter()

	api.BaseRoutes.Climbs = api.BaseRoutes.Root.PathPrefix("/climb").Subrouter()

	api.InitForecast()
	api.InitCrag(api.Deps)
	// api.InitClimbs()

	return api, nil
}
