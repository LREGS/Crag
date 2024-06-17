package models

type Crag struct {
	Id        int     `json:"id"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type CragPayload struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Climb struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Grade  string `json:"grade"`
	CragID int    `json:"cragId"`
}

type ClimbPayload struct {
	Name   string `json:"name"`
	Grade  string `json:"grade"`
	CragID int    `json:"cragId"`
}

type Report struct {
	Id      int    `json:"id"`
	Content string `json:"content"`
	Author  string `json:"author"`
	Crag    *Crag  `json:"crag"`
}

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

// Selects the data we want from the response

type DBForecast struct {
	Id                  int     `json:"id"`
	Time                string  `json:"time"`
	ScreenTemperature   float64 `json:"screenTemperature"`
	FeelsLikeTemp       float64 `json:"feelsLikeTemp"`
	WindSpeed           float64 `json:"windSpeed"`
	WindDirection       float64 `json:"windDirection"`
	TotalPrecipAmount   float64 `json:"totalPrecipAmount"`
	ProbOfPrecipitation int     `json:"probOfPrecipitation"`
	Latitude            float64 `json:"latitude"`
	Longitude           float64 `json:"longitude"`
}

type DBForecastPayload struct {
	Time                string  `json:"time"`
	ScreenTemperature   float64 `json:"screenTemperature"`
	FeelsLikeTemp       float64 `json:"feelsLikeTemp"`
	WindSpeed           float64 `json:"windSpeed"`
	WindDirection       int     `json:"windDirection"`
	TotalPrecipAmount   float64 `json:"totalPrecipAmount"`
	ProbOfPrecipitation int     `json:"probOfPrecipitation"`
	Latitude            float64 `json:"latitude"`
	Longitude           float64 `json:"longitude"`
	CragId              int     `json:"cragId"`
}
