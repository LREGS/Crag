package main

import (
	"log"
	"net/http"
)

func main() {

	server := &CragServer{NewInMemoryCragtStore()}

	log.Fatal(http.ListenAndServe(":6969", server))
}
