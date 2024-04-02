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

	store, err := store.NewSqlStore(&store.StoreConfig{DbConnection: db})
	if err != nil {
		log.Fatalf("Could not create store because of error: %s", err)
	}

	srv := server.NewServer(store)

	err = http.ListenAndServe(":6969", srv)
	if err != nil {
		log.Fatalf("could not start srv because of err: %s", err)
	}

}
