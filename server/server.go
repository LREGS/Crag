package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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

		vars := mux.Vars(r)
		key := vars["key"]

		cragID, err := strconv.Atoi(key)
		if err != nil {
			fmt.Errorf("couldnt convert key to int: %s", err)
		}

		res, err := CragStore.GetCrag(cragID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Printf("problem getting crag because of error %s", err)
			return
		}
		if res == nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Printf("Crag with Id %d not found", cragID)
			return
		}

		err = encode(w, r, http.StatusOK, res)
		if err != nil {
			fmt.Printf("error encoding: %s", err)
			w.WriteHeader(http.StatusNotFound)

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
