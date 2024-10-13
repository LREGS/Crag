package main

import (
	"net/http"
)

func main() {

	srv := http.FileServer(http.Dir("./static"))
	if err := http.ListenAndServe(":8282", srv); err != nil {
		panic("failed starting server " + err.Error())
	}

}
