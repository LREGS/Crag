package crag

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	store "github.com/lregs/Crag/SqlStore"
)

type Handler struct {
	store store.CragStore
}

func NewHandler(store store.CragStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRouters(r *mux.Router) {
	r.PathPrefix("/{key}").HandlerFunc(h.handleGetCrag()).Methods("GET")
}

func (h *Handler) handleGetCrag() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		key := vars["key"]

		cragID, err := strconv.Atoi(key)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Printf("error converting key to integer: %s", err)
			return
		}
		res, err := h.store.GetCrag(cragID)
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

func encode[T any](w http.ResponseWriter, r *http.Request, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json %w", err)
	}
	return nil

}
