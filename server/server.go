package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	store "github.com/lregs/Crag/SqlStore"
	"github.com/lregs/Crag/models"
)

// we need to pass the whole sql store, not just crag store pls
func NewServer(store store.CragStore) http.Handler {
	mux := mux.NewRouter()
	subrouter := mux.PathPrefix("/api/v1").Subrouter()
	cragRouter := subrouter.PathPrefix("/crags").Subrouter()

	addRoutes(mux, store)
	addCragRoutes(cragRouter, store)

	return mux
}

func addCragRoutes(r *mux.Router, s store.CragStore) {
	r.PathPrefix("/{key}").HandlerFunc(handleGetCrag(s)).Methods("GET")
}
func addRoutes(mux *mux.Router, cragStore store.CragStore) {

	// mux.PathPrefix("/crags/{key}").HandlerFunc(handlePostCrag(store)).Methods("POST")
	mux.PathPrefix("/crags/{key}").HandlerFunc(handleGetCrag(cragStore)).Methods("GET")
	mux.PathPrefix("/crags/{key}").HandlerFunc(handleDelCragById(cragStore)).Methods("DELETE")
	mux.PathPrefix("/crags").HandlerFunc(handlePostCrag(cragStore)).Methods("POST")
	// mux.HandleFunc("/", handleRoot()).Methods("GET")

}

func handleDelCragById(CragStore store.CragStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		IdStr := vars["key"]

		Id, err := strconv.Atoi(IdStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		err = CragStore.DeleteCragByID(Id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

	}
}

func handleGetCrag(CragStore store.CragStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		key := vars["key"]

		cragID, err := strconv.Atoi(key)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Printf("error converting key to integer: %s", err)
			return
		}
		res, err := CragStore.GetCrag(cragID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Printf("problem getting crag because of error %s", err)
			return
		}

		err = encode(w, r, http.StatusOK, res)
		if err != nil {
			fmt.Printf("error encoding: %s", err)
			w.WriteHeader(http.StatusNotFound)

		}
	}
}

func handlePostCrag(CragStore store.CragStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var crag models.Crag

		//not sure if mental block or miss-understanding but could not
		//get the decode function to work because of could not infer type.
		err := json.NewDecoder(r.Body).Decode(&crag)
		fmt.Print(crag)
		if err != nil {
			http.Error(w, "error decoding request body", http.StatusBadRequest)
		}

		err = CragStore.StoreCrag(&crag)
		if err != nil {
			http.Error(w, "Could not store crag", http.StatusBadRequest)
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

// func decode[T any](r *http.Request) (T, error) {
// 	var v T
// 	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
// 		return v, fmt.Errorf("decode json %w", err)
// 	}
// 	return v, nil
// }
