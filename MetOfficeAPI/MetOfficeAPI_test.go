package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
				t.Fatalf("get forecast failed %s", err)
			}

			var testData Forecast
			if err := json.Unmarshal(tc.body, &testData); err != nil {
				t.Fatalf("failed decoding test data %s", err)
			}

			if tc.code == 500 {
				return
			}
			assert.Equal(t, testData, res)

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
