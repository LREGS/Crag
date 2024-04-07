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

type getResponse struct {
	Data  interface{}
	Error string
}

func (r *getResponse) GetError() string {
	return r.Error
}

func (r *getResponse) GetData() interface{} {
	return r.Data
}

func NewHanlder(store store.ForecastStore) *Handler {
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

		var payload models.DBForecastPayload

		err := util.Decode(r, &payload)
		if err != nil {
			util.WriteResponse(w, http.StatusBadRequest, nil, err.Error())
			return
		}

		//shouldnt this be a copy?!
		res, err := h.store.AddForecast(&payload)
		if err != nil {
			util.WriteResponse(w, http.StatusBadRequest, nil, err.Error())
			return
		}

		util.WriteResponse(w, http.StatusOK, res, "")

	}
}

func (h *Handler) handleGetForecastByCragId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := &getResponse{}
		vars := mux.Vars(r)
		key, err := strconv.Atoi(vars["Id"])
		if err != nil {
			response.Error = "Could not get id from request"
			util.RWriteResponse(w, http.StatusBadRequest, response)
		}

		data, err := h.store.GetForecastByCragId(key)
		if err != nil {
			response.Error = fmt.Sprintf("Store failed: %s", err)
		}

		response.Data = data

		util.WriteResponse(w, 200, data, "")

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
