package crag

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	store "github.com/lregs/CragWeather/Crag/SqlStore"
	"github.com/lregs/CragWeather/Crag/models"
	"github.com/lregs/CragWeather/Crag/util"
)

type Handler struct {
	store store.CragStore
}

func NewHandler(store store.CragStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(ctx context.Context, r *mux.Router) {
	// "crags/...
	r.HandleFunc("/", h.Post(ctx)).Methods("POST")
	r.PathPrefix("/{key}").HandlerFunc(h.GetById(ctx)).Methods("GET")
	r.PathPrefix("/{key}").HandlerFunc(h.DeleteById(ctx)).Methods("DELETE")
}

func (h *Handler) Post(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var crag models.CragPayload
		if err := util.Decode(r, &crag); err != nil {
			util.WriteError(w, http.StatusInternalServerError, decodeError, err)
			return
		}

		stored, err := h.store.StoreCrag(ctx, crag)
		if err != nil {
			util.WriteError(w, http.StatusInternalServerError, storeError, err)
			return
		}

		if err := util.Encode(w, 200, stored); err != nil {
			util.WriteError(w, http.StatusInternalServerError, encodeError, err)
			return
		}

	}
}

func (h *Handler) GetById(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		cragID, err := strconv.Atoi(vars["key"])
		if err != nil {
			util.WriteError(w, http.StatusBadRequest, varsErorr, err)
			return
		}

		res, err := h.store.GetCrag(ctx, cragID)
		if err != nil {
			util.WriteError(w, http.StatusInternalServerError, storeError, err)
			return

		}

		if err = util.Encode(w, http.StatusOK, res); err != nil {
			util.WriteError(w, http.StatusInternalServerError, encodeError, err)
			return

		}

	}
}

func (h *Handler) DeleteById(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		Id, err := strconv.Atoi(vars["key"])
		if err != nil {
			util.WriteError(w, http.StatusInternalServerError, varsErorr, err)
			return
		}

		data, err := h.store.DeleteCragByID(ctx, Id)
		if err != nil {
			util.WriteError(w, http.StatusInternalServerError, storeError, err)
			return
		}

		if err := util.Encode(w, 200, data); err != nil {
			util.WriteError(w, http.StatusInternalServerError, encodeError, err)
			return
		}
	}
}
