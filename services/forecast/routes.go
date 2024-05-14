package forecast

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	store "github.com/lregs/Crag/SqlStore"
	"github.com/lregs/Crag/models"
	"github.com/lregs/Crag/util"
)

type Handler struct {
	store store.ForecastStore
}

func NewHandler(store store.ForecastStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) RegisterRoutes(ctx context.Context, r *mux.Router) {
	// "/forecast"
	r.HandleFunc("", h.Post(ctx)).Methods("POST")
	// r.HandleFunc("/{Id}", h.GetByCragId(ctx)).Methods("GET")
	r.HandleFunc("/all", h.GetAllForecasts(ctx)).Methods("GET")
	r.HandleFunc("/{Id}", h.handleDeleteForecastById(ctx)).Methods("DELETE")
}

func (h *Handler) Post(ctx context.Context) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "POST" {
			http.Error(w, "Wrong request method", http.StatusMethodNotAllowed)
			return
		}

		payload := models.DBForecastPayload{}

		err := util.Decode(r, &payload)
		if err != nil {
			util.WriteError(w, http.StatusBadRequest, decodeError, err)
			return
		}

		res, err := h.store.StoreForecast(ctx, payload)
		if err != nil {
			util.WriteError(w, http.StatusInternalServerError, storeError, err)
			return
		}

		if err = util.Encode(w, http.StatusOK, &res); err != nil {
			util.WriteError(w, http.StatusBadRequest, encodeError, err)
			return
		}
	}
}

func (h *Handler) GetByCragId(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		key, err := strconv.Atoi(vars["Id"])
		if err != nil {
			util.WriteError(w, http.StatusBadRequest, varsErorr, err)
			return
		}

		data, err := h.store.GetForecastByCragId(ctx, key)
		if err != nil {
			util.WriteError(w, http.StatusInternalServerError, storeError, err)
			return
		}

		err = util.Encode(w, http.StatusOK, data)
		if err != nil {
			util.WriteError(w, http.StatusInternalServerError, encodeError, err)
			return
		}

	}
}

func (h *Handler) GetAllForecasts(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		data, err := h.store.GetAllForecastsByCragId(ctx)
		if err != nil {
			util.WriteError(w, http.StatusInternalServerError, storeError, err)
			return
		}

		if err = util.Encode(w, 200, data); err != nil {
			util.WriteError(w, http.StatusInternalServerError, encodeError, err)
			return
		}
	}
}

func (h *Handler) handleDeleteForecastById(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		key, err := strconv.Atoi(vars["Id"])
		if err != nil {
			util.WriteError(w, http.StatusBadRequest, varsErorr, err)
			return
		}

		data, err := h.store.DeleteForecastById(ctx, key)
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
