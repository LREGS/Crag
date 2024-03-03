package main

type Crag struct {
	Id        int
	Name      string
	Latitude  float64
	Longitude float64
	Climbs    []Climb      //many to one relationship not db field
	Reports   []Report     //many to one relationship not db field
	Forecast  []DBForecast //I dont think I have an int key in my forecast pointing to crag and this needs fixingf
}

type Climb struct {
	Id    int
	Name  string
	Grade string
	Crag  *Crag
}

type Report struct {
	Id      int
	Content string
	Author  string
	Crag    *Crag
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
	Type        string `json:"type"`
	Geometry    Geometry
	Properties  Properties
	Coordinates []float64 `json:"coordinates"`
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
	Id                  int
	Time                string
	ScreenTemperature   float64
	FeelsLikeTemp       float64
	WindSpeed           float64
	WindDirection       float64
	TotalPrecipAmount   float64
	ProbOfPrecipitation float64
	Latitude            float64
	Longitude           float64
}
