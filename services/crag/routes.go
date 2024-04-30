package crag

import (
	"net/http"

	"github.com/gorilla/mux"
	store "github.com/lregs/Crag/SqlStore"
	"github.com/lregs/Crag/models"
	"github.com/lregs/Crag/util"
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
	// r.PathPrefix("/{key}").HandlerFunc(h.handleGetCrag()).Methods("GET")
	// r.PathPrefix("/{key}").HandlerFunc(h.handleDelCragById()).Methods("DELETE")
	// r.PathPrefix("/{key}").HandlerFunc(h.handlePostCrag()).Methods("POST")
}

func (h *Handler) Post() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var crag models.Crag

		data, err := util.Decode(r, crag)
		if err != nil {
			http.Error(w, "error decoding request body", http.StatusBadRequest)
		}

		err = h.store.StoreCrag(crag)
		if err != nil {
			http.Error(w, "Could not store crag", http.StatusBadRequest)
		}

	}
}

// func (h *Handler) GetById() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 		vars := mux.Vars(r)
// 		key := vars["key"]

// 		cragID, err := strconv.Atoi(key)
// 		if err != nil {
// 			w.WriteHeader(http.StatusBadRequest)
// 			fmt.Printf("error converting key to integer: %s", err)
// 			return
// 		}
// 		res, err := h.store.GetCrag(cragID)
// 		if err != nil {
// 			w.WriteHeader(http.StatusNotFound)
// 			fmt.Printf("problem getting crag because of error %s", err)
// 			return
// 		}

// 		err = encode(w, r, http.StatusOK, res)
// 		if err != nil {
// 			fmt.Printf("error encoding: %s", err)
// 			w.WriteHeader(http.StatusNotFound)

// 		}

// 	}
// }

// func (h *Handler) DelById() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		vars := mux.Vars(r)
// 		IdStr := vars["key"]

// 		Id, err := strconv.Atoi(IdStr)
// 		if err != nil {
// 			w.WriteHeader(http.StatusBadRequest)
// 		}

// 		err = h.store.DeleteCragByID(Id)
// 		if err != nil {
// 			w.WriteHeader(http.StatusBadRequest)
// 		}

// 	}
// }
