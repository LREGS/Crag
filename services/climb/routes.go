package climb

import (
	"net/http"

	"github.com/gorilla/mux"
	store "github.com/lregs/Crag/SqlStore"
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

}

func (h *Handler) handlePostClimb() {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
