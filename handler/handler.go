package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"
	data "workspaces/github.com/lregs/Crag/data"
	helpers "workspaces/github.com/lregs/Crag/helper"

	_ "github.com/lib/pq"
)

var Db *sql.DB

func init() {
	DbUsername := os.Getenv("DB_USERNAME")
	DbPassword := os.Getenv("DB_PASSWORD")
	intPassword, errs := strconv.Atoi(DbPassword)
	helpers.CheckError(errs)

	var err error
	Db, err = sql.Open("postgres", fmt.Sprintf("user=%s dbname=crag password=%d sslmode=disable", DbUsername, intPassword))
	if err != nil {
		panic(err)
	}
}

func main() {
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}
	http.HandleFunc("/crag/", handleCragRequests)
	server.ListenAndServe()
}

func handleCragRequests(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.Method {
	case "GET":
		err = handleGetCrag(w, r, Db)
	case "POST":
		err = handlePostCrag(w, r)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	helpers.CheckError(err)
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

func handlePostCrag(w http.ResponseWriter, r *http.Request) (err error) {
	len := r.ContentLength
	body := make([]byte, len)

	r.Body.Read(body)

	var crag data.Crag

	json.Unmarshal(body, &crag)
	err = crag.Create(Db)
	helpers.CheckError(err)

	w.WriteHeader(200)
	return
}
