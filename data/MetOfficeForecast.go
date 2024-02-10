package data

import (
	"database/sql"
	"strings"

	"strconv"

	helpers "workspaces/github.com/lregs/Crag/helper"

	_ "github.com/lib/pq"
)

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

func MarshalForecastToDB(cragID int, forecast Forecast) (DBForecast, error) {

	Features := forecast.Features[0]
	Longitude := Features.Coordinates[0]
	Latitude := Features.Coordinates[1]

	TimeSeries := Features.Properties.TimeSeries[0]

	forecastDB := DBForecast{
		Time:                TimeSeries.Time,
		ScreenTemperature:   TimeSeries.ScreenTemperature,
		FeelsLikeTemp:       TimeSeries.FeelsLikeTemperature,
		WindSpeed:           TimeSeries.WindSpeed10m,
		WindDirection:       float64(TimeSeries.WindDirectionFrom10m),
		TotalPrecipAmount:   TimeSeries.TotalPrecipAmount,
		ProbOfPrecipitation: TimeSeries.TotalPrecipAmount,
		Latitude:            Latitude,
		Longitude:           Longitude,
	}

	return forecastDB, nil

}

func (forecast *DBForecast) Create(db *sql.DB) (err error) {
	query := `
		insert into forecast(
		Time,
		ScreenTemperature,
		FeelsLikeTemp,
		WindSpeed,
		WindDirection,
		totalPrecitipitation,
		ProbofPrecipitation,
		Latitude,
		Longitude
		)
		values(
		$1,$2,$3,$4,$5,
		$6,$7,$8,$9,$10
		)
		returning id	
	`
	err = db.QueryRow(query, forecast.Time, forecast.ScreenTemperature,
		forecast.FeelsLikeTemp, forecast.WindSpeed, forecast.WindDirection,
		forecast.TotalPrecipAmount, forecast.ProbOfPrecipitation,
		forecast.Latitude, forecast.Longitude).Scan(&forecast.Id)

	return nil
}

func GetForecast(Id int, db *sql.DB) (forecast DBForecast, err error) {
	forecast = DBForecast{}

	query := `
		select Id, 
		Time,
		ScreenTemperature,
		FeelsLikeTemp,
		WindSpeed,
		WindDirection,
		totalPrecitipitation,
		ProbofPrecipitation,
		Latitude,
		Longitude
		from forecast where id = $1"
		)
	`

	err = db.QueryRow(query, Id).Scan(&forecast.Time, &forecast.ScreenTemperature,
		&forecast.FeelsLikeTemp, &forecast.WindSpeed, &forecast.WindDirection,
		&forecast.TotalPrecipAmount, &forecast.ProbOfPrecipitation,
		&forecast.Latitude, &forecast.Longitude)

	return forecast, nil

}

func (forecast *DBForecast) UpdateForecast(Id int, db *sql.DB, updates map[string]interface{}) (err error) {

	query := "update forcast set "

	params := make([]interface{}, 0)
	i := 1

	for key, value := range updates {
		query += key + " = $" + strconv.Itoa(i) + ","
		params = append(params, value)
	}

	query = strings.TrimSuffix(query, ",")

	query += " Where id = $1"

	//not sure this works but lets see
	_, err = db.Exec(query, append(params, Id)...)
	helpers.CheckError(err)

	return nil

}

func (forecast *DBForecast) DeleteForecast(Id int, db *sql.DB) (err error) {
	_, err = db.Exec("delete from forecast where id = $1", forecast.Id)
	return
}
