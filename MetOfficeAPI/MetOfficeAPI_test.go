package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockStore struct {
}

func (s *MockStore) Add(ctx context.Context, name string, payload ForecastPayload) error {
	return nil
}

func (s *MockStore) Get() (map[string]*ForecastTotals, error) {
	return nil, nil
}

func GetLastUpdate() time.Time {
	return time.Time{}
}

func TestGetForecast(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		code int
		body []byte
	}{
		{
			name: "bad code",
			code: http.StatusInternalServerError,
			body: []byte{},
		},
		{
			name: "available forecast",
			code: http.StatusOK,
			body: GetTestForecast(t),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			t.Logf("running %s", tc.name)

			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.code)
				_, err := w.Write(tc.body)
				if err != nil {
					t.Fatalf("failled writing body %s", err)
				}
			}))

			defer testServer.Close()

			log := NewLogger("testLog.txt")

			api := NewMetAPI("", log)
			api.BaseURL = testServer.URL

			res, err := api.GetForecast(api.BaseURL)
			if err != nil {
				assert.Equal(t, "code 500", err.Error())
				assert.Equal(t, Forecast{}, res)
			}

			if tc.code != 500 {

				var testData Forecast
				if err := json.Unmarshal(tc.body, &testData); err != nil {
					t.Fatalf("failed decoding test data %s", err)
				}
				assert.Equal(t, testData, res)
			}

		})
	}
}

func TestMakeURL(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name     string
		coords   []float64
		baseKey  string
		expected string
	}{
		{
			name:     "correct",
			coords:   []float64{53.04567339439013, -4.021447439922229},
			baseKey:  "string",
			expected: fmt.Sprintf("https://data.hub.api.metoffice.gov.uk/sitespecific/v0/point/hourly?latitude=53.04567339439013longitude=-4.021447439922229"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			log := NewLogger("testLog.txt")
			api := NewMetAPI(tc.baseKey, log)

			assert.Equal(t, tc.expected, api.CreateURL(tc.coords))

		})
	}
}

func TestCalculateTotals(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name          string
		TestData      []TimeSeriesData
		expectedAvg   map[string]*ForecastTotals
		expectedError bool
	}{
		{
			name: "Valid Data",
			TestData: []TimeSeriesData{
				{
					Time:              "2024-08-05T00:00Z",
					ScreenTemperature: 20.5,
					MaxScreenAirTemp:  25.0,
					MinScreenAirTemp:  15.0,
					TotalPrecipAmount: 1.5,
				},
				{
					Time:              "2024-08-05T01:00Z",
					ScreenTemperature: 19.5,
					MaxScreenAirTemp:  25.0,
					MinScreenAirTemp:  15.0,
					TotalPrecipAmount: 0.5,
				},
				{
					Time:              "2024-08-05T02:00Z",
					ScreenTemperature: 18.5,
					MaxScreenAirTemp:  25.0,
					MinScreenAirTemp:  15.0,
					TotalPrecipAmount: 0.0,
				},
				{
					Time:              "2024-08-06T00:00Z",
					ScreenTemperature: 21.0,
					MaxScreenAirTemp:  26.0,
					MinScreenAirTemp:  16.0,
					TotalPrecipAmount: 0.2,
				},
			},

			expectedAvg: map[string]*ForecastTotals{
				"05": {
					HighestTemp: 25.0,
					LowestTemp:  15.0,
					AvgTemp:     19.5,
					TotalPrecip: 2.0,
					Datapoints:  3,
				},
				"06": {
					HighestTemp: 26.0,
					LowestTemp:  16.0,
					AvgTemp:     21.0,
					TotalPrecip: 0.2,
					Datapoints:  1,
				},
			},
			expectedError: false,
		},
		// {
		// 	name:          "incorrect Data",
		// 	TestData:      []TimeSeriesData{},
		// 	expectedAvg:   nil,
		// 	expectedError: true,
		// },
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			log := NewLogger("tLog.txt")
			api := NewMetAPI(" ", log)

			d := api.CalculateTotals(tc.TestData)
			t.Logf("%s %v", tc.name, d)
			assert.Equal(t, tc.expectedAvg, d)

		})
	}
}

func TestFindWindows(t *testing.T) {

	cases := []struct {
		name     string
		data     []TimeSeriesData
		expected [][]time.Time
	}{
		{
			name: "two hour window",
			data: []TimeSeriesData{
				{
					Time:              "2024-08-05T00:00Z",
					TotalPrecipAmount: 1.0,
				},
				{
					Time:              "2024-08-05T01:00Z",
					TotalPrecipAmount: 0.0,
				},
				{
					Time:              "2024-08-05T02:00Z",
					TotalPrecipAmount: 0.0,
				},
				{
					Time:              "2024-08-05T03:00Z",
					TotalPrecipAmount: 1.0,
				},
			},
			expected: [][]time.Time{
				{time.Date(2024, 8, 5, 1, 0, 0, 0, time.UTC), time.Date(2024, 8, 5, 2, 0, 0, 0, time.UTC)},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			log := NewLogger("tLog.txt")
			api := NewMetAPI(" ", log)

			wins := api.FindWindows(tc.data)

			assert.Equal(t, tc.expected, wins)
		})
	}
}

func GetTestForecast(t *testing.T) []byte {
	t.Helper()
	file, err := os.ReadFile("./test/sampleData.json")
	if err != nil {
		t.Fatalf("failed getting json %s", err)
	}
	return file
}
