package climb

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	store "github.com/lregs/Crag/SqlStore"
	"github.com/lregs/Crag/models"
	"github.com/lregs/Crag/util"
)

type Handler struct {
	log   *log.Logger
	store store.ClimbStore
}

func NewHandler(log *log.Logger, store store.ClimbStore) *Handler {
	return &Handler{
		store: store,
		log:   log,
	}
}

func (h *Handler) RegisterRoutes(ctx context.Context, r *mux.Router) {

	//subrouter of /climb
	r.HandleFunc("", h.Post(ctx, h.log)).Methods("POST")
	r.HandleFunc("/crag/{cragId}", h.GetByCragId(ctx)).Methods("GET")
	r.HandleFunc("/all", h.GetAll(ctx)).Methods("GET")
	r.HandleFunc("/{Id}", h.GetById(ctx)).Methods("GET")
	r.HandleFunc("/", h.Update(ctx)).Methods("PUT")
	r.HandleFunc("/{Id}", h.Delete(ctx)).Methods("DELETE")

}

func (h *Handler) Post(ctx context.Context, log *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var climb models.ClimbPayload
		if err := util.Decode(r, &climb); err != nil {
			util.WriteError(w, http.StatusBadRequest, decodeError, err)
			return
		}

		storedData, err := h.store.StoreClimb(ctx, climb)
		if err != nil {
			log.Printf("error message is %s", err.Error())
			util.WriteError(w, http.StatusInternalServerError, storeError, err)
			// http.Error(w, fmt.Sprintf("store err %s", err), http.StatusInternalServerError)
			return
		}

		err = util.Encode(w, http.StatusCreated, storedData)
		if err != nil {
			log.Printf("error message is %s", err.Error())

			util.WriteError(w, http.StatusInternalServerError, encodeError, err)
			return
		}
	}
}

func (h *Handler) GetByCragId(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key, err := strconv.Atoi(vars["cragId"])
		if err != nil {
			util.WriteError(w, http.StatusBadRequest, varsErorr, err)
			return
		}

		res, err := h.store.GetClimbsByCragId(ctx, key)
		if err != nil {
			util.WriteError(w, http.StatusInternalServerError, storeError, err)
			return
		}

		err = util.Encode(w, http.StatusOK, res)
		if err != nil {
			util.WriteError(w, http.StatusInternalServerError, encodeError, err)
			return
		}

	}
}

func (h *Handler) GetAll(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		res, err := h.store.GetAllClimbs(ctx)
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

func (h *Handler) GetById(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		key, err := strconv.Atoi(vars["Id"])
		if err != nil {
			util.WriteError(w, http.StatusBadRequest, varsErorr, err)
			return
		}

		res, err := h.store.GetClimbById(ctx, key)
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

func (h *Handler) Update(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var updatedClimb models.Climb
		err := util.Decode(r, &updatedClimb)
		if err != nil {
			util.WriteError(w, http.StatusBadRequest, decodeError, err)
			return
		}

		resData, err := h.store.UpdateClimb(ctx, updatedClimb)
		if err != nil {
			util.WriteError(w, http.StatusInternalServerError, storeError, err)
			return
		}

		err = util.Encode(w, http.StatusOK, resData)
		if err != nil {
			util.WriteError(w, http.StatusInternalServerError, encodeError, err)
			return

		}

	}
}

func (h *Handler) Delete(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		key, err := strconv.Atoi(vars["Id"])
		if err != nil {
			util.WriteError(w, http.StatusBadRequest, varsErorr, err)
			return
		}

		res, err := h.store.DeleteClimb(ctx, key)
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
