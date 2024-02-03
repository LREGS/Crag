// package main

// import (
// 	"database/sql"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"os"

// 	_ "github.com/lib/pq"

// 	client "workspaces/github.com/lregs/Crag/client"
// 	helpers "workspaces/github.com/lregs/Crag/helper"

// 	"workspaces/github.com/lregs/Crag/utils"

// 	"github.com/go-chi/chi/v5"
// )

// var Db *sql.DB

// func init() {
// 	DbUsername := os.Getenv("DB_USERNAME")
// 	DbPassword := os.Getenv("DB_PASSWORD")

// 	var err error
// 	Db, err = sql.Open("postgres", fmt.Sprintf("user=%s dbname=crag password=%s sslmode=disable", DbUsername, DbPassword))
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func main() {

// 	const port = "8080"

// 	r := chi.NewRouter()

// 	// r.Get("/DryCrag", whichCragDry)
// 	// coords := getCragCoords()
// 	// fmt.Println(coords)

// 	// //WORKING TEST REQUEST
// 	// coords := []float64{53.122664, -3.998611}
// 	client := client.DefaultClient()
// 	// helpers.CheckError(err)
// 	headers := helpers.ReturnHeaders()

// 	forecastsToGet := getCragCoords()

// 	for _, value := range *forecastsToGet {
// 		url, err := helpers.MetOfficeURL(value)
// 		helpers.CheckError(err)
// 		f, err := utils.GetForecast(url, headers, client)
// 		fmt.Println(f.Features[0].Properties.TimeSeries[0])

// 	}
// 	// f, err := utils.GetForecast(url, headers, client)
// 	// fmt.Println(f.Features[0].Properties.TimeSeries)

// 	http.ListenAndServe(":8080", r)

// }

// func whichCragDry() {

// }

// func getCragCoords() *map[string][]float64 {
// 	file, err := os.Open("data/crags.json")
// 	helpers.CheckError(err)
// 	defer file.Close()

// 	fileInfo, err := file.Stat()
// 	helpers.CheckError(err)

// 	size := fileInfo.Size()

// 	content := make([]byte, size)
// 	_, err = file.Read(content)
// 	helpers.CheckError(err)

// 	var coordMap = make(map[string][]float64)

// 	err = json.Unmarshal(content, &coordMap)
// 	helpers.CheckError(err)

// 	// coords := make([][]float64, 0)
// 	// for _, vals := range coordMap {

// 	// 	coords = append(coords, vals)

// 	// }

// 	return &coordMap

// }

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	data "workspaces/github.com/lregs/Crag/data"
	helpers "workspaces/github.com/lregs/Crag/helper"
	m "workspaces/github.com/lregs/Crag/mw"

	_ "github.com/lib/pq"
)

var Db *sql.DB

func init() {
	// DbUsername := os.Getenv("DB_USERNAME")
	// DbPassword := os.Getenv("DB_PASSWORD")

	var err error
	Db, err = sql.Open("postgres", fmt.Sprintf("user=william dbname=crag password=1 sslmode=disable"))
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

func handleGetCrag(w http.ResponseWriter, r *http.Request) (err error) {
	id, err := m.IDFromRequest(r)
	helpers.CheckError(err)

	crag, err := data.GetCrag(id, Db)
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

func handlePutCrag(w http.ResponseWriter, r *http.Request) (err error) {
	id, err := m.IDFromRequest(r)
	helpers.CheckError(err)

	crag, err := data.GetCrag(id, Db)
	helpers.CheckError(err)

	len := r.ContentLength
	body := make([]byte, len)

	r.Body.Read(body)
	json.Unmarshal(body, &crag)

	err = crag.UpdateCrag(Db)
	helpers.CheckError(err)

	w.WriteHeader(200)
	return
}

func handleDeleteCrag(w http.ResponseWriter, r *http.Request) (err error) {
	id, err := m.IDFromRequest(r)
	helpers.CheckError(err)

	crag, err := data.GetCrag(id, Db)
	helpers.CheckError(err)

	err = crag.DeleteCrag(Db)
	helpers.CheckError(err)
	w.WriteHeader(200)
	return

}
