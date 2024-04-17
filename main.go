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
	"fmt"

	_ "github.com/lib/pq"
	met "github.com/lregs/Crag/services/metoffice"
)

func main() {
	// // Connect to the database
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

	d, err := met.GetForecast([]float64{53.12000233374393, -4.000659549362343})
	if err != nil {
		fmt.Print(err)
	}

	fmt.Print(d.Features)

}
