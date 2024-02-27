package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewServer(cragStore *cragStore) http.Handler {
	mux := mux.NewRouter()

}
