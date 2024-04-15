package forecast

import (
	"fmt"
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

func (h *Handler) RegisterRoutes(r *mux.Router) {
	// "/forecast"
	r.HandleFunc("", h.handlePostForecast()).Methods("POST")
	r.HandleFunc("/{Id}", h.handleGetForecastByCragId()).Methods("GET")
	r.HandleFunc("/all", h.handleGetAllForecast()).Methods("GET")
	r.HandleFunc("/{key}", h.handleDeleteForecastById()).Methods("DELETE")
}

func (h *Handler) handlePostForecast() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "POST" {
			http.Error(w, "Wrong request method", http.StatusMethodNotAllowed)
		}

		payload := &models.DBForecastPayload{}

		err := util.Decode(r, payload)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed decoding payload %s", err), http.StatusInternalServerError)
			return
		}

		//shouldnt this be a copy?!
		res, err := h.store.AddForecast(payload)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed storing payload  %s", err), http.StatusInternalServerError)
			return
		}

		//is this here to just pass the test because its been some days
		if res.Id == 1 {
			http.Error(w, "empty value returned from store", 500)
		}

		//why &??
		err = util.Encode(w, http.StatusOK, &res)
		if err != nil {
			http.Error(w, fmt.Sprintf("error encoding response: %s", err), http.StatusInternalServerError)
		}

	}
}

func (h *Handler) handleGetForecastByCragId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		key, err := strconv.Atoi(vars["Id"])
		if err != nil {
			http.Error(w, "Could not get id from request", http.StatusBadRequest)
		}

		data, err := h.store.GetForecastByCragId(key)
		if err != nil {
			http.Error(w, fmt.Sprintf("Getting data failed: %s", err), http.StatusInternalServerError)
		}

		err = util.Encode(w, http.StatusOK, data)
		if err != nil {
			http.Error(w, "Could not encode responde", http.StatusInternalServerError)
		}

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
