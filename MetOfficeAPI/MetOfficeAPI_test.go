package main

import (
	"context"
	"encoding/json"
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

// func TestScheduleMetOffice(t *testing.T) {

// 	cases := []struct {
// 		name string
// 		t    time.Time
// 	}{
// 		{
// 			name: "Time Now",
// 			t:    time.Now(),
// 		},
// 		{
// 			name: "Immediate update required",
// 			t:    time.Now().Add(-2 * time.Hour),
// 		},
// 	}

// 	for _, tc := range cases {
// 		t.Run(tc.name, func(t *testing.T) {

// 			log := NewLogger("testLog")
// 			t.Log("started")
// 			api := NewMetAPI(os.Getenv("apikey"), log)

// 			api.ScheduleMetOffice(tc.t)

// 		})
// 	}

// }

func GetTestForecast(t *testing.T) []byte {
	t.Helper()
	file, err := os.ReadFile("./test/sampleData.json")
	if err != nil {
		t.Fatalf("failed getting json %s", err)
	}
	return file
}
