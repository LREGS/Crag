package mw

import (
	"net/http"
	"path"
	"strconv"
)

func IDFromRequest(r *http.Request) (int, error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	return id, err
}
