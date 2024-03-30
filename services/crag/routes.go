package crag

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	store "github.com/lregs/Crag/SqlStore"
	"github.com/lregs/Crag/models"
)

type Handler struct {
	store store.CragStore
}

func NewHandler(store store.CragStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	// "crags/..."
	r.HandleFunc("/", h.handlePostCrag()).Methods("POST")
	r.PathPrefix("/{key}").HandlerFunc(h.handleGetCrag()).Methods("GET")
	r.PathPrefix("/{key}").HandlerFunc(h.handleDelCragById()).Methods("DELETE")
	r.PathPrefix("/{key}").HandlerFunc(h.handlePostCrag()).Methods("POST")
}

func (h *Handler) handlePostCrag() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		crag := &models.Crag{}

		data, err := decode(r, crag)
		if err != nil {
			http.Error(w, "error decoding request body", http.StatusBadRequest)
		}

		err = h.store.StoreCrag(data)
		if err != nil {
			http.Error(w, "Could not store crag", http.StatusBadRequest)
		}

	}
}

// this isnt really complete because im usre we're going to want to get forecast with crag
// but at the same time forecast is linked to cragID so I guess when we get crag to display on front end
// you just get the corresponding forecast at the same time this way their storage etc is seperate and linked
// only through the integer key
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

func (h *Handler) handleDelCragById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		IdStr := vars["key"]

		Id, err := strconv.Atoi(IdStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		err = h.store.DeleteCragByID(Id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

	}
}

func encode[T any](w http.ResponseWriter, r *http.Request, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func decode[T any](r *http.Request, v T) (T, error) {
	// var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}
