package store

import (
	"errors"
	"reflect"

	"github.com/lregs/Crag/models"
)

type SqlForecastStore struct {
	Store *SqlStore
}

func NewForecastStore(sqlStore *SqlStore) *SqlForecastStore {
	store := &SqlForecastStore{sqlStore}
	return store
}

const storeForecast = `insert into forecast(	
	Time, 
	ScreenTemperature,
	FeelsLikeTemp,
	WindSpeed, 
	WindDirection,
	totalPrecipitation, 
	ProbofPrecipitation, 
	Latitude, 
	Longitude, 
	CragID
	)values(
	$1,$2,$3,$4,$5,$6,$7,$8,$9,$10
	) RETURNING *`

func (fs *SqlForecastStore) StoreForecast(forecast models.DBForecastPayload) (models.DBForecast, error) {

	var storedForecast models.DBForecast

	err := fs.validatePayload(forecast)
	if err != nil {
		return storedForecast, err
	}

	err = fs.Store.masterX.QueryRow(
		storeForecast,
		forecast.Time,
		forecast.ScreenTemperature,
		forecast.FeelsLikeTemp,
		forecast.WindSpeed,
		forecast.WindDirection,
		forecast.TotalPrecipAmount,
		forecast.ProbOfPrecipitation,
		forecast.Latitude,
		forecast.Longitude,
		forecast.CragId).Scan(
		&storedForecast.Id,
		&storedForecast.Time,
		&storedForecast.ScreenTemperature,
		&storedForecast.FeelsLikeTemp,
		&storedForecast.WindSpeed,
		&storedForecast.WindDirection,
		&storedForecast.TotalPrecipAmount,
		&storedForecast.ProbOfPrecipitation,
		&storedForecast.Latitude,
		&storedForecast.Longitude,
		&storedForecast.CragId)
	if err != nil {
		return storedForecast, err
	}

	return storedForecast, nil
}

const getForecastByCrag = `select * from forecast where CragId = $1`

func (fs *SqlForecastStore) GetForecastByCragId(CragId int) ([]models.DBForecast, error) {
	//we're returning every forecast, need some function/ http endpoint that will serve
	// presented data from the forecast (total rainfall etc)
	rows, err := fs.Store.masterX.Query(getForecastByCrag, CragId)
	if err != nil {
		return nil, err
	}

	var results []models.DBForecast

	for rows.Next() {
		var forecast models.DBForecast
		err := rows.Scan(
			&forecast.Id,
			&forecast.Time,
			&forecast.ScreenTemperature,
			&forecast.FeelsLikeTemp,
			&forecast.WindSpeed,
			&forecast.WindDirection,
			&forecast.TotalPrecipAmount,
			&forecast.ProbOfPrecipitation,
			&forecast.Latitude,
			&forecast.Longitude,
			&forecast.CragId)
		if err != nil {
			return nil, err
		}
		results = append(results, forecast)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

const getAllForecast = `select * from forecast`

func (fs *SqlForecastStore) GetAllForecastsByCragId() (map[int][]models.DBForecast, error) {

	//this is returning every forecast for every crag we have, not every forecast based on the crag Id

	rows, err := fs.Store.masterX.Query(getAllForecast)
	if err != nil {
		return nil, err
	}

	results := make(map[int][]models.DBForecast)

	for rows.Next() {
		var forecast models.DBForecast
		err := rows.Scan(
			&forecast.Id,
			&forecast.Time,
			&forecast.ScreenTemperature,
			&forecast.FeelsLikeTemp,
			&forecast.WindSpeed,
			&forecast.WindDirection,
			&forecast.TotalPrecipAmount,
			&forecast.ProbOfPrecipitation,
			&forecast.Latitude,
			&forecast.Longitude,
			&forecast.CragId)
		if err != nil {
			return nil, err
		}
		results[forecast.CragId] = append(results[forecast.CragId], forecast)

	}
	return results, nil
}

const deleteForecastById = `DELETE FROM forecast where Id = $1`

// I should be returning an instance of the deleted data
func (fs *SqlForecastStore) DeleteForecastById(Id int) error {
	_, err := fs.Store.masterX.Exec(deleteForecastById, Id)
	if err != nil {
		return err
	}
	return nil
}

func (fs *SqlForecastStore) validatePayload(data models.DBForecastPayload) error {
	if reflect.DeepEqual(models.DBForecastPayload{}, data) {
		return errors.New("input cannot be empty")
	}
	return nil
}

func (fs *SqlForecastStore) validateDBForecast(data models.DBForecast) error {
	if reflect.DeepEqual(models.DBForecast{}, data) {
		return errors.New("db value returned empty")
	}
	return nil
}
