package server

import (
	"fmt"
	"net/http"
)

func CragServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "cold")
}
