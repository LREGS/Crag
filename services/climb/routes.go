package climb

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
	store store.ClimbStore
}

func NewHandler(store store.ClimbStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(r *mux.Router) {

	//subrouter of /climb
	r.HandleFunc("", h.Post()).Methods("POST")
	r.HandleFunc("/crag/{cragId}", h.GetByCragId()).Methods("GET")
	r.HandleFunc("/all", h.GetAll()).Methods("GET")
	r.HandleFunc("/{Id}", h.GetById()).Methods("GET")
	r.HandleFunc("/", h.Update()).Methods("PUT")
	r.HandleFunc("/{Id}", h.Delete()).Methods("DELETE")

}

func (h *Handler) Post() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		climb := models.ClimbPayload{}

		if err := util.Decode(r, &climb); err != nil {
			util.WriteError(w, http.StatusInternalServerError, fmt.Errorf(decodeError, err))
		}

		storedData, err := h.store.StoreClimb(climb)
		if err != nil {
			util.WriteError(w, http.StatusInternalServerError, fmt.Errorf(storeError, err))
		}

		err = util.Encode(w, http.StatusCreated, storedData)
		if err != nil {
			util.WriteError(w, http.StatusInternalServerError, fmt.Errorf(encodeError, err))
		}

	}
}

func (h *Handler) GetByCragId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key, err := strconv.Atoi(vars["cragId"])
		if err != nil {
			util.WriteError(w, http.StatusBadRequest, fmt.Errorf(varsErorr, err))
			return
		}

		res, err := h.store.GetClimbsByCragId(key)
		if err != nil {
			util.WriteError(w, http.StatusInternalServerError, fmt.Errorf(storeError, err))
			return
		}

		err = util.Encode(w, http.StatusOK, res)
		if err != nil {
			util.WriteError(w, http.StatusInternalServerError, fmt.Errorf(encodeError, err))
			return
		}

	}
}

func (h *Handler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		res, err := h.store.GetAllClimbs()
		if err != nil {
			util.WriteError(w, http.StatusInternalServerError, fmt.Errorf(storeError, err))
		}

		if err = util.Encode(w, http.StatusOK, res); err != nil {
			util.WriteError(w, http.StatusInternalServerError, fmt.Errorf(encodeError, err))
		}
	}

}

func (h *Handler) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		key, err := strconv.Atoi(vars["Id"])
		if err != nil {
			util.WriteError(w, http.StatusBadRequest, fmt.Errorf(varsErorr, err))
			return
		}

		res, err := h.store.GetClimbById(key)
		if err != nil {
			util.WriteError(w, http.StatusInternalServerError, fmt.Errorf(storeError, err))
		}

		if err = util.Encode(w, http.StatusOK, res); err != nil {
			util.WriteError(w, http.StatusInternalServerError, fmt.Errorf(encodeError, err))
		}

	}
}

func (h *Handler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var updatedClimb models.Climb
		err := util.Decode(r, &updatedClimb)
		if err != nil {
			util.WriteError(w, http.StatusBadRequest, fmt.Errorf(decodeError, err))
		}

		resData, err := h.store.UpdateClimb(updatedClimb)
		if err != nil {
			util.WriteError(w, http.StatusInternalServerError, fmt.Errorf(storeError, err))
		}

		err = util.Encode(w, http.StatusOK, resData)
		if err != nil {
			util.WriteError(w, http.StatusInternalServerError, fmt.Errorf(encodeError, err))

		}

	}
}

func (h *Handler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		key, err := strconv.Atoi(vars["Id"])
		if err != nil {
			util.WriteError(w, http.StatusBadRequest, fmt.Errorf(varsErorr, err))
		}

		res, err := h.store.DeleteClimb(key)
		if err != nil {
			util.WriteError(w, http.StatusInternalServerError, fmt.Errorf(storeError, err))
		}

		if err = util.Encode(w, http.StatusOK, res); err != nil {
			util.WriteError(w, http.StatusInternalServerError, fmt.Errorf(encodeError, err))
		}

	}
}
