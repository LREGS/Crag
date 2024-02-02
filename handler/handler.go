package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"path"
	"strconv"
	data "workspaces/github.com/lregs/Crag/data"
	helpers "workspaces/github.com/lregs/Crag/helpers"

	_ "github.com/lib/pq"
)

func handleCragRequests(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var err error
	switch r.Method {
	case "GET":
		err = handleGetCrag(w, r, db)
	}
}

func handleGetCrag(w http.ResponseWriter, r *http.Request, db *sql.DB) (err error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	helpers.CheckError(err)

	crag, err := data.GetCrag(id, db)
	helpers.CheckError(err)

	output, err := json.MarshalIndent(&crag, "", "\t\t")
	helpers.CheckError(err)

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return

}
