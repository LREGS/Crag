package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// type ForecastData struct{
//   Features []Feature `json:"features"`

// }
// type Feature struct{
//   Geometry...
//   TimeSeries []TimeSeriesEnty
// }

// type TimeSeriesEntry struct{
//   temp
//   precipitation etc
// }
// maybe worth trying to unmarshall the json into structs to make it easier to pass around and access the data vs a map

func GetForecast(coords []float32, headers map[string]interface{}) map[string]interface{} {
	err := godotenv.Load()
	if err != nil {
		log.Println("error loading env")
	}
	client := defaultClient()

	apiUrl := fmt.Sprintf(("https://api-metoffice.apiconnect.ibmcloud.com/v0/forecasts/point/hourly?latitude=%f&longitude=%f"), coords[1], coords[0])

	req, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		log.Println(err)
	}

	jsonHeaders, err := json.Marshal(headers)
	if err != nil {
		log.Println(err)
	}
	req.Header.Add("Headers", string(jsonHeaders))

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	if res.StatusCode != http.StatusOK {
		log.Println(fmt.Sprintf("Status %v", res.Status))
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}

	var forecastData map[string]interface{}
	err = json.Unmarshal(body, &forecastData)
	if err != nil {
		log.Println(err)
	}
	return forecastData
}

func defaultClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
	}
}

func createMetOfficeRequest(coords []float32, apiUrl string) (*http.Request, error) {

}
func sendRequest(req *http.Request, client *http.Client) (*http.Response, error) {
}
func parseResponse(res *http.Response) (map[string]interface{}, error) {
}
func KGetForecast(coords []float32, headers map[string]interface{}) map[string]interface{} {
}

func GetEnv(variables []string) (map[string]string, error) {

	if len(variables) == 0 {
		return nil, errors.New("No variables passed to GetEnv Function")
	}

	err := godotenv.Load()
	if err != nil {
		log.Println("error loading env")
	}

	envVariables := make(map[string]string)

	for _, variable := range variables {
		envVariables[variable] = os.Getenv(variable)
	}
	return envVariables, nil
}
