package forecast

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/lregs/Crag/models"
	"github.com/lregs/Crag/util"
	"github.com/stretchr/testify/assert"
)

type MockForecastStore struct {
	forecast []models.DBForecast
}

func (fs *MockForecastStore) Validate(*models.DBForecast) error {
	return nil
}

func (fs *MockForecastStore) returnDBForecast(p *models.DBForecastPayload, Id int) models.DBForecast {
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

	return storedForecast
}

func (fs *MockForecastStore) StoreForecast(newForecast models.DBForecastPayload) (models.DBForecast, error) {

	if newForecast.Time == "" {
		f := models.DBForecast{}
		return f, errors.New("no empty values allowed in forecast entry to db")
	}

	fToStore := fs.returnDBForecast(&newForecast, (len(fs.forecast) + 1))

	// if fToStore.Time == "" {
	// 	return fToStore, errors.New("invalid data")
	// }

	fs.forecast = append(fs.forecast, fToStore)

	f := models.DBForecast{
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
	}

	return f, nil
}

func (fs *MockForecastStore) GetForecastByCragId(CragId int) ([]models.DBForecast, error) {

	return []models.DBForecast{
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
			CragId:              2,
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
		{
			Id:                  3,
			Time:                "2024-04-06T13:00:00Z",
			ScreenTemperature:   22.3,
			FeelsLikeTemp:       20.1,
			WindSpeed:           12.5,
			WindDirection:       200.0,
			TotalPrecipAmount:   0.8,
			ProbOfPrecipitation: 40.0,
			Latitude:            41.01,
			Longitude:           41.11,
			CragId:              3,
		},
	}, nil

}

func (fs *MockForecastStore) GetAllForecastsByCragId() (map[int][]models.DBForecast, error) {
	data := map[int][]models.DBForecast{
		2: {
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
				CragId:              2,
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
		3: {
			{
				Id:                  3,
				Time:                "2024-04-06T13:00:00Z",
				ScreenTemperature:   22.3,
				FeelsLikeTemp:       20.1,
				WindSpeed:           12.5,
				WindDirection:       200.0,
				TotalPrecipAmount:   0.8,
				ProbOfPrecipitation: 40.0,
				Latitude:            41.01,
				Longitude:           41.11,
				CragId:              3,
			},
		},
	}

	return data, nil
}
func (fs *MockForecastStore) DeleteForecastById(Id int) (models.DBForecast, error) {
	if len(fs.forecast) == 0 {
		return models.DBForecast{}, errors.New("No forecasts to delete")
	}

	if Id > len(fs.forecast)+1 {
		return models.DBForecast{}, errors.New("forecast doesn't exist to delete")
	}

	return models.DBForecast{Id: 1,
		Time:                "2024-04-06T12:00:00Z",
		ScreenTemperature:   20.5,
		FeelsLikeTemp:       18.2,
		WindSpeed:           10.0,
		WindDirection:       180.0,
		TotalPrecipAmount:   0.5,
		ProbOfPrecipitation: 30.0,
		Latitude:            40.01,
		Longitude:           40.11,
		CragId:              1}, nil
}

func TestPost(t *testing.T) {

	store := &MockForecastStore{
		forecast: []models.DBForecast{
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

	handler := NewHandler(store)
	router := mux.NewRouter()

	router.PathPrefix("/forecast").HandlerFunc(handler.Post()).Methods("POST")

	t.Run("Testing Valid Post", func(t *testing.T) {

		payload := models.DBForecastPayload{
			Time:                "2024-04-06T13:00:00Z",
			ScreenTemperature:   22.3,
			FeelsLikeTemp:       20.1,
			WindSpeed:           12.5,
			WindDirection:       200.0,
			TotalPrecipAmount:   0.8,
			ProbOfPrecipitation: 40.0,
			Latitude:            41.01,
			Longitude:           41.11,
			CragId:              2}

		body, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("could not marshal body because of err: %s", err)
		}
		res, req, err := util.NewPostRequest(body, "/forecast")
		if err != nil {
			t.Fatalf("error %s making new request", err)
		}

		router.ServeHTTP(res, req)

		switch res.Code {
		case 200:
			var data models.DBForecast

			if err := json.Unmarshal(res.Body.Bytes(), &data); err != nil {
				t.Fatalf("Error decoding response: %s", err)
			}

			expected := models.DBForecast{
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
				CragId:              1}

			assert.Equal(t, expected, data)

		default:
			var errMsg map[string]string

			if err := json.Unmarshal(res.Body.Bytes(), &errMsg); err != nil {
				t.Fatalf("error decoding response into error %s", err)
			}

			t.Fatalf("unexpected code :%d, error %s", res.Code, errMsg["Error"])
		}
	})

	//t.Run("testing invalid data") - at some point we do also maybe want data validation at the store layer?

	t.Run("Testing Invalid Request Method", func(t *testing.T) {
		response := httptest.NewRecorder()
		request, err := http.NewRequest("GET", "/forecast", nil)
		if err != nil {
			t.Fatalf("creating request failed %s", err)
		}

		router.ServeHTTP(response, request)

		if response.Code != 405 {
			t.Fatalf("Accepted incorrect method, code %d", response.Code)
		}

	})

	t.Run("Testing Empty Data", func(t *testing.T) {
		payload := models.DBForecastPayload{}

		b, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("cannot unmarshal, %s", err)
		}

		res, req, err := util.NewPostRequest(b, "/forecast")
		if err != nil {
			t.Fatalf("cannot create request %s", err)
		}

		router.ServeHTTP(res, req)

	})
}

func TestGetByCragId(t *testing.T) {
	store := &MockForecastStore{
		forecast: []models.DBForecast{
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
				CragId:              2,
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
			{
				Id:                  3,
				Time:                "2024-04-06T13:00:00Z",
				ScreenTemperature:   22.3,
				FeelsLikeTemp:       20.1,
				WindSpeed:           12.5,
				WindDirection:       200.0,
				TotalPrecipAmount:   0.8,
				ProbOfPrecipitation: 40.0,
				Latitude:            41.01,
				Longitude:           41.11,
				CragId:              3,
			},
		},
	}

	handler := NewHandler(store)
	router := mux.NewRouter()

	router.PathPrefix("/forecast/{Id}").HandlerFunc(handler.GetByCragId()).Methods("GET")

	t.Run("Valid CragID", func(t *testing.T) {

		res, req := util.NewGetRequest("/forecast/2")
		router.ServeHTTP(res, req)

		switch res.Code {
		case 200:
			var forecasts []models.DBForecast
			if err := json.Unmarshal(res.Body.Bytes(), &forecasts); err != nil {
				t.Fatalf("failed decoding")
			}

			assert.Equal(t, store.forecast, forecasts)
		}

	})

	t.Run("Invalid Request Id", func(t *testing.T) {
		res, req := util.NewGetRequest("/forecast/dank")
		router.ServeHTTP(res, req)

		if res.Code != 400 {
			t.Fatal("handler accepted incorrect method")
		}

		//something in our api error response it work as expected which is annoying

	})

	t.Run("Invalid Method", func(t *testing.T) {
		res := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/forecast/1", nil)
		if err != nil {
			t.Fatal("creating request failed")
		}
		router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusMethodNotAllowed, res.Code)
	})
}

// func TestGetAllForecasts(t *testing.T){

// 	store := &MockForecastStore{
// 		forecast: []models.DBForecast{
// 			{
// 				Id:                  1,
// 				Time:                "2024-04-06T12:00:00Z",
// 				ScreenTemperature:   20.5,
// 				FeelsLikeTemp:       18.2,
// 				WindSpeed:           10.0,
// 				WindDirection:       180.0,
// 				TotalPrecipAmount:   0.5,
// 				ProbOfPrecipitation: 30.0,
// 				Latitude:            40.01,
// 				Longitude:           40.11,
// 				CragId:              2,
// 			},
// 			{
// 				Id:                  2,
// 				Time:                "2024-04-06T13:00:00Z",
// 				ScreenTemperature:   22.3,
// 				FeelsLikeTemp:       20.1,
// 				WindSpeed:           12.5,
// 				WindDirection:       200.0,
// 				TotalPrecipAmount:   0.8,
// 				ProbOfPrecipitation: 40.0,
// 				Latitude:            41.01,
// 				Longitude:           41.11,
// 				CragId:              2,
// 			},
// 			{
// 				Id:                  3,
// 				Time:                "2024-04-06T13:00:00Z",
// 				ScreenTemperature:   22.3,
// 				FeelsLikeTemp:       20.1,
// 				WindSpeed:           12.5,
// 				WindDirection:       200.0,
// 				TotalPrecipAmount:   0.8,
// 				ProbOfPrecipitation: 40.0,
// 				Latitude:            41.01,
// 				Longitude:           41.11,
// 				CragId:              3,
// 			},
// 		},
// 	}
// 	handler := NewHandler(store)
// 	router := mux.NewRouter()
//  	router.PathPrefix("/forecast/all").HandlerFunc(handler.GetAllForecasts()).Methods("GET")

// }

func TestGetAllForecasts(t *testing.T) {
	store := &MockForecastStore{
		forecast: []models.DBForecast{
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
				CragId:              2,
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
			{
				Id:                  3,
				Time:                "2024-04-06T13:00:00Z",
				ScreenTemperature:   22.3,
				FeelsLikeTemp:       20.1,
				WindSpeed:           12.5,
				WindDirection:       200.0,
				TotalPrecipAmount:   0.8,
				ProbOfPrecipitation: 40.0,
				Latitude:            41.01,
				Longitude:           41.11,
				CragId:              3,
			},
		},
	}
	handler := NewHandler(store)
	router := mux.NewRouter()

	router.PathPrefix("/forecast/all").HandlerFunc(handler.GetAllForecasts()).Methods("GET")

	t.Run("Testing Valid Request", func(t *testing.T) {
		res, req := util.NewGetRequest("/forecast/all")

		router.ServeHTTP(res, req)
		//little reminder for future me, when we initialise a value like this its nil.
		//To deocode, we need to provide a pointer to this value, not a copy of it
		//or else we will get nil and itll be really annoying
		var data map[int][]models.DBForecast

		_, err := util.DecodeResponse(res.Body, &data)
		if err != nil {
			t.Fatalf("Could not decode response %s", err)
		}

		testData := map[int][]models.DBForecast{
			2: {
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
					CragId:              2,
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
			3: {
				{
					Id:                  3,
					Time:                "2024-04-06T13:00:00Z",
					ScreenTemperature:   22.3,
					FeelsLikeTemp:       20.1,
					WindSpeed:           12.5,
					WindDirection:       200.0,
					TotalPrecipAmount:   0.8,
					ProbOfPrecipitation: 40.0,
					Latitude:            41.01,
					Longitude:           41.11,
					CragId:              3,
				},
			},
		}

		assert.Equal(t, testData, data)
	})
}

func TestDeleteForecast(t *testing.T) {
	store := &MockForecastStore{
		forecast: []models.DBForecast{
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
				CragId:              2,
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
			{
				Id:                  3,
				Time:                "2024-04-06T13:00:00Z",
				ScreenTemperature:   22.3,
				FeelsLikeTemp:       20.1,
				WindSpeed:           12.5,
				WindDirection:       200.0,
				TotalPrecipAmount:   0.8,
				ProbOfPrecipitation: 40.0,
				Latitude:            41.01,
				Longitude:           41.11,
				CragId:              3,
			},
		},
	}
	handler := NewHandler(store)
	router := mux.NewRouter()

	router.PathPrefix("/forecast/{Id}").HandlerFunc(handler.handleDeleteForecastById()).Methods("GET")

	t.Run("Testing Valid ID", func(t *testing.T) {
		res, req := util.NewGetRequest("/forecast/1")

		router.ServeHTTP(res, req)

		assert.Equal(t, 200, res.Code)
	})

}
