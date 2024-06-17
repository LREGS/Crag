package server

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	store "github.com/lregs/CragWeather/Crag/SqlStore"
	"github.com/lregs/CragWeather/Crag/services/climb"
	"github.com/lregs/CragWeather/Crag/services/crag"
)

func NewServer(ctx context.Context, log *log.Logger, store *store.SqlStore) http.Handler {
	mux := mux.NewRouter()
	subrouter := mux.PathPrefix("/api/v1").Subrouter()
	cragRouter := subrouter.PathPrefix("/crags").Subrouter()
	climbRouter := subrouter.PathPrefix("/climb").Subrouter()

	cragHandler := crag.NewHandler(store.Stores.CragStore)
	cragHandler.RegisterRoutes(ctx, cragRouter)

	climbHandler := climb.NewHandler(log, store.Stores.ClimbStore)
	climbHandler.RegisterRoutes(ctx, climbRouter)
	// addRoutes(mux, store)

	return mux
}
