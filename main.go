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
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	store "github.com/lregs/Crag/SqlStore"
	"github.com/lregs/Crag/server"
)

func main() {

	conn, err := pgxpool.Connect(context.Background(), "postgres://postgres:postgres@CragDb:5432/postgres")
	if err != nil {
		panic(err)
	}

	// Connect to the database
	// db, err := sql.Open("postgres", "host=CragDb user=postgres password=postgres dbname=postgres sslmode=disable")
	// if err != nil {
	// 	panic(err)
	// }

	// defer db.Close()

	log := NewLogger("log.txt")
	log.Println("log created")

	if err := initdb(log, conn); err != nil {
		panic(err)
	}

	store, err := store.NewSqlStore(&store.StoreConfig{DbConnection: conn})
	if err != nil {
		log.Fatalf("Could not create store because of error: %s", err)
	}

	time.NewTimer(20 * time.Second)

	store.Stores.ForecastStore.Populate(context.Background(), log)

	store.Stores.ForecastStore.Refresh(context.Background(), log)

	srv := server.NewServer(context.Background(), log, store)

	err = http.ListenAndServe(":6969", srv)
	if err != nil {
		log.Fatalf("could not start srv because of err: %s", err)
	}

	// if err := gocron.Every(1).Day().At("03:30").Do(store.Stores.ForecastStore.Refresh()); err != nil {
	// 	log.Printf("Refresh failed %s", err)
	// }

}

func NewLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("bad file")
	}
	return log.New(logfile, "[main]", log.Ldate|log.Ltime|log.Lshortfile)
}

// I just need to edit the posgres docker image to include the execution of this on start-up or init - maybe

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
	Id Int, 
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

func initdb(log *log.Logger, db *pgxpool.Pool) error {
	if _, err := db.Exec(context.Background(), initDb); err != nil {
		log.Printf("init db failed %s", err)
	}

	log.Printf("db inited")
	return nil

}
