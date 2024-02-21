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
	"log"
	Services "workspaces/github.com/lregs/Crag/Services"
	"workspaces/github.com/lregs/Crag/api"
	"workspaces/github.com/lregs/Crag/app"
)

func main() {
	server := app.NewServer()

	Services, err := Services.InitServices(server.Store)
	if err != nil {
		log.Fatal(err)
	}

	//should probably have an init function like everything else - gotta stay dry bro
	var Deps api.Dependecnies
	Deps.Services = Services
	Deps.Store = server.Store

	api, err := api.Init(server, &Deps)
	if err != nil {
		log.Fatal(err)
	}

	server.Router.Handle("/api", api.BaseRoutes.Root)
}
