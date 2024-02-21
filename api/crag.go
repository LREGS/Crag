package api

import "workspaces/github.com/lregs/Crag/web"

func (api *API) InitCrag() {
	api.BaseRoutes.Crag.Handle("", web.APIHandler(getCrags)),Methods.("GET")
}
