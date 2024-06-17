package store

import (
	"context"
	"errors"
	"log"
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

func (fs *SqlForecastStore) StoreForecast(ctx context.Context, forecast models.DBForecastPayload) (models.DBForecast, error) {

	var storedForecast models.DBForecast

	err := fs.validatePayload(forecast)
	if err != nil {
		return storedForecast, err
	}

	err = fs.Store.masterX.QueryRow(
		ctx,
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
		&storedForecast.Longitude)
	if err != nil {
		return storedForecast, err
	}

	return storedForecast, nil
}

const getForecastByCrag = `select * from forecast where Id = $1`

func (fs *SqlForecastStore) GetForecastByCragId(ctx context.Context, CragId int) ([]models.DBForecast, error) {
	//we're returning every forecast, need some function/ http endpoint that will serve
	// presented data from the forecast (total rainfall etc)
	rows, err := fs.Store.masterX.Query(ctx, getForecastByCrag, CragId)
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
			&forecast.Longitude)
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

func (fs *SqlForecastStore) GetAllForecastsByCragId(ctx context.Context) (map[int][]models.DBForecast, error) {

	//this is returning every forecast for every crag we have, not every forecast based on the crag Id

	rows, err := fs.Store.masterX.Query(ctx, getAllForecast)
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
			&forecast.Longitude)
		if err != nil {
			return nil, err
		}

		//please change this back

		// results[forecast.CragId] = append(results[forecast.CragId], forecast)

	}
	return results, nil
}

const deleteForecastById = `DELETE FROM forecast where Id = $1 returning *`

// I should be returning an instance of the deleted data
func (fs *SqlForecastStore) DeleteForecastById(ctx context.Context, Id int) (models.DBForecast, error) {

	var forecast models.DBForecast
	if err := fs.Store.masterX.QueryRow(ctx, deleteForecastById, Id).Scan(
		&forecast.Id,
		&forecast.Time,
		&forecast.ScreenTemperature,
		&forecast.FeelsLikeTemp,
		&forecast.WindSpeed,
		&forecast.WindDirection,
		&forecast.TotalPrecipAmount,
		&forecast.ProbOfPrecipitation,
		&forecast.Latitude,
		&forecast.Longitude); err != nil {
		return models.DBForecast{}, err
	}
	//in this stage do we want to be validating here or does it go back through validation middleware I dont know?! Maybe middleware is only validating data from the client
	return forecast, nil
}

func (fs *SqlForecastStore) validatePayload(data models.DBForecastPayload) error {
	if reflect.DeepEqual(models.DBForecastPayload{}, data) {
		return errors.New("input cannot be empty")
	}
	return nil
}

//	func (fs *SqlForecastStore) validateDBForecast(data models.DBForecast) error {
//		if reflect.DeepEqual(models.DBForecast{}, data) {
//			return errors.New("db value returned empty")
//		}
//		return nil

// const copyCSV = `COPY forecast FROM STDIN WITH CSV HEADER`

func (fs *SqlForecastStore) Populate(ctx context.Context, log *log.Logger) {

	// payload, _ := met.GetPayload(log, []float64{53.12000233374393, -4.000659549362343})

	// _, err := fs.Store.masterX.CopyFrom(
	// 	ctx,
	// 	pgx.Identifier{"forecast"},
	// 	[]string{
	// 		"id",
	// 		"time",
	// 		"screentemperature",
	// 		"feelsliketemp",
	// 		"windspeed",
	// 		"winddirection",
	// 		"totalprecipitation",
	// 		"probofprecipitation",
	// 		"latitude",
	// 		"longitude"},
	// 	pgx.CopyFromRows(payload),
	// )
	// if err != nil {
	// 	log.Printf("failed to populate db %s", err)
	// }

}

const drop = `DROP TABLE forecast`

const createTable = `CREATE TABLE forecast (
    Id Int, 
    Time VARCHAR(255) UNIQUE,
    ScreenTemperature DOUBLE PRECISION,
    FeelsLikeTemp DOUBLE PRECISION, 
    WindSpeed DOUBLE PRECISION,
    WindDirection DOUBLE PRECISION,
    totalPrecipitation DOUBLE PRECISION,
    ProbofPrecipitation INT,
    Latitude DOUBLE PRECISION,
    Longitude DOUBLE PRECISION
);`

func (fs *SqlForecastStore) Refresh(ctx context.Context, log *log.Logger) {

	log.Print("Deleting data from forecast")
	_, err := fs.Store.masterX.Exec(ctx, drop)
	if err != nil {
		log.Printf("failed dropping %s", err)
	}

	log.Print("Creating tables")
	_, err = fs.Store.masterX.Exec(ctx, createTable)
	if err != nil {
		log.Printf("failed creating table %s", err)
	}

	fs.Populate(ctx, log)

}

// func forecast2csv(log *log.Logger, f models.Forecast) *os.File {

// 	// fmt.Println(f)

// 	d := f.Features[0].Properties.TimeSeries

// 	result := make([][]string, len(d))

// 	//header
// 	result[0] = []string{"Id", "Time", "ScreenTemperature", "FeelsLikeTemp", "WindSpeed",
// 		"WindDirection", "totalPrecipitation", "ProbofPrecipitation", "Latitude", "Longitude"}

// 	for i := 1; i < len(d); i++ {
// 		result[i] = []string{
// 			strconv.FormatFloat(d[i].FeelsLikeTemperature, 'f', -1, 64),
// 			strconv.FormatFloat(d[i].WindSpeed10m, 'f', -1, 64),
// 			strconv.Itoa(d[i].WindDirectionFrom10m),
// 			strconv.FormatFloat(d[i].TotalPrecipAmount, 'f', -1, 64),
// 			strconv.Itoa(d[i].ProbOfPrecipitation),
// 			strconv.FormatFloat(f.Features[0].Geometry.Coordinates[0], 'f', -1, 64),
// 			strconv.FormatFloat(f.Features[0].Geometry.Coordinates[1], 'f', -1, 64),
// 		}
// 	}

// 	file, err := os.Create("forecast.csv")
// 	if err != nil {
// 		log.Printf("failed creating file %s", err)
// 	}
// 	w := csv.NewWriter(file)
// 	w.WriteAll(result)

// 	return file

// }
