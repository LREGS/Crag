package web

import (
	"net/http"
	"workspaces/github.com/lregs/Crag/app"
)

type Handler struct {
	Srv        *app.Server
	HandleFunc func(http.ResponseWriter, *http.Request)
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	appInstance := app.New()

	c := &Context
}
