package main

import (
	"fmt"
	// "log"
	// "os"

	// "github.com/joho/godotenv"
	"workspaces/github.com/lregs/Crag/utils"
)

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("error loading .env")
	// }
	// id := os.Getenv("CLIENT_ID")
	// secret := os.Getenv("CLIENT_SECRET")

	// fmt.Printf(id, secret)
	var coords = []float32{53.089600, -4.049700}

	forecast := utils.GetForecast(coords)
	fmt.Print(forecast["features"])
}
