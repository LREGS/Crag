package helpers

import (
	"errors"
	"fmt"
)

func MetOfficeURL(coords []float32) (string, error) {
	//long[0],lat[1]
	if len(coords) > 2 {
		return "", errors.New("Too many arguments provided in coords")
	}
	if len(coords) < 2 {
		return "", errors.New("Not enough arguments provided in coords")
	}

	return fmt.Sprintf("https://api-metoffice.apiconnect.ibmcloud.com/v0/forecasts/point/hourly?latitude=%f&longitude=%f", coords[0], coords[1]), nil

}
