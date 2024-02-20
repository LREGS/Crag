package app

import (
	"net/http"
	"workspaces/github.com/lregs/Crag/Store"

	"github.com/gorilla/mux"
)

type Server struct {
	RootRouter *mux.Router
	Router     *mux.Router
	Server     *http.Server
	// logger *log.Logger
	store *Store.SqlStore
}

func (s *Server) Store() *Store.SqlStore {
	return s.store
}

func NewServer() *Server {

	rootRouter := mux.NewRouter()
	Router := mux.NewRouter()

	HTTPServer := http.Server{Addr: "127.0.0.1:8080"}

	s := &Server{
		RootRouter: rootRouter,
		Router:     Router,
		Server:     &HTTPServer,
	}

	store, sErr := Store.New()
	if sErr != nil {
		return nil
	}
	s.store = store

	return s

}
