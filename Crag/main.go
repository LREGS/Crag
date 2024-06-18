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
	store "github.com/lregs/CragWeather/Crag/SqlStore"
	"github.com/lregs/CragWeather/Crag/server"
)

func main() {

	log := NewLogger("log.txt")
	log.Println("log created")

	conn, err := pgxpool.Connect(context.Background(), "postgres://postgres:postgres@CragDb:5432/postgres")
	if err != nil {
		log.Printf("db connection failed %s", err)
	}

	store, err := store.NewSqlStore(&store.StoreConfig{DbConnection: conn})
	if err != nil {
		log.Fatalf("Could not create store because of error: %s", err)
	}

	srv := server.NewServer(context.Background(), log, store)

	log.Println("starting server")
	err = http.ListenAndServe(":6969", srv)
	if err != nil {
		log.Fatalf("could not start srv because of err: %s", err)
	}

}

func NewLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("bad file")
	}
	return log.New(logfile, "[main]", log.Ldate|log.Ltime|log.Lshortfile)
}
