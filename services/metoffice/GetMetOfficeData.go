package met

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/lregs/Crag/models"
)

//do I need to have a struct that has the methods or just the functions I dont know

// returns the forecast for a crag based on its stored coords
func GetForecast(coords []float64) (models.Forecast, error) {
	var forecast models.Forecast

	client := http.Client{}

	url := fmt.Sprintf("https://data.hub.api.metoffice.gov.uk/sitespecific/v0/point/hourly?latitude=%f&longitude=%f", coords[0], coords[1])

	headers, err := getHeaders()
	if err != nil {
		return forecast, err
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return forecast, err
	}

	req.Header = http.Header{

		KEY PLEASE 
		"accept": {headers.Accept},
	}

	res, err := client.Do(req)
	if err != nil {
		return forecast, err

	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return forecast, fmt.Errorf("code %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return forecast, err

	}

	err = json.Unmarshal(body, &forecast)
	if err != nil {
		return forecast, err
	}

	// defer res.Body.Close()
	// err = json.NewDecoder(res.Body).Decode(&forecast)
	// if err != nil {
	// 	return forecast, err
	// }

	return forecast, nil

}

type metOfficeHeaders struct {
	ClientId     string `json:"X-IBM-Client-Id"`
	ClientSecret string `json:"X-IBM-Client-Secret"`
	Accept       string `json:"accept"`
}

func getHeaders() (metOfficeHeaders, error) {
	// env, err := util.GetEnv([]string{"CLIENT_ID", "CLIENT_SECRET"})
	// if err != nil {
	// 	return metOfficeHeaders{}, nil
	// }

	return metOfficeHeaders{ClientId: "e281f46e4980a322a84ed4592e2ae920", ClientSecret: "3a66408a791e46b7e6762b4370713072", Accept: "application/json"}, nil
}

func GetPayload(log *log.Logger, coords []float64) [][]interface{} {

	forecast, err := GetForecast(coords)
	if err != nil {
		log.Println(err)
	}

	timeSeries := forecast.Features[0].Properties.TimeSeries

	payload := make([][]interface{}, len(timeSeries))

	for i := 0; i < len(timeSeries); i++ {

		payload[i] = []interface{}{
			i + 1, //Id
			timeSeries[i].Time,
			timeSeries[i].ScreenTemperature,
			timeSeries[i].FeelsLikeTemperature,
			timeSeries[i].WindSpeed10m,
			timeSeries[i].WindDirectionFrom10m,
			timeSeries[i].TotalPrecipAmount,
			timeSeries[i].ProbOfPrecipitation,
			forecast.Features[0].Geometry.Coordinates[0],
			forecast.Features[0].Geometry.Coordinates[1],
		}

	}

	return payload

}
