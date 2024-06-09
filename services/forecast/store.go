package forecast

//This file will take the data from the met office api, extract the data it needs and store it within redis

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

type MetPayload struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

// Selects the data we want from the response

type Forecast struct {
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

// func redisClient() *redis.Client {

// 	rc := redis.NewClient(&redis.Options{
// 		Addr:     "redis-19441.c233.eu-west-1-1.ec2.redns.redis-cloud.com:19441",
// 		Password: "",
// 		DB:       0,
// 	})

// 	c := cache.New(&cache.Options{
// 		Redis: rc,
// 	})
// 	return rc

// }

func totalValues([]Forecast) {

}
