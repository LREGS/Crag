package main

import "time"

//Model the incoming data from the met office hourly api

type TimeSeriesData struct {
	Time                      string  `json:"time"`
	ScreenTemperature         float64 `json:"screenTemperature"`
	MaxScreenAirTemp          float64 `json:"maxScreenAirTemp"`
	MinScreenAirTemp          float64 `json:"minScreenAirTemp"`
	ScreenDewPointTemperature float64 `json:"screenDewPointTemperature"`
	FeelsLikeTemperature      float64 `json:"feelsLikeTemperature"`
	WindSpeed10m              float64 `json:"windSpeed10m"`
	WindDirectionFrom10m      int     `json:"windDirectionFrom10m"`
	WindGustSpeed10m          float64 `json:"windGustSpeed10m"`
	Max10mWindGust            float64 `json:"max10mWindGust"`
	Visibility                int     `json:"visibility"`
	ScreenRelativeHumidity    float64 `json:"screenRelativeHumidity"`
	Mslp                      int     `json:"mslp"`
	UvIndex                   int     `json:"uvIndex"`
	SignificantWeatherCode    int     `json:"significantWeatherCode"`
	PrecipitationRate         float64 `json:"precipitationRate"`
	TotalPrecipAmount         float64 `json:"totalPrecipAmount"`
	TotalSnowAmount           float64 `json:"totalSnowAmount"`
	ProbOfPrecipitation       int     `json:"probOfPrecipitation"`
}

type Feature struct {
	Type        string     `json:"type"`
	Geometry    Geometry   `json:"geometry"`
	Properties  Properties `json:"properties"`
	Coordinates []float64  `json:"coordinates"`
}

type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type Properties struct {
	RequestPointDistance float64          `json:"requestPointDistance"`
	ModelRunDate         string           `json:"modelRunDate"`
	TimeSeries           []TimeSeriesData `json:"timeSeries"`
}

type Forecast struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

type ForecastTotals struct {
	HighestTemp  float64
	LowestTemp   float64
	AvgTemp      float64
	AvgWindSpeed float64
	TotalPrecip  float64
	Datapoints   int
}

// A window is a gap in precipitation of more than 1hour and gives the amount of rain the area has just experienced.

type Window struct {
	Time                       []string
	AvgTemp                    float64
	AvgWindSpeed               float64
	WindDirection              int
	PrecipInLastXAmountofHours int
}

// Met office weather model is updated and released every hour for the next 72 hours.

type ForecastPayload struct {
	LastModelRunTime string
	ForecastTotals   map[string]*ForecastTotals
	Windows          [][]time.Time
}

type Crag struct {
	Name      string
	Latitude  float64
	Longitude float64
}
