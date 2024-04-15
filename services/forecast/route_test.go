package forecast

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/lregs/Crag/models"
	"github.com/lregs/Crag/util"
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

func (fs *MockForecastStore) AddForecast(newForecast *models.DBForecastPayload) (models.DBForecast, error) {

	fToStore := fs.returnDBForecast(newForecast, (len(fs.forecast) + 1))

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
	if CragId == 0 {
		return nil, errors.New("id 0 not valid")
	}

	// res := []models.DBForecast{}

	// for _, f := range fs.forecast {
	// 	if f.CragId == CragId {
	// 		res = append(res, f)
	// 	}
	// }

	// if len(res) == 0 {
	// 	return res, errors.New("no forecast where found for crag Id")
	// } else {
	// 	return res, nil
	// }

	res := []models.DBForecast{
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
	}
	return res, nil

}

func (fs *MockForecastStore) GetAllForecastsByCragId() (map[int][]models.DBForecast, error) {

	if len(fs.forecast) == 0 {
		return nil, errors.New("DB is empty")
	} else {
		res := make(map[int][]models.DBForecast, 0)

		for i, f := range fs.forecast {
			res[i] = append(res[i], f)
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

	handler := NewHanlder(store)
	router := mux.NewRouter()

	router.PathPrefix("/forecast").HandlerFunc(handler.handlePostForecast()).Methods("POST")

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

		response := httptest.NewRecorder()

		body, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("could not marshal body because of err: %s", err)
		}
		request, err := util.NewPostRequest(body, "/forecast")
		if err != nil {
			t.Fatalf("error %s making new request", err)
		}

		router.ServeHTTP(response, request)

		switch response.Code {
		case 200:
			var data models.DBForecast

			_, err := util.DecodeResponse(response.Body, data)
			if err != nil {
				t.Fatalf("Error decoding response: %s", err)
			}

			if data.Time == "" {
				t.Fatalf("inv")
			}

		}

		t.Run("Testing Invalid Data", func(t *testing.T) {

			payload := models.Crag{Id: 2, Name: "dank", Longitude: 1.1, Latitude: 2.2}
			body, err := json.Marshal(payload)
			if err != nil {
				t.Fatalf("marshall failed: %s", err)
			}
			response := httptest.NewRecorder()
			request, err := util.NewPostRequest(body, "/forecast")
			if err != nil {
				t.Fatalf("new post request failed %s", err)
			}

			router.ServeHTTP(response, request)

			if response.Code != 500 {
				t.Fatalf("Server did not handle incorrect data type")
			}

		})

		t.Run("Testing Invalid Request Method", func(t *testing.T) {
			response := httptest.NewRecorder()
			request := util.NewGetRequest("/forecast")

			router.ServeHTTP(response, request)

			if response.Code != 400 {
				t.Fatalf("Accepted incorrect method")
			}

		})

		// var res util.Response

		// //just check response code to decide whether to unmarshaoll into error or not ffs

		// _, err = util.DecodeResponse(response.Body, &res)
		// if err != nil {
		// 	t.Fatalf("could not decode because of err :%s", err)
		// }

		// if res.Error != "" {
		// 	t.Fatalf("Error storing forecast: %s", res.Error)
		// }

		// if len(store.forecast) != 3 {
		// 	t.Fatal("Nothing stored")
		// }

		// util.AssertStatus(t, response.Code, http.StatusOK)

	})

	// t.Run("Testing Invalid Request Type", func(t *testing.T) {
	// 	response := httptest.NewRecorder()
	// 	request := util.NewGetRequest("/forecast")

	// 	router.ServeHTTP(response, request)

	// 	util.AssertStatus(t, response.Code, http.StatusMethodNotAllowed)
	// })

	// t.Run("Testing Invalid Data Type", func(t *testing.T) {
	// 	response := httptest.NewRecorder()

	// 	payload := models.Climb{
	// 		Id:     1,
	// 		Name:   "test",
	// 		Grade:  "v2",
	// 		CragID: 2,
	// 	}

	// 	body, err := json.Marshal(payload)
	// 	if err != nil {
	// 		t.Fatalf("Marshal failed: %s", err)
	// 	}

	// 	request, err := util.NewPostRequest(body, "/forecast")
	// 	if err != nil {
	// 		t.Fatalf("Getting request failed: %s", err)
	// 	}

	// 	router.ServeHTTP(response, request)

	// 	util.CheckErrs(t, response, "invalid data")

	// 	util.AssertStatus(t, response.Code, http.StatusBadRequest)

	// })

}

func TestGetForecastByCragId(t *testing.T) {
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

	handler := NewHanlder(store)
	router := mux.NewRouter()

	router.PathPrefix("/forecast/{Id}").HandlerFunc(handler.handleGetForecastByCragId()).Methods("GET")

	t.Run("Valid CragID", func(t *testing.T) {

		response := httptest.NewRecorder()
		request := util.NewGetRequest("/forecast/2")
		router.ServeHTTP(response, request)

		if response.Code == 200 {
			var res []models.DBForecast
			_, err := util.DecodeResponse(response.Body, &res)
			if err != nil {
				t.Fatalf("Could not decode response :%s", err)
			}
		}

		// util.AssertStatus(t, response.Code, 200)

		// var r getResponse

		// _, err := util.DResponse(response.Body, &r)
		// if err != nil {
		// 	t.Fatalf("Could not decode response because of err: %s", err)
		// }

		// respErr := r.GetError()
		// if respErr != "" {
		// 	t.Fatalf("error in response %s", respErr)
		// }

		// data, ok := r.Data.([]models.DBForecast)
		// if !ok {
		// 	t.Fatalf("no data returned")
		// }
		// if data[0].Id != 1 {
		// 	t.Fatalf("got %d wanted %d", data[0].Id, 1)
		// }

	})
}
