package server

import (
	"net/http"

	"github.com/gorilla/mux"
	store "github.com/lregs/Crag/SqlStore"
	"github.com/lregs/Crag/services/climb"
	"github.com/lregs/Crag/services/crag"
)

// we need to pass the whole sql store, not just crag store pls
func NewServer(store *store.SqlStore) http.Handler {
	mux := mux.NewRouter()
	subrouter := mux.PathPrefix("/api/v1").Subrouter()
	cragRouter := subrouter.PathPrefix("/crags").Subrouter()
	climbRouter := subrouter.PathPrefix("/climb").Subrouter()

	cragHandler := crag.NewHandler(store.Stores.CragStore)
	cragHandler.RegisterRoutes(cragRouter)

	climbHandler := climb.NewHandler(store.Stores.ClimbStore)
	climbHandler.RegisterRoutes(climbRouter)
	// addRoutes(mux, store)

	return mux
}
