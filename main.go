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

	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	store "github.com/lregs/Crag/SqlStore"
	"github.com/lregs/Crag/server"
	met "github.com/lregs/Crag/services/metoffice"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron"
)

func main() {

	log := NewLogger("log.txt")
	log.Println("log created")

	conn, err := pgxpool.Connect(context.Background(), "postgres://postgres:postgres@CragDb:5432/postgres")
	if err != nil {
		log.Printf("db connection failed %s", err)
	}

	if err := initdb(log, conn); err != nil {
		panic(err)
	}

	store, err := store.NewSqlStore(&store.StoreConfig{DbConnection: conn})
	if err != nil {
		log.Fatalf("Could not create store because of error: %s", err)
	}

	rc := redis.NewClient(&redis.Options{
		Addr:     "redis-19441.c233.eu-west-1-1.ec2.redns.redis-cloud.com:19441",
		Password: "N9jHgekt2GxfqkHpQtNHL7jmwUCkq3zA",
		DB:       0,
	})

	client := http.Client{}

	c := cron.New()
	c.AddFunc("0 0 * * *", func() {
		log.Println("cron started")

		//we dont want to be making uneccassary requests. At the moment we're getting the data and then trying to store it
		//and checking if it needs to be stored. This should be happening before the data is even requested.

		data, err := met.GetForecast(client, []float64{53.121482791166194, -3.9988571454802284})
		if err != nil {
			log.Printf("couldn't get forecast %s", err)
		}
		payload, err := met.GetPayload(log, data)
		if err != nil {
			log.Printf("error getting payload %s", err)
		}
		if err := met.StoreData(log, context.Background(), rc, payload); err != nil {
			log.Printf("failed storing during cron: %s", err)
		}
	})

	c.Start()
	log.Println("cron Started")

	srv := server.NewServer(context.Background(), log, store)

	log.Println("init redis")
	if err = initRedis(log, client, rc); err != nil {
		log.Printf("error init db %s", err)
	}

	log.Println("starting server")
	err = http.ListenAndServe(":6969", srv)
	if err != nil {
		log.Fatalf("could not start srv because of err: %s", err)
	}

}

func initRedis(log *log.Logger, client http.Client, rdb *redis.Client) error {

	exists, err := rdb.Exists(context.Background(), "LastUpdated").Result()
	if err != nil {
		log.Printf("failed checking key for last update %s", err)
		return err
	}

	if exists != 0 {
		lastUpdate, err := rdb.Get(context.Background(), "LastUpdated").Result()
	}

	//check if data needs updating:

	data, err := met.GetForecast(client, []float64{53.121482791166194, -3.9988571454802284})
	if err != nil {
		log.Printf("couldn't get forecast %s", err)
	}
	payload, err := met.GetPayload(log, data)
	if err != nil {
		log.Printf("error getting payload %s", err)
	}
	err = met.StoreData(log, context.Background(), rc, payload)
	if err != nil {
		log.Printf("error storing %s", err)
	}
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
