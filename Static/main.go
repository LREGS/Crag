package main

import (
	"fmt"
	"log"
	"net/http"
	"unicode/utf8"
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

// creates strings and headers for the table columns
// TODO: also needs to take windows and return the windows string too
func formatForecastData(forecast map[string]*ForecastTotals) [][]string {

	colsMap := make([][]string, len(forecast))
	columnWidths := make([]int, 5) // we want this to be able to grow as columns grow actually but we know how many columns beforehand just no magic number
	for k, v := range forecast {
		cols := []string{
			k,
			fmt.Sprintf(" Temp %d/%d/%d ", int(v.HighestTemp), int(v.LowestTemp), int(v.AvgTemp)),
			fmt.Sprintf(" Total Precip %d ", int(v.TotalPrecip)),
			fmt.Sprintf(" Wind %dmp ↓ ", int(v.AvgWindSpeed)),
			fmt.Sprintf(" 1/2/3 "),
		}
		// I think this is ok or should it be outside of the other for loop
		for i, v := range cols {
			if columnWidths[i] < utf8.RuneCountInString(v) {
				columnWidths[i] = utf8.RuneCountInString(v)
			}
		}
		// key is date
		colsMap = append(colsMap, cols)
	}

	return colsMap
}
