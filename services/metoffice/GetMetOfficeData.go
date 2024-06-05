package met

import (
	"encoding/json"
	"errors"
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

	//this should be recieving a client so im not making a new one with every request plls
	client := http.Client{}

	url := fmt.Sprintf("https://data.hub.api.metoffice.gov.uk/sitespecific/v0/point/hourly?latitude=%f&longitude=%f", coords[0], coords[1])

	//need this back online

	// if err := godotenv.Load(); err != nil {
	// 	return forecast, err
	// }

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return forecast, err
	}

	req.Header = http.Header{

		"apikey": {"eyJ4NXQiOiJOak16WWpreVlUZGlZVGM0TUdSalpEaGtaV1psWWpjME5UTXhORFV4TlRZM1ptRTRZV1JrWWc9PSIsImtpZCI6ImdhdGV3YXlfY2VydGlmaWNhdGVfYWxpYXMiLCJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJzdWIiOiJ3aWxsaWFtLWN1bGx1bUBob3RtYWlsLmNvLnVrQGNhcmJvbi5zdXBlciIsImFwcGxpY2F0aW9uIjp7Im93bmVyIjoid2lsbGlhbS1jdWxsdW1AaG90bWFpbC5jby51ayIsInRpZXJRdW90YVR5cGUiOm51bGwsInRpZXIiOiJVbmxpbWl0ZWQiLCJuYW1lIjoic2l0ZV9zcGVjaWZpYy1hYjE3ZTkyMy1kODI2LTQ5ZDQtYWZhMC01ODY3ZTQxODMwNzciLCJpZCI6NDIwNCwidXVpZCI6IjA4NzBjYjI5LWIyNDYtNDE3OS05OWQzLTE1ZDg5Njc5MGE0MSJ9LCJpc3MiOiJodHRwczpcL1wvYXBpLW1hbmFnZXIuYXBpLW1hbmFnZW1lbnQubWV0b2ZmaWNlLmNsb3VkOjQ0M1wvb2F1dGgyXC90b2tlbiIsInRpZXJJbmZvIjp7IndkaF9zaXRlX3NwZWNpZmljX2ZyZWUiOnsidGllclF1b3RhVHlwZSI6InJlcXVlc3RDb3VudCIsImdyYXBoUUxNYXhDb21wbGV4aXR5IjowLCJncmFwaFFMTWF4RGVwdGgiOjAsInN0b3BPblF1b3RhUmVhY2giOnRydWUsInNwaWtlQXJyZXN0TGltaXQiOjAsInNwaWtlQXJyZXN0VW5pdCI6InNlYyJ9fSwia2V5dHlwZSI6IlBST0RVQ1RJT04iLCJzdWJzY3JpYmVkQVBJcyI6W3sic3Vic2NyaWJlclRlbmFudERvbWFpbiI6ImNhcmJvbi5zdXBlciIsIm5hbWUiOiJTaXRlU3BlY2lmaWNGb3JlY2FzdCIsImNvbnRleHQiOiJcL3NpdGVzcGVjaWZpY1wvdjAiLCJwdWJsaXNoZXIiOiJKYWd1YXJfQ0kiLCJ2ZXJzaW9uIjoidjAiLCJzdWJzY3JpcHRpb25UaWVyIjoid2RoX3NpdGVfc3BlY2lmaWNfZnJlZSJ9XSwidG9rZW5fdHlwZSI6ImFwaUtleSIsImlhdCI6MTcxNTcxMjE4MywianRpIjoiYjFkMDZjZDctZDViNy00OWFkLThhM2YtNzhjMDRkZjZjM2ZhIn0=.PzEmpYP8PqjABpjJN8z4LbfILUgALEybfjJIJ2IrV9gwWV9GoLTGjqFrajJ0QkSJID80HuOlDLp0psYZLgSVe-l1DpAe6FtwYoBP6TaTN8PdiLle5m7JFlIR-sYd_iXDHUpAnjWrNh1u_Ofz8bPcQZ8F5szg9DbZQL_umsw-ST5L01tl3PEmqBkZiZ20rCMTxK9OMpoukfX2iPX4US48sIny6XVExLAXvZXt4uFqnChInFJXkIALksndiUm8OL_sDwbraPZKd1MUgII-SBtbJcp-nWCI3J9oNvKrK53HlqOWPTTcrsLnZkcJLmhQbvXXEgqyRnJM5Usa297EveFWYA=="},
		"accept": {"application/json"},
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

// type metOfficeHeaders struct {
// 	ApiKey string `json:"apikey"`
// 	Accept string `json:"accept"`
// }

// func getHeaders() (metOfficeHeaders, error) {
// 	env, err := util.GetEnv([]string{"apikey"})
// 	if err != nil {
// 		return metOfficeHeaders{}, nil
// 	}

// 	return metOfficeHeaders{apikey: env[0], Accept: "application/json"}, nil
// }

func GetPayload(log *log.Logger, coords []float64) ([][]interface{}, error) {

	//if get forecast fails we get an index out of range error because of the timeSeries
	//im not sure why the error is obviously being returned as nil but tis annoying

	forecast, err := GetForecast(coords)
	if err != nil {
		log.Println(err)
	}

	if len(forecast.Features) == 0 {
		return nil, errors.New("empty forecast")
	}

	timeSeries := forecast.Features[0].Properties.TimeSeries

	//not sure I need [][]interface{} as this was for sql copy
	payload := make([][]interface{}, len(timeSeries))

	for i := 0; i < (len(timeSeries) - 1); i++ {

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

	return payload, nil

}

// this doesnt show prob of precipitation because that needs to be hourly, not totals

//These stats will provide header // outline stats for each crag but the full hourly will also be available on inspection I guess for all crags
//but a condensed leaderboard style report will also be available

// window is a weather window and the average weather conditions within the weather window are shown.
// I'm not sure how useful the whole scale averages are for multi day windows but I guess the full hourly forecast will still be presented somewhere alongside the windows

// This will represent a single window in time detailed by the time string.
// each window is a minimum of 2 hours long
type Window struct {
	Time          []string
	AvgTemp       float64
	AvgWindSpeed  float64
	WindDirection int
}

type ForecastTotals struct {
	HighestTemp   float64
	LowestTemp    float64
	AvgTemp       float64
	AvgWindSpeed  float64
	WindDirection int
	TotalPrecip   float64
	Windows       []Window
}

func GetRedisPayload(log *log.Logger, forecast models.Forecast) (map[string][]ForecastTotals, error) {

	//if get forecast fails we get an index out of range error because of the timeSeries
	//im not sure why the error is obviously being returned as nil but tis annoying

	if len(forecast.Features) == 0 {
		return nil, errors.New("empty forecast")
	}

	timeSeries := forecast.Features[0].Properties.TimeSeries

	totals := make(map[string][]ForecastTotals, len(timeSeries))

	time := timeSeries[0].Time[8:10]

	total := ForecastTotals{
		HighestTemp:   0.0,
		LowestTemp:    timeSeries[0].ScreenTemperature,
		AvgWindSpeed:  0.00,
		WindDirection: 0,
		TotalPrecip:   0,
		Windows:       []Window{},
	}

	for i := 0; i < (len(timeSeries) - 1); i++ {

		if time != timeSeries[i].Time[8:10] {
			totals[time] = append(totals[time], total)
			time = timeSeries[i].Time[8:10]
			total = ForecastTotals{
				HighestTemp:   0.0,
				LowestTemp:    timeSeries[0].ScreenTemperature,
				AvgWindSpeed:  0.00,
				WindDirection: 0,
				TotalPrecip:   0,
				Windows:       []Window{}}
		}

		windows := findWindows(timeSeries)
		total.Windows = windows

		if timeSeries[i].ScreenTemperature > total.HighestTemp {
			total.HighestTemp = timeSeries[i].ScreenTemperature
		}

		if timeSeries[i].ScreenTemperature < total.LowestTemp {
			total.LowestTemp = timeSeries[i].ScreenTemperature
		}

		// if timeSeries[i].TotalPrecipAmount != 0.00 {
		// 	total.TotalPrecip += timeSeries[i].TotalPrecipAmount
		// 	window = Window{}
		// } else {
		// 	window. = append(window, timeSeries[i].TotalPrecipAmount)
		// }

		// if timeSeries[i].TotalPrecipAmount > 0.00 {

		// } else {
		// 	if !(timeSeries[i+1].TotalPrecipAmount > 0.00) {
		// 		window = append(window, timeSeries[i].Time, timeSeries[i+1].Time)
		// 		windows = append(windows, window)
		// 	} else {
		// 		total.Windows = windows
		// 		window = []string{}
		// 	}
		// }

		total.AvgWindSpeed = total.AvgWindSpeed + timeSeries[i].WindSpeed10m
		total.AvgWindSpeed = total.AvgWindSpeed / float64(len(timeSeries))

	}

	return totals, nil

}

func findWindows(data []models.TimeSeriesData) []Window {
	windows := []Window{}
	time := []string{}

	win := Window{}

	for _, d := range data {
		if d.TotalPrecipAmount == 0 {
			time = append(time, d.Time)
			win.AvgTemp += d.ScreenTemperature
			win.AvgWindSpeed += d.WindSpeed10m
			//need to find a way to generalise
			win.WindDirection = d.WindDirectionFrom10m
			continue
		} else {
			win.Time = time
			windows = append(windows, win)
			win = Window{}
			continue
		}

	}

	return windows

}
