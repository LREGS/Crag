package handlers

import (
	"database/sql"
	"net/http"
	m "workspaces/github.com/lregs/Crag/mw"
	helpers "workspaces/github.com/lregs/Crag/helper"



	_ "github.com/lib/pq"
)

 func HandleForecastRequests(w http.ResponseWriter, r *http.Request){
	var err error
	switch r.Method {
	case "GET":
		err = handleGetCrag(w, r)
	case "POST":
		err = handlePostCrag(w, r)
	case "PUT":
		err = handlePutCrag(w, r)
	case "DELETE":
		err = handleDeleteCrag(w, r)

	}
	helpers.CheckError(err)

 }

 func handGetForecast(w http.ResponseWriter, r *http.Request){
	id, err := m.IDFromRequest(r)
	helpers.CheckError(err)

	forecast, err := store.ForecastStore.GetForecast()
 }