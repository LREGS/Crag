package api

import (
	h "workspaces/github.com/lregs/Crag/handler"

	"github.com/gorilla/mux"
)

func main() {
	mainRouter := mux.NewRouter()

	CragRouter := mainRouter.PathPrefix("/crag").SubRouter()
	ForecastRouter := mainRouter.PathPrefix("/forecast").SubRouter()

	CragRouter.HandleFunc("/", h.HandleCragRequests).Methods("GET", "PUT", "POST", "DELETE")

	ForecastRouter.HandleFunc("/", h.HandleForecastRequests).Method("GET", "PUT", "POST", "DELETE")

}
