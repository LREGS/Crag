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
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	store "github.com/lregs/Crag/SqlStore"
	"github.com/lregs/Crag/server"
)

func main() {
	// Connect to the database
	db, err := sql.Open("postgres", "host=CragDb user=postgres password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	log := NewLogger("log.txt")
	log.Println("log created")

	if err := initdb(log, db); err != nil {
		log.Panicf("Creating db failed %s", err)
	}

	store, err := store.NewSqlStore(&store.StoreConfig{DbConnection: db})
	if err != nil {
		log.Fatalf("Could not create store because of error: %s", err)
	}

	// store.Stores.ForecastStore.Populate(log)

	// time.NewTimer(20 * time.Second)

	// store.Stores.ForecastStore.Refresh(log)

	srv := server.NewServer(log, store)

	err = http.ListenAndServe(":6969", srv)
	if err != nil {
		log.Fatalf("could not start srv because of err: %s", err)
	}

	// returnD, err, _ := met.GetForecast([]float64{53.12000233374393, -4.000659549362343})
	// if err != nil {
	// 	fmt.Print(err)
	// }

	// jsons, err := json.Marshal(returnD.Features[0].Properties.TimeSeries)
	// if err != nil {
	// 	fmt.Printf("error encoding %s", err)
	// }

	// txt, err := os.Create("txt.txt")
	// if err != nil {
	// 	return
	// }

	// file, err := os.Create("forecast.csv")
	// if err != nil {
	// 	return
	// }

	// defer file.Close()
	// // _, _ = file.Write(jsons)
	// _, _ = txt.Write(body)
	// fmt.Print(db.Close())

	// if err := gocron.Every(1).Day().At("03:30").Do(store.Stores.ForecastStore.Refresh()); err != nil {
	// 	log.Printf("Refresh failed %s", err)
	// }

	// f2slice := f2csv(returnD)

	// w := csv.NewWriter(file)
	// w.WriteAll(f2slice)

}

// func f2csv(f models.Forecast) [][]string {

// 	// fmt.Println(f)

// 	d := f.Features[0].Properties.TimeSeries

// 	result := make([][]string, len(d))

// 	//header
// 	result[0] = []string{"Id", "Time", "ScreenTemperature", "FeelsLikeTemp", "WindSpeed",
// 		"WindDirection", "totalPrecipitation", "ProbofPrecipitation", "Latitude", "Longitude"}

// 	for i := 1; i < len(d); i++ {
// 		result[i] = []string{
// 			strconv.FormatFloat(d[i].FeelsLikeTemperature, 'f', -1, 64),
// 			strconv.FormatFloat(d[i].WindSpeed10m, 'f', -1, 64),
// 			strconv.Itoa(d[i].WindDirectionFrom10m),
// 			strconv.FormatFloat(d[i].TotalPrecipAmount, 'f', -1, 64),
// 			strconv.Itoa(d[i].ProbOfPrecipitation),
// 			strconv.FormatFloat(f.Features[0].Geometry.Coordinates[0], 'f', -1, 64),
// 			strconv.FormatFloat(f.Features[0].Geometry.Coordinates[1], 'f', -1, 64),
// 		}
// 	}

// 	return result
// }

func NewLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("bad file")
	}
	return log.New(logfile, "[main]", log.Ldate|log.Ltime|log.Lshortfile)
}

var initDb = `DROP TABLE IF EXISTS forecast;
DROP TABLE IF EXISTS report;
DROP TABLE IF EXISTS climb;
DROP TABLE IF EXISTS crag;

-- Create tables
CREATE TABLE crag (
	Id SERIAL PRIMARY KEY, 
	Name TEXT UNIQUE, 
	Latitude DOUBLE PRECISION,
	Longitude DOUBLE PRECISION
);

CREATE TABLE climb (
	Id SERIAL PRIMARY KEY,
	Name VARCHAR(255) UNIQUE,
	Grade VARCHAR(255),
	CragID INTEGER REFERENCES crag(Id)
);

CREATE TABLE report (
	Id SERIAL PRIMARY KEY, 
	Content VARCHAR(255),
	Author VARCHAR(255),
	CragID INTEGER REFERENCES crag(Id)
);

CREATE TABLE forecast (
	Id SERIAL PRIMARY KEY, 
	Time VARCHAR(255) UNIQUE,
	ScreenTemperature DOUBLE PRECISION,
	FeelsLikeTemp DOUBLE PRECISION, 
	WindSpeed DOUBLE PRECISION,
	WindDirection DOUBLE PRECISION,
	totalPrecipitation DOUBLE PRECISION,
	ProbofPrecipitation INT,
	Latitude DOUBLE PRECISION,
	Longitude DOUBLE PRECISION
);`

func initdb(log *log.Logger, db *sql.DB) error {
	if _, err := db.Exec(initDb); err != nil {
		log.Printf("init db failed %s", err)
	}

	log.Printf("db inited")
	return nil

}
