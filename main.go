package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	Store *InMemoryCragStore
}

func NewServer(store cragStore) http.Handler {
	mux := mux.NewRouter()
	addRoutes(mux, store)

	// var handler http.Handler = mux

	//handler = middleware(handler)

	return http.Handler(mux)

}

func addRoutes(mux *mux.Router, store cragStore) {
	mux.PathPrefix("/crags/{key}").HandlerFunc(handlePostCrag(store)).Methods("POST")
	mux.PathPrefix("/crags/{key}").HandlerFunc(handleGetCrag(store)).Methods("GET")
	mux.HandleFunc("/", handleRoot()).Methods("GET")

}

func handlePostCrag(store cragStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		crag := vars["key"]
		store.addCrag(crag)
		w.WriteHeader(http.StatusAccepted)
	}
}

func handleGetCrag(store cragStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		cragNames := store.getNames()
		for _, name := range cragNames {
			if name == vars["key"] {
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, name)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}
	}
}

func handleRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Welcome to the root path!")
	}

}

func NewInMemoryCragStore() *InMemoryCragStore {
	return &InMemoryCragStore{[]string{}}
}

type InMemoryCragStore struct {
	Names []string
}

func (i *InMemoryCragStore) addCrag(name string) {
	i.Names = append(i.Names, name)
}

func (i *InMemoryCragStore) getNames() []string {
	return i.Names
}

func main() {
	store := NewInMemoryCragStore()
	server := NewServer(store)

	err := http.ListenAndServe(":6969", server)
	if err != nil {
		fmt.Println("Error starting server")
	}
}

type cragStore interface {
	addCrag(crag string)
	getNames() []string
}

// package main

// import (
// 	"database/sql"
// 	"fmt"

// 	_ "github.com/lib/pq"
// )

// func main() {
// 	// Connect to the database
// 	db, err := sql.Open("postgres", "host=CragDb user=postgres password=postgres dbname=postgres sslmode=disable")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer db.Close()

// 	// Ping the database to check the connection
// 	err = db.Ping()
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println("Connected to the PostgreSQL database!")

// 	// Perform a simple query
// 	rows, err := db.Query("SELECT version()")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer rows.Close()

// 	// Iterate over the rows
// 	for rows.Next() {
// 		var version string
// 		if err := rows.Scan(&version); err != nil {
// 			panic(err)
// 		}
// 		fmt.Println("PostgreSQL version:", version)
// 	}
// 	if err := rows.Err(); err != nil {
// 		panic(err)
// 	}
// }
