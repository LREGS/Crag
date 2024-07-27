package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// type MetOfficeAPI interface{

// }

const baseHourlyURL string = "https://data.hub.api.metoffice.gov.uk/sitespecific/v0/point/hourly?"

type api interface {
	MakeAPICall(url string, headers http.Header) ([]byte, error)
}

type Mapi struct {
	client http.Client
}

func (a *Mapi) MakeAPICall(url string, header http.Header) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := a.client.Do(req)
	if err != nil {
		return nil, err

	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return b, nil

}

type MetOfficeAPI struct {
	log     *log.Logger
	BaseURL string
	APIKey  string
	Client  http.Client
	Mapi
}

func NewMetAPI(apikey string, log *log.Logger) *MetOfficeAPI {
	return &MetOfficeAPI{log: log, BaseURL: baseHourlyURL, APIKey: apikey, Client: http.Client{}}
}

// Returns the hourly forecast for a 72hour period from the met office data hub api, hourly
func (mAPI *MetOfficeAPI) GetForecast(url string) (Forecast, error) {

	var forecast Forecast

	// url := mAPI.BaseURL + fmt.Sprintf("latitude=%f&longitude=%f", coords[0], coords[1])

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return forecast, err
	}

	req.Header = http.Header{

		"apikey": {mAPI.APIKey},
		"accept": {"application/json"},
	}

	res, err := mAPI.Client.Do(req)
	if err != nil {
		return forecast, err

	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return forecast, fmt.Errorf("code %d", res.StatusCode)
	}

	if err := json.NewDecoder(res.Body).Decode(&forecast); err != nil {
		return forecast, err
	}

	return forecast, nil

}

func (mAPI *MetOfficeAPI) CreateURL(coords []float64) string {
	return fmt.Sprintf("%slatitude=%f&longitude=%f", mAPI.BaseURL, coords[0], coords[1])
}

func (mAPI *MetOfficeAPI) GetHeaders() http.Header {
	return http.Header{

		"apikey": {mAPI.APIKey},
		"accept": {"application/json"},
	}
}

func (mAPI *MetOfficeAPI) GetPayload(log *log.Logger, forecast Forecast) (ForecastPayload, error) {

	//1 hourly spot has three days worth of data. This function provides the totals for
	//all three days

	data := forecast.Features[0].Properties.TimeSeries

	if len(data) == 0 {
		return ForecastPayload{}, errors.New("Forecast provided is empty")
	}

	totals := map[string]*ForecastTotals{}

	//we can try and calculate a running average using the mean calculator and then adding that value as the avg instead of current
	//clumsyness

	for _, val := range data {

		_, ok := totals[val.Time[8:10]]
		if !ok {
			totals[val.Time[8:10]] = &ForecastTotals{}
			entry := totals[val.Time[8:10]]
			entry.AvgTemp += val.ScreenTemperature
			entry.AvgWindSpeed += val.WindSpeed10m
			entry.HighestTemp = val.MaxScreenAirTemp
			entry.LowestTemp = val.MinScreenAirTemp
			//this is just assiging it as the first given value
			entry.WindDirection = val.WindDirectionFrom10m
			entry.TotalPrecip += val.TotalPrecipAmount
		} else {
			entry := totals[val.Time[8:10]]
			entry.AvgTemp += val.ScreenTemperature
			entry.AvgWindSpeed += val.WindSpeed10m
			entry.TotalPrecip += val.TotalPrecipAmount
		}

	}

	for day, entry := range totals {
		entry.AvgTemp = float64((entry.AvgTemp / float64((len(data) / 3))))
		entry.AvgWindSpeed = float64((entry.AvgWindSpeed / float64((len(data) / 3))))
		totals[day] = entry

	}

	return ForecastPayload{LastModelRunTime: forecast.Features[0].Properties.ModelRunDate, ForecastTotals: totals, Windows: mAPI.FindWindows(log, forecast)}, nil

}

func (mAPI *MetOfficeAPI) FindWindows(log *log.Logger, forecast Forecast) [][]time.Time {

	//probably should be returning and handling our errors in this

	var startOfWindow string
	var endOfWindow string

	windows := [][]time.Time{}
	for i := 0; i < len(forecast.Features[0].Properties.TimeSeries); i++ {

		currForecast := forecast.Features[0].Properties.TimeSeries[i]

		if currForecast.TotalPrecipAmount != 0.00 {
			if startOfWindow != "" {
				endOfWindow = currForecast.Time
				startTime, err := Str2Time(startOfWindow)
				if err != nil {
					log.Printf("faield converting start string %s", err)
					continue
				}
				endTime, err := Str2Time(endOfWindow)
				if err != nil {
					log.Printf("faield converting end string %s", err)
					continue
				}
				windows = append(windows, []time.Time{startTime, endTime})
				//	reset pointer
				// log.Println(startOfWindow, endOfWindow)
				startOfWindow = ""

			}

			continue
		}

		if startOfWindow == "" {
			startOfWindow = currForecast.Time
		}

		if i == (len(forecast.Features[0].Properties.TimeSeries) - 1) {
			if currForecast.TotalPrecipAmount != 0.00 && startOfWindow != "" {
				endOfWindow = forecast.Features[0].Properties.TimeSeries[i-1].Time
				startTime, err := Str2Time(startOfWindow)
				if err != nil {
					log.Printf("faield converting start string %s", err)
					continue
				}
				endTime, err := Str2Time(endOfWindow)
				if err != nil {
					log.Printf("faield converting end string %s", err)
					continue
				}
				windows = append(windows, []time.Time{startTime, endTime})
			}

			if startOfWindow != "" {
				endOfWindow = forecast.Features[0].Properties.TimeSeries[i-1].Time
				startTime, err := Str2Time(startOfWindow)
				if err != nil {
					log.Printf("faield converting start string %s", err)
					continue
				}
				endTime, err := Str2Time(endOfWindow)
				if err != nil {
					log.Printf("faield converting end string %s", err)
					continue
				}
				windows = append(windows, []time.Time{startTime, endTime})
			}

		}

		// log.Println(startOfWindow, endOfWindow)

	}

	return windows
}

// TODO: Add track previous potential windows, to somehow examine likelihood of when certain place could be dry based on previous weather
