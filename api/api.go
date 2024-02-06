package api

import "github.com/gorilla/mux"

type Routes struct {
	Root    *mux.Router
	APIRoot *mux.Router

	Forecast *mux.Router
	Crag     *mux.Router
	Climbs   *mux.Router
	//Users
	//Comments
	//Posts
	//Auth
}

type API struct {
	srv       *Server
	BaseRoute *Routes
}
