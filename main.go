package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	client "workspaces/github.com/lregs/Crag/client"
	helpers "workspaces/github.com/lregs/Crag/helper"

	// helpers "workspaces/github.com/lregs/Crag/helper"

	"workspaces/github.com/lregs/Crag/utils"

	"github.com/go-chi/chi/v5"
)

type Coordinates struct {
	Milestone  []float64 `json:"milestone"`
	Gorlan     []float64 `json:"gorlan"`
	Elephant   []float64 `json:"elephant"`
	PortYsgo   []float64 `json:"Port Ysgo"`
	Cratcliffe []float64 `json:"cratcliffe"`
}

func main() {

	const port = "8080"

	r := chi.NewRouter()

	// r.Get("/DryCrag", whichCragDry)
	// coords := getCragCoords()
	// fmt.Println(coords)

	// //WORKING TEST REQUEST
	// coords := []float64{53.122664, -3.998611}
	client := client.DefaultClient()
	// helpers.CheckError(err)
	headers := helpers.ReturnHeaders()

	forecastsToGet := getCragCoords()

	for _, value := range *forecastsToGet {
		url, err := helpers.MetOfficeURL(value)
		helpers.CheckError(err)
		f, err := utils.GetForecast(url, headers, client)
		fmt.Println(f.Features[0].Properties.TimeSeries[0])

	}
	// f, err := utils.GetForecast(url, headers, client)
	// fmt.Println(f.Features[0].Properties.TimeSeries)

	http.ListenAndServe(":8080", r)

}

func whichCragDry() {

}

func getCragCoords() *map[string][]float64 {
	file, err := os.Open("data/crags.json")
	helpers.CheckError(err)
	defer file.Close()

	fileInfo, err := file.Stat()
	helpers.CheckError(err)

	size := fileInfo.Size()

	content := make([]byte, size)
	_, err = file.Read(content)
	helpers.CheckError(err)

	var coordMap = make(map[string][]float64)

	err = json.Unmarshal(content, &coordMap)
	helpers.CheckError(err)

	// coords := make([][]float64, 0)
	// for _, vals := range coordMap {

	// 	coords = append(coords, vals)

	// }

	return &coordMap

}
