package forecast

import (
	"errors"
	"testing"

	"github.com/lregs/Crag/models"
)

type MockForecastStore struct {
	forecast []*models.DBForecast
}

func (fs *MockForecastStore) returnDBForecast(p *models.DBForecastPayload, Id int) *models.DBForecast {
	var storedForecast models.DBForecast
	storedForecast.Id = Id
	storedForecast.Time = p.Time
	storedForecast.ScreenTemperature = p.ScreenTemperature
	storedForecast.FeelsLikeTemp = p.FeelsLikeTemp
	storedForecast.WindSpeed = p.WindSpeed
	storedForecast.WindDirection = p.WindDirection
	storedForecast.TotalPrecipAmount = p.TotalPrecipAmount
	storedForecast.ProbOfPrecipitation = p.ProbOfPrecipitation
	storedForecast.Latitude = p.Latitude
	storedForecast.Longitude = p.Longitude
	storedForecast.CragId = p.CragId

	return &storedForecast
}

func (fs *MockForecastStore) AddForecast(newForecast *models.DBForecastPayload) (*models.DBForecast, error) {

	fToStore := fs.returnDBForecast(newForecast, (len(fs.forecast) + 1))

	fs.forecast = append(fs.forecast, fToStore)

	return fs.forecast[(len(fs.forecast) - 1)], nil
}

func (fs *MockForecastStore) GetForecastByCragId(CragId int) ([]models.DBForecast, error) {
	if CragId == 0 {
		return nil, errors.New("id 0 not valid")
	}

	res := []models.DBForecast{}

	for _, f := range fs.forecast {
		if f.CragId == CragId {
			res = append(res, *f)
		}
	}

	if len(res) == 0 {
		return res, errors.New("no forecast where found for crag Id")
	} else {
		return res, nil
	}
}

func (fs *MockForecastStore) GetAllForecastsByCragId() (map[int][]models.DBForecast, error) {

	if len(fs.forecast) == 0 {
		return nil, errors.New("DB is empty")
	} else {
		res := make(map[int][]models.DBForecast, 0)

		for i, f := range fs.forecast {
			res[i] = append(res[i], *f)
		}
		return res, nil
	}
}
func (fs *MockForecastStore) DeleteForecastById(Id int) error {
	if len(fs.forecast) == 0 {
		return errors.New("No forecasts to delete")
	}

	if Id > len(fs.forecast)+1 {
		return errors.New("forecast doesn't exist to delete")
	}

	return nil
}

func TestAddForecast(t *testing.T) {

	store := &MockForecastStore{
		forecast: []*models.DBForecast{
			{
				Id:                  1,
				Time:                "2024-04-06T12:00:00Z",
				ScreenTemperature:   20.5,
				FeelsLikeTemp:       18.2,
				WindSpeed:           10.0,
				WindDirection:       180.0,
				TotalPrecipAmount:   0.5,
				ProbOfPrecipitation: 30.0,
				Latitude:            40.01,
				Longitude:           40.11,
				CragId:              1,
			},
			{
				Id:                  2,
				Time:                "2024-04-06T13:00:00Z",
				ScreenTemperature:   22.3,
				FeelsLikeTemp:       20.1,
				WindSpeed:           12.5,
				WindDirection:       200.0,
				TotalPrecipAmount:   0.8,
				ProbOfPrecipitation: 40.0,
				Latitude:            41.01,
				Longitude:           41.11,
				CragId:              2,
			},
		},
	}

}
