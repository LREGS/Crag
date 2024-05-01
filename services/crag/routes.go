package crag

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	store "github.com/lregs/Crag/SqlStore"
	"github.com/lregs/Crag/util"
)

type Handler struct {
	store store.CragStore
}

func NewHandler(store store.CragStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	// "crags/...
	// r.HandleFunc("/", h.Post()).Methods("POST")
	r.PathPrefix("/{key}").HandlerFunc(h.GetById()).Methods("GET")
	// r.PathPrefix("/{key}").HandlerFunc(h.handleDelCragById()).Methods("DELETE")
	// r.PathPrefix("/{key}").HandlerFunc(h.handlePostCrag()).Methods("POST")
}

// func (h *Handler) Post() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 		var crag models.CragPayload

// 		if err := util.Decode(r, &crag); err != nil {
// 			util.WriteError(w, http.StatusInternalServerError, fmt.Errorf(decodeError, err))
// 		}

// 		stored, err := h.store.StoreCrag(crag)
// 		if err != nil {
// 			util.WriteError(w, http.StatusInternalServerError, fmt.Errorf(storeError, err))
// 		}

// 	}
// }

func (h *Handler) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		cragID, err := strconv.Atoi(vars["key"])

		res, err := h.store.GetCrag(cragID)
		if err != nil {
			util.WriteError(w, http.StatusInternalServerError, fmt.Errorf(storeError, err))
		}

		if err = util.Encode(w, http.StatusOK, res); err != nil {
			util.WriteError(w, http.StatusInternalServerError, fmt.Errorf(encodeError, err))

		}

	}
}

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
