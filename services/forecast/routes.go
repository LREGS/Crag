package forecast

import (
	"net/http"

	"github.com/gorilla/mux"
	store "github.com/lregs/Crag/SqlStore"
)

type Handler struct {
	store store.ForecastStore
}

func NewHanlder(store store.ForecastStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) RegisterRouts(r *mux.Router) {
	// "/forecast"
	r.HandleFunc("", h.handlePostForecast()).Methods("POST")
	r.HandleFunc("/{key}", h.handleGetForecastByCragId()).Methods("GET")
	r.HandleFunc("/all", h.handleGetAllForecast()).Methods("GET")
	r.HandleFunc("/{key}", h.handleDeleteForecastById()).Methods("DELETE")
}

func (h *Handler) handlePostForecast() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (h *Handler) handleGetForecastByCragId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (h *Handler) handleGetAllForecast() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (h *Handler) handleDeleteForecastById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
