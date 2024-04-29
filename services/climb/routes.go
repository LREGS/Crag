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

	r.HandleFunc("", h.Post()).Methods("POST")
	r.HandleFunc("/crag/{cragId}", h.GetByCragId()).Methods("GET")
	// r.HandleFunc("/all", h.HandleGetAllClimbs()).Methods("GET")
	// r.HandleFunc("/{Id}", h.HandleGetClimbById()).Methods("GET")
	// r.HandleFunc("/", h.HandleUpdateClimb()).Methods("PUT")
	// r.HandleFunc("/{Id}", h.HandleDeleteClimb()).Methods("DELETE")

}

func (h *Handler) Post() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		climb := models.ClimbPayload{}

		if err := util.Decode(r, &climb); err != nil {
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

func (h *Handler) GetByCragId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key, err := strconv.Atoi(vars["cragId"])
		if err != nil {
			util.WriteError(w, http.StatusInternalServerError, fmt.Errorf("could not convert key (%d) to integer: %s, URL is %v", key, err, r.URL))
			return
		}

		res, err := h.store.GetClimbsByCragId(key)
		if err != nil {
			util.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error getting crags %s", err))
			return
		}

		err = util.Encode(w, http.StatusOK, res)
		if err != nil {
			util.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error ecoding %s", err))
			return
		}

	}
}

// func (h *Handler) HandleGetAllClimbs() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		//please check this error
// 		res, _ := h.store.GetAllClimbs()

// 		err := util.Encode(w, http.StatusOK, res)
// 		if err != nil {
// 			util.WriteError(w, http.StatusBadRequest, fmt.Errorf("error encoding %s", err))
// 		}

// 		// w.WriteHeader(http.StatusOK)
// 	}

// }

// func (h *Handler) HandleGetClimbById() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		vars := mux.Vars(r)
// 		key, err := strconv.Atoi(vars["Id"])
// 		if err != nil {
// 			util.WriteError(w, http.StatusBadRequest, fmt.Errorf("error getting key becuase of error: %s", err))
// 		}
// 		res, err := h.store.GetClimbById(key)
// 		if err != nil {
// 			util.WriteError(w, http.StatusBadRequest, fmt.Errorf("could not get climb because of err: %s", err))
// 		}

// 		err = util.Encode(w, http.StatusOK, res)
// 		if err != nil {
// 			util.WriteError(w, http.StatusBadRequest, fmt.Errorf("could not encode because of error: %s", err))
// 		}
// 	}
// }

// func (h *Handler) HandleUpdateClimb() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 		updatedClimb := &models.Climb{}

// 		err := util.Decode(r, updatedClimb)
// 		if err != nil {
// 			util.WriteError(w, http.StatusBadRequest, fmt.Errorf("Coud not decode request because of err: %s", err))
// 		}

// 		resData, err := h.store.UpdateClimb(updatedClimb)
// 		if err != nil {
// 			util.WriteError(w, http.StatusBadRequest, fmt.Errorf("Could not update climb because of err: %s", err))
// 		}

// 		err = util.Encode(w, http.StatusOK, resData)
// 		if err != nil {
// 			util.WriteError(w, http.StatusBadRequest, fmt.Errorf("Coud not encode request because of err: %s", err))

// 		}

// 	}
// }

// func (h *Handler) HandleDeleteClimb() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 		vars := mux.Vars(r)
// 		key, err := strconv.Atoi(vars["Id"])
// 		if err != nil {
// 			util.WriteError(w, http.StatusBadRequest, fmt.Errorf("could not get id from url"))
// 		}

// 		err = h.store.DeleteClimb(key)
// 		if err != nil {
// 			util.WriteError(w, http.StatusBadRequest, fmt.Errorf("failed deleting climb: %s", err))
// 		}
// 		w.WriteHeader(http.StatusNoContent)

// 	}
// }
