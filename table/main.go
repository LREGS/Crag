package table

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	srv := http.FileServer(http.Dir("./static"))

	//func to pull data from redis db
	//use formatForecastData to get lists of strings

	fmt.Println("Starting Static HTML server on :6996")
	if err := http.ListenAndServe(":6996", srv); err != nil {
		log.Println("failed " + err.Error())
		panic("failed starting server " + err.Error())
	}

}
