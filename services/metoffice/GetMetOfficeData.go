package met

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/lregs/Crag/models"
	"github.com/lregs/Crag/util"
)

//do I need to have a struct that has the methods or just the functions I dont know

// returns the forecast for a crag based on its stored coords
func GetForecast(coords []float64) (*models.Forecast, error) {
	var forecast models.Forecast

	client := http.Client{}

	url := fmt.Sprintf("https://api-metoffice.apiconnect.ibmcloud.com/v0/forecasts/point/hourly?latitude=%f&longitude=%f", coords[0], coords[1])

	headers, err := getHeaders()
	if err != nil {
		return &forecast, err
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &forecast, err
	}

	req.Header = http.Header{
		"X-IBM-Client-Id":     {headers.ClientId},
		"X-IBM-Client-Secret": {headers.ClientSecret},
		"accept":              {headers.Accept},
	}

	res, err := client.Do(req)
	if err != nil {
		return &forecast, err

	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return &forecast, err
	}

	fmt.Print(string(body))

	// var ResponseData = make(map[string]interface{})

	err = json.Unmarshal(body, &forecast)
	if err != nil {
		return &forecast, err
	}

	// defer res.Body.Close()
	// err = json.NewDecoder(res.Body).Decode(&forecast)
	// if err != nil {
	// 	return forecast, err
	// }

	return &forecast, nil

}

type metOfficeHeaders struct {
	ClientId     string `json:"X-IBM-Client-Id"`
	ClientSecret string `json:"X-IBM-Client-Secret"`
	Accept       string `json:"accept"`
}

func getHeaders() (metOfficeHeaders, error) {
	env, err := util.GetEnv([]string{"CLIENT_ID", "CLIENT_SECRET"})
	if err != nil {
		return metOfficeHeaders{}, nil
	}

	return metOfficeHeaders{ClientId: env[0], ClientSecret: env[1], Accept: "application/json"}, nil
}
