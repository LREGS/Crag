package store

import (
	"context"
	"testing"

	"github.com/lregs/Crag/models"
	"github.com/stretchr/testify/assert"
)

func TestStoreForecast(t *testing.T) {
	//probably need more testing to make sure it saved correctly
	MockStore := returnPrePopulatedMockStore(t, true, false)
	forecast := testPayload()

	t.Run("Testing Store Forecast", func(t *testing.T) {
		storedForecast, err := MockStore.Stores.ForecastStore.StoreForecast(context.Background(), forecast)

		if err != nil {
			t.Fatalf("post forecast request failed because of err: %s", err)
		}

		assert.Equal(t, testForecst(), storedForecast)

	})

	t.Run("testing invalid forecast", func(t *testing.T) {
		_, err := MockStore.Stores.ForecastStore.StoreForecast(context.Background(), models.DBForecastPayload{})
		if err == nil {
			t.Fatal("stored empty values")
		}
	})

}

func TestGetForecastByCrag(t *testing.T) {
	MockStore := returnPrePopulatedMockStore(t, true, true)

	t.Run("Testing get by CragId", func(t *testing.T) {
		const Id = 1
		results, err := MockStore.Stores.ForecastStore.GetForecastByCragId(context.Background(), Id)
		if err != nil {
			t.Fatalf("Could not perform sql task because of this error: %s", err)
		}
		if len(results) > 0 != true {
			t.Fatalf("No forecasts were returned")
		}

	})
}

func TestGetAllForecasts(t *testing.T) {
	MockStore := returnPrePopulatedMockStore(t, true, true)

	t.Run("Testing get all forecasts", func(t *testing.T) {
		results, err := MockStore.Stores.ForecastStore.GetAllForecastsByCragId(context.Background())
		if err != nil {
			t.Fatalf("could not get forecasts because of error: %s", err)
		}

		if len(results) > 0 != true {
			t.Fatalf("No forecasts were returned")
		}
		t.Logf("crag %+v was returned", results[0])
	})
}

func TestDeleteForecast(t *testing.T) {
	MockStore := returnPrePopulatedMockStore(t, true, true)

	t.Run("Testing delete forecast", func(t *testing.T) {

		const Id = 1
		deletedData, err := MockStore.Stores.ForecastStore.DeleteForecastById(context.Background(), Id)
		if err != nil {
			t.Fatalf("could not delete item becasue of err: %s", err)
		}

		assert.Equal(t, testForecst(), deletedData)

	})
}

func testPayload() models.DBForecastPayload {
	forecast := models.DBForecastPayload{
		Time:                "11",
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

func testForecst() models.DBForecast {
	forecast := models.DBForecast{
		Id:                  1,
		Time:                "11",
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
