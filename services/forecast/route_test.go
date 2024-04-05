package forecast

import "github.com/lregs/Crag/models"

type MockForecastStore struct {
	forecast []*models.DBForecastPayload
}

func (fs *MockForecastStore) AddForecast(newForecast *models.DBForecastPayload) (*models.DBForecast, error) {
	fs.forecast = append(fs.forecast, newForecast)

	last := len(fs.forecast) - 1
	var storedForecast models.DBForecast
	storedForecast.Id = len(fs.forecast)
	storedForecast.Time = fs.forecast[last].Time
	storedForecast.ScreenTemperature = fs.forecast[last].ScreenTemperature
	storedForecast.FeelsLikeTemp = fs.forecast[last].FeelsLikeTemp
	storedForecast.WindSpeed = fs.forecast[last].WindSpeed
	storedForecast.WindDirection = fs.forecast[last].WindDirection
	storedForecast.TotalPrecipAmount = fs.forecast[last].TotalPrecipAmount
	storedForecast.ProbOfPrecipitation = fs.forecast[last].ProbOfPrecipitation
	storedForecast.Latitude = fs.forecast[last].Latitude
	storedForecast.Longitude - fs.forecast[last].Longitude
	storedForecast.CragId = fs.forecast[last].CragId

	return nil, nil
}

func (fs *MockForecastStore) GetForecastByCragId(CragId int) ([]models.DBForecast, error) {
	return nil, nil
}
func (fs *MockForecastStore) GetAllForecasts() (map[int][]models.DBForecast, error) { return nil, nil }
func (fs *MockForecastStore) DeleteForecastById(*models.DBForecast, error) error    { return nil, nil }
