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
)

// "github.com/joho/godotenv"
// "workspaces/github.com/lregs/Crag/utils"

func main() {
	client := client.DefaultClient()
	coords := []float32{53.122677, -4.013838}

	url, err := helpers.MetOfficeURL(coords)
	helpers.CheckError(err)

	headers := h.ReturnHeaders()

	fmt.Println(headers)

	f, err := utils.GetForecast(url, headers, client)

	fmt.Println(f)

}
