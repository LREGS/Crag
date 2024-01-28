package main

// "fmt"
// "log"
// "os"

import (
	"fmt"
	client "workspaces/github.com/lregs/Crag/client"
	h "workspaces/github.com/lregs/Crag/headers"
	helpers "workspaces/github.com/lregs/Crag/helper"
	utils "workspaces/github.com/lregs/Crag/utils"
	data "workspaces/github.com/lregs/Crag/utils"

)

// "github.com/joho/godotenv"
// "workspaces/github.com/lregs/Crag/utils"

func main() {
	client := client.DefaultClient()

	coords := make(map[string][]float64)

	coords["milestone"] = []float64{53.122664, -3.998611}
	coords["gorlan"] = []float64{53.141574, -4.026437}

	url, err := helpers.MetOfficeURL(coords)
	helpers.CheckError(err)

	headers := h.ReturnHeaders()

	fmt.Println(headers)

	f, err := utils.GetForecast(url, headers, client)

	fmt.Println(f)

}
 func allForecasts(in <- chan int, out chan<- int){
	client := client.DefaultClient()
	url, err := helpers.MetOfficeURL(coords)
	helpers.CheckError(err)
	headers := h.ReturnHeaders()


	go func(){
		for val := range in {
			f, err := utils.GetForecast(url, headers, client)
			out <- result
		}
	}
 }