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
	//should recieve a subrouter /climb/X

	r.HandleFunc("/", h.handlePostClimb()).Methods("POST")
	r.HandleFunc("/crag/{cragID}", h.handleGetClimbsByCrag()).Methods("GET")
	r.HandleFunc("/all", h.HandleGetAllClimbs()).Methods("GET")
	r.HandleFunc("/{Id}", h.HandleGetClimbById()).Methods("GET")
	r.HandleFunc("/", h.HandleUpdateClimb()).Methods("PUT")

}

func (h *Handler) handlePostClimb() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		climb := &models.Climb{}
		// data, err := util.Decode(r, climb)
		// if err != nil {
		// 	http.Error(w, "error decoding request body", http.StatusBadRequest)
		// }
		if err := util.Decode(r, climb); err != nil {
			errString := fmt.Sprintf("error decoding request body: %s", err)
			http.Error(w, errString, http.StatusBadRequest)
		}

		storedData, err := h.store.StoreClimb(climb)
		if err != nil {
			http.Error(w, "error storing climb", http.StatusBadRequest)
		}

		err = util.Encode(w, http.StatusCreated, storedData)
		if err != nil {
			http.Error(w, "error encoding response", http.StatusBadRequest)
		}

	}
}

func (h *Handler) handleGetClimbsByCrag() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key, err := strconv.Atoi(vars["cragID"])
		if err != nil {
			util.WriteError(w, http.StatusBadRequest, fmt.Errorf("could not convert key to integer: %s", err))
		}

		res, err := h.store.GetClimbsByCrag(key)
		if err != nil {
			util.WriteError(w, http.StatusBadRequest, fmt.Errorf("error getting crags %s", err))
		}

		for _, climb := range res {
			err := h.store.Validate(climb)
			if err != nil {
				util.WriteError(w, http.StatusBadRequest, fmt.Errorf("could not validate because of error %s", err))
			}
		}

		err = util.Encode(w, http.StatusOK, res)
		if err != nil {
			util.WriteError(w, http.StatusBadRequest, fmt.Errorf("error ecoding %s", err))
		}

	}
}

func (h *Handler) HandleGetAllClimbs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := h.store.GetAllClimbs()

		err := util.Encode(w, http.StatusOK, res)
		if err != nil {
			util.WriteError(w, http.StatusBadRequest, fmt.Errorf("error encoding %s", err))
		}

		// w.WriteHeader(http.StatusOK)
	}

}

func (h *Handler) HandleGetClimbById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key, err := strconv.Atoi(vars["Id"])
		if err != nil {
			util.WriteError(w, http.StatusBadRequest, fmt.Errorf("error getting key becuase of error: %s", err))
		}
		res, err := h.store.GetClimbById(key)
		if err != nil {
			util.WriteError(w, http.StatusBadRequest, fmt.Errorf("could not get climb because of err: %s", err))
		}

		err = util.Encode(w, http.StatusOK, res)
		if err != nil {
			util.WriteError(w, http.StatusBadRequest, fmt.Errorf("could not encode because of error: %s", err))
		}
	}
}

func (h *Handler) HandleUpdateClimb() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		updatedClimb := &models.Climb{}

		err := util.Decode(r, updatedClimb)
		if err != nil {
			util.WriteError(w, http.StatusBadRequest, fmt.Errorf("Coud not decode request because of err: %s", err))
		}

		resData, err := h.store.UpdateClimb(updatedClimb)
		if err != nil {
			util.WriteError(w, http.StatusBadRequest, fmt.Errorf("Could not update climb because of err: %s", err))
		}

		err = util.Encode(w, http.StatusOK, resData)
		if err != nil {
			util.WriteError(w, http.StatusBadRequest, fmt.Errorf("Coud not encode request because of err: %s", err))

		}

	}
}
