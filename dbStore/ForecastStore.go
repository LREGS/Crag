package dbStore

//we want to find a way of cleaning the database every x amount of days to remove old forecasts

import (
	"workspaces/github.com/lregs/Crag/data"
)

// in case I forget this will be part of the Forecast service struct, that provides services to the forecast data mode
type SqlForecastStore struct {
	*SqlStore
}

func newSqlForecastStore(sqlStore *SqlStore) ForecastStore {
	fs := &SqlForecastStore{
		SqlStore: sqlStore,
	}
	return fs
}

func (fs SqlForecastStore) StoreForecast(forecast data.Forecast) error {
	DBforecast, err := fs.MarshalForecastToDB(forecast)
	if err != nil {
		return err
	}

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
	err = fs.SqlStore.masterX.QueryRow(
		query,
		DBforecast.Time,
		DBforecast.ScreenTemperature,
		DBforecast.FeelsLikeTemp,
		DBforecast.WindSpeed,
		DBforecast.WindDirection,
		DBforecast.TotalPrecipAmount,
		DBforecast.ProbOfPrecipitation,
		DBforecast.Latitude,
		DBforecast.Longitude).Scan(&DBforecast.Id)

	return nil

}
func (fs SqlForecastStore) GetForecastByID(Id int) (forecast data.DBForecast, err error) {
	forecast = data.DBForecast{}

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

	err = fs.SqlStore.masterX.QueryRow(query, Id).Scan(
		&forecast.Time,
		&forecast.ScreenTemperature,
		&forecast.FeelsLikeTemp,
		&forecast.WindSpeed,
		&forecast.WindDirection,
		&forecast.TotalPrecipAmount,
		&forecast.ProbOfPrecipitation,
		&forecast.Latitude,
		&forecast.Longitude)

	return forecast, nil

}

func (fs SqlForecastStore) GetMultipleForecastByID(Ids []int) (map[int][]interface{}, error) {
	r := make(map[int][]interface{})

	for id := range Ids {
		forecast, err := fs.GetForecastByID(id)
		if err != nil {
			return nil, err
		}
		r[id] = append(r[id], forecast)
	}
	return r, nil

}

func (fs SqlForecastStore) UpdateForecast(forecast data.DBForecast) error {

	query := `
	update forecast 
	set
	Time = $1,
	ScreenTemperature = $2,
	FeelsLikeTemp = $3,
	WindSpeed = $4,
	WindDirection = $5,
	totalPrecitipitation = $6,
	ProbofPrecipitation = $7,
	Latitude = $8,
	Longitude = $9
	where id = $10"
	)
`

	_, err := fs.SqlStore.masterX.Exec(query, forecast.Time, forecast.ScreenTemperature, forecast.FeelsLikeTemp, forecast.WindSpeed, forecast.WindDirection, forecast.TotalPrecipAmount, forecast.ProbOfPrecipitation, forecast.Latitude, forecast.Longitude, forecast.Id)
	if err != nil {
		return err
	}
	//we probably want to return success or use this logger thing -  its use the error thing :)
	return nil
}

func (fs SqlForecastStore) DeleteForecast(Id int) error {
	_, err := fs.masterX.Query("delete from forecast where id = $1", Id)
	if err != nil {
		return err
	}
	return nil
}
func (fs SqlForecastStore) GetForecastByDate(time string, id int) (err error, forecast data.DBForecast) {

	forecast = data.DBForecast{Id: id}

	query := `
		select, 
		ScreenTemperature,
		FeelsLikeTemp,
		WindSpeed,
		WindDirection,
		totalPrecitipitation,
		ProbofPrecipitation,
		Latitude,
		Longitude
		from forecast where id = $1 AND time = $2"
		)
	`

	err = fs.SqlStore.masterX.QueryRow(query, id, time).Scan(
		&forecast.Time,
		&forecast.ScreenTemperature,
		&forecast.FeelsLikeTemp,
		&forecast.WindSpeed,
		&forecast.WindDirection,
		&forecast.TotalPrecipAmount,
		&forecast.ProbOfPrecipitation,
		&forecast.Latitude,
		&forecast.Longitude)

	return nil, forecast
}

func (fs SqlForecastStore) MarshalForecastToDB(forecast data.Forecast) (data.DBForecast, error) {
	//why error if im not handling case where there is an error?
	Features := forecast.Features[0]
	Longitude := Features.Coordinates[0]
	Latitude := Features.Coordinates[1]

	TimeSeries := Features.Properties.TimeSeries[0]

	forecastDB := data.DBForecast{
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

func (fs SqlForecastStore) DryestForecast() (err error, DryestCrags []data.DBForecast) {

	query := `
	select * frome forecast where TotalPrecipAmount = (SELECT MIN(TotalPrecipAmount) from forecast)
`

	var DryestForecast []data.DBForecast

	rows, err := fs.SqlStore.masterX.Query(query)

	if err != nil {
		return err, nil
	}

	for rows.Next() {
		f := data.DBForecast{}
		err := rows.Scan(
			&f.Time,
			&f.ScreenTemperature,
			&f.FeelsLikeTemp,
			&f.WindSpeed,
			&f.WindDirection,
			&f.TotalPrecipAmount,
			&f.ProbOfPrecipitation,
			&f.Latitude,
			&f.Longitude)

		if err != nil {
			return err, nil
		}

		DryestForecast = append(DryestForecast, f)

	}

	return nil, DryestForecast

}

// func (fs SqlForecastStore) GetForecastByValues(Id int, values []string) (data.DBForecast, error){

// 	// not sure if I should use naked interfaces here because of the lack of type safety
// 	var returnedValues interface{}

// 	forecastByValue := make(map[int]interface{})
// 	query := "select"

// 	for _, value := range values {
// 		if value == values[(len(values) - 1)]{
// 			query += " " + value
// 		}
// 		query += " " + value + ","
// 	}

// 	query += " where id = $1"

// 	_, err := fs.SqlStore.masterX.Exec(query, Id)
// 	if err != nil{
// 		return nil, err
// 	}

// }
