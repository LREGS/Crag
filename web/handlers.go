package web

import (
	"net/http"
)

type Handler struct {
	handlerFunc func(http.ResponseWriter, *http.Request)
	//requires session - when auth is implemented
}

func APIHandler(h func(http.ResponseWriter, *http.Request)) *Handler {
	return &Handler{
		handlerFunc: h,
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.handlerFunc(w, r)
}
