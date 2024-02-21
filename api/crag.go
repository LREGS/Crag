package api

import (
	"fmt"
	"net/http"
	"path"
	"strconv"
	"workspaces/github.com/lregs/Crag/Store"
	"workspaces/github.com/lregs/Crag/web"
)

func (api *API) InitCrag() {
	api.BaseRoutes.Crag.Handle("", web.APIHandler(getCrag)).Methods("GET") ///crags shows all crags
}

func getCrag(w http.ResponseWriter, r *http.Request, store *Store.SqlStore) {
	//request should look like api/crag/id
	/*this is worrying pattern because user could have any input as id
	and that is going to our db as a command and it doesn't feel like we're totally in control of what is happening within this query
	maybe this is normally handled elswhere, or maybe I should be veryfying this here?! Or using context to pass the id reaquired?*/
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		//this isnt right but will work for now - we need to be implementing a log
		fmt.Println(err)
	}

	crag, err := api.Deps.Store.Stores.Crag.GetCrag(id)

}
