package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	store "github.com/lregs/Crag/SqlStore"
)

func NewServer(store store.CragStore) http.Handler {
	mux := mux.NewRouter()

	addRoutes(mux, store)

	return mux
}

func addRoutes(mux *mux.Router, cragStore store.CragStore) {

	// mux.PathPrefix("/crags/{key}").HandlerFunc(handlePostCrag(store)).Methods("POST")
	mux.PathPrefix("/crags/{key}").HandlerFunc(handleGetCrag(cragStore)).Methods("GET")
	// mux.HandleFunc("/", handleRoot()).Methods("GET")

}

func handleGetCrag(CragStore store.CragStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const Id = 1
		res, err := CragStore.GetCrag(Id)
		if err != nil {
			fmt.Printf("problem getting crag because of error %s", err)
		}
		err = encode(w, r, http.StatusOK, res)
		if err != nil {
			fmt.Printf("error encoding: %s", err)
		}
	}
}

func encode[T any](w http.ResponseWriter, r *http.Request, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json %w", err)
	}
	return nil

}

// func handlePostCrag(store cragStore) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		vars := mux.Vars(r)
// 		crag := vars["key"]
// 		store.addCrag(crag)
// 		w.WriteHeader(http.StatusAccepted)
// 	}
// }

// func handleRoot() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Fprint(w, "Welcome to the root path!")
// 	}

// }

// func NewServer(store cragStore) http.Handler {
// 	mux := mux.NewRouter()
// 	addRoutes(mux, store)

// 	// var handler http.Handler = mux

// 	//handler = middleware(handler)

// 	return http.Handler(mux)

// }
