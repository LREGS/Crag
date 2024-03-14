package store

import (
	"testing"
	"time"

	"github.com/lregs/Crag/models"
)

func TestAddForecast(t *testing.T) {
	//probably need more testing to make sure it saved correctly
	MockStore := returnPrePopulatedMockStore(t, true, false)
	forecast := newForecast()

	t.Run("Testing Add Forecast", func(t *testing.T) {
		_, err := MockStore.Stores.ForecastStore.AddForecast(forecast)

		if err != nil {
			t.Fatalf("post forecast request failed because of err: %s", err)
		}

	})

}

func TestGetForecastByCrag(t *testing.T) {
	MockStore := returnPrePopulatedMockStore(t, true, true)

	t.Run("Testing get by CragId", func(t *testing.T) {
		const Id = 1
		results, err := MockStore.Stores.ForecastStore.GetForecastByCragId(Id)
		if err != nil {
			t.Fatalf("Could not perform sql task because of this error: %s", err)
		}
		if len(resluts) > 0 != true {
			t.Fatalf("No forecasts were returned")
		}

	})
}

func newForecast() models.DBForecast {
	forecast := models.DBForecast{
		Id:                  1,
		Time:                time.Now().Format(time.RFC3339),
		ScreenTemperature:   25.0,
		FeelsLikeTemp:       24.0,
		WindSpeed:           10.0,
		WindDirection:       180.0,
		TotalPrecipAmount:   0.0,
		ProbOfPrecipitation: 0.0,
		Latitude:            51.5074,
		Longitude:           -0.1278,
		CragId:              1,
	}

	return forecast
}
