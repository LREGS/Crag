package api

import (
	"net/http"
	store "workspaces/github.com/lregs/Crag/dbStore"

	"github.com/gorilla/mux"
)

type Server struct {
	RootRouter *mux.Router

	//Router for API
	Router *mux.Router

	Server *http.Server

	//storage layer for the server
	store *store.SqlStore
}

func NewServer() *Server {
	rootRouter := mux.NewRouter()
	Router := mux.NewRouter()
	HTTPServer := http.Server{Addr: "127.0.0.1:8080"}
	store, err := store.New()
	if err != nil {
		panic(err)
	}

	s := &Server{
		RootRouter: rootRouter,
		Router:     Router,
		Server:     &HTTPServer,
		store:      store,
	}

	return s

}
