package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	rate "golang.org/x/time/rate"
)

// type MetOfficeAPI interface{

// }

type MetAPI interface {
	GetForecast(url string) (Forecast, error)
	CreateURL(coords []float64) string
	GetHeaders() http.Header
	CalculateTotals(data []TimeSeriesData) map[string]*ForecastTotals
	FindWindows(data []TimeSeriesData) [][]time.Time
	GetRateLimit() rate.Limit
}

const baseHourlyURL string = "https://data.hub.api.metoffice.gov.uk/sitespecific/v0/point/hourly?"

type api interface {
	Call(url string, headers http.Header) (*http.Response, error)
}

type Api struct {
	client http.Client
}

func NewApi() *Api {
	return &Api{client: http.Client{}}
}

func (a *Api) Call(url string, header http.Header) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if req.Header != nil {
		req.Header = header
	}

	res, err := a.client.Do(req)
	if err != nil {
		return nil, err

	}

	return res, nil

}

type MetOfficeAPI struct {
	log     *log.Logger
	BaseURL string
	// I dont understand having this and not just returning from env?
	APIKey      string
	client      api
	store       Store
	ratelimiter *rate.Limiter
}

func NewMetAPI(apikey string, log *log.Logger) *MetOfficeAPI {
	return &MetOfficeAPI{
		log:         log,
		BaseURL:     baseHourlyURL,
		APIKey:      apikey,
		client:      NewApi(),
		ratelimiter: rate.NewLimiter(rate.Every(time.Minute/100), 100),
	}
}

func (mAPI *MetOfficeAPI) GetRateLimit() rate.Limit {
	return mAPI.ratelimiter.Limit()
}

// Returns the hourly forecast for a 72hour period from the met office data hub api, hourly
func (mAPI *MetOfficeAPI) GetForecast(url string) (Forecast, error) {

	var forecast Forecast

	res, err := mAPI.client.Call(url, mAPI.GetHeaders())
	if err != nil {
		return Forecast{}, nil
	}

	if res.StatusCode != 200 {
		return forecast, fmt.Errorf("failed getting forecast: code %d", res.StatusCode)
	}

	if err := json.NewDecoder(res.Body).Decode(&forecast); err != nil {
		return forecast, err
	}

	return forecast, nil

}

func (mAPI *MetOfficeAPI) CreateURL(coords []float64) string {
	// I probbaly just want to store the urls with the crag information that way I only need to create them once
	var sb strings.Builder
	sb.WriteString(mAPI.BaseURL)
	sb.WriteString("latitude=")
	sb.WriteString(strconv.FormatFloat(coords[0], 'f', -1, 64))
	sb.WriteString("&")
	sb.WriteString("longitude=")
	sb.WriteString(strconv.FormatFloat(coords[1], 'f', -1, 64))
	return sb.String()
}

func (mAPI *MetOfficeAPI) GetHeaders() http.Header {
	return http.Header{

		"apikey": {mAPI.APIKey},
		"accept": {"application/json"},
	}
}

func (mAPI *MetOfficeAPI) CalculateTotals(data []TimeSeriesData) map[string]*ForecastTotals {

	var totals = make(map[string]*ForecastTotals)

	for _, val := range data {
		day := val.Time[8:10]
		entry, ok := totals[day]
		if !ok {
			entry = &ForecastTotals{
				HighestTemp: val.MaxScreenAirTemp,
				LowestTemp:  val.MinScreenAirTemp,
				Datapoints:  0,
			}
			totals[day] = entry
		}

		entry.AvgTemp += val.ScreenTemperature
		entry.AvgWindSpeed += val.WindSpeed10m
		entry.TotalPrecip += val.TotalPrecipAmount
		entry.Datapoints++
	}

	for _, entry := range totals {
		if entry.Datapoints > 0 {
			entry.AvgTemp /= float64(entry.Datapoints)
			entry.AvgWindSpeed /= float64(entry.Datapoints)

		}
	}
	return totals
}

func (mAPI *MetOfficeAPI) FindWindows(data []TimeSeriesData) [][]time.Time {

	//probably should be returning and handling our errors in this

	var startOfWindow string
	var endOfWindow string

	windows := [][]time.Time{}
	for i := 0; i < len(data); i++ {

		currForecast := data[i]

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

		if i == (len(data) - 1) {
			if currForecast.TotalPrecipAmount != 0.00 && startOfWindow != "" {
				endOfWindow = data[i-1].Time
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
				endOfWindow = data[i-1].Time
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
