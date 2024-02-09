package api

import "github.com/gorilla/mux"

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

type API struct {
	srv        *Server
	BaseRoutes *Routes
}

func Init(srv *Server) (*API, error) {
	api := &API{
		srv:        srv,
		BaseRoutes: &Routes{},
	}

	api.BaseRoutes.Root = srv.Router

	api.BaseRoutes.Forecast = api.BaseRoutes.Root.PathPrefix("/forecast").SubRouter()

	api.BaseRoutes.Crag = api.BaseRoutes.Root.PathPrefix("/crag").SubRouter()

	api.BaseRoutes.Climbs = api.BaseRoutes.Root.PathPrefix("/climb").SubRouter()

	api.InitForecast()
	api.InitCrag()
	api.InitClimbs()
}
