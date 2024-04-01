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
