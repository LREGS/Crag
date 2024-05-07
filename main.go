// package main

// import (
// 	"fmt"
// 	"net/http"
// )

// type Server struct {
// 	Store *InMemoryCragStore
// }

// func main() {
// 	store := NewInMemoryCragStore()
// 	server := NewServer(store)

// 	err := http.ListenAndServe(":6969", server)
// 	if err != nil {
// 		fmt.Println("Error starting server")
// 	}
// }

package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/lregs/Crag/models"
	met "github.com/lregs/Crag/services/metoffice"
)

func main() {
	// Connect to the database
	// db, err := sql.Open("postgres", "host=CragDb user=postgres password=postgres dbname=postgres sslmode=disable")
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()

	// store, err := store.NewSqlStore(&store.StoreConfig{DbConnection: db})
	// if err != nil {
	// 	log.Fatalf("Could not create store because of error: %s", err)
	// }

	// srv := server.NewServer(store)

	// err = http.ListenAndServe(":6969", srv)
	// if err != nil {
	// 	log.Fatalf("could not start srv because of err: %s", err)
	// }

	returnD, err, _ := met.GetForecast([]float64{53.12000233374393, -4.000659549362343})
	if err != nil {
		fmt.Print(err)
	}

	// jsons, err := json.Marshal(returnD.Features[0].Properties.TimeSeries)
	// if err != nil {
	// 	fmt.Printf("error encoding %s", err)
	// }

	// txt, err := os.Create("txt.txt")
	// if err != nil {
	// 	return
	// }

	file, err := os.Create("forecast.csv")
	if err != nil {
		return
	}

	defer file.Close()
	// _, _ = file.Write(jsons)
	// _, _ = txt.Write(body)
	// fmt.Print(db.Close())

	f2slice := f2csv(returnD)

	w := csv.NewWriter(file)
	w.WriteAll(f2slice)

}

func f2csv(f models.Forecast) [][]string {

	// fmt.Println(f)

	d := f.Features[0].Properties.TimeSeries

	result := make([][]string, len(d))

	//header
	result[0] = []string{"Id", "Time", "ScreenTemperature", "FeelsLikeTemp", "WindSpeed",
		"WindDirection", "totalPrecipitation", "ProbofPrecipitation", "Latitude", "Longitude"}

	for i := 1; i < len(d); i++ {
		result[i] = []string{
			strconv.Itoa(i),
			d[i].Time,
			strconv.FormatFloat(d[i].ScreenTemperature, 'f', -1, 64),
			strconv.FormatFloat(d[i].FeelsLikeTemperature, 'f', -1, 64),
			strconv.FormatFloat(d[i].WindSpeed10m, 'f', -1, 64),
			strconv.Itoa(d[i].WindDirectionFrom10m),
			strconv.FormatFloat(d[i].TotalPrecipAmount, 'f', -1, 64),
			strconv.Itoa(d[i].ProbOfPrecipitation),
			strconv.FormatFloat(f.Features[0].Geometry.Coordinates[0], 'f', -1, 64),
			strconv.FormatFloat(f.Features[0].Geometry.Coordinates[1], 'f', -1, 64),
		}
	}

	return result
}
