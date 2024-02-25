package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGETCrags(t *testing.T) {

	store := StubCragStore{
		map[string]string{
			"stanage":   "cold",
			"milestone": "dry",
		},
		nil,
	}

	server := &CragServer{&store}

	t.Run("returns forecast of stanage", func(t *testing.T) {
		request := newGetForecastRequest("stanage")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertBodyResponse(t, response.Body.String(), "cold")

	})

	t.Run("returns forecast of milestone", func(t *testing.T) {
		request := newGetForecastRequest("milestone")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertBodyResponse(t, response.Body.String(), "dry")

	})

	t.Run("returns 404 on missing crags", func(t *testing.T) {
		request := newGetForecastRequest("unkown")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusNotFound)
	})
}

func newGetForecastRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/crags/%s", name), nil)
	return req

}

func assertBodyResponse(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q and wanted %q", got, want)
	}
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("status response was expected, got %d, want %d", got, want)
	}
}

type StubCragStore struct {
	crags     map[string]string
	forecasts []string
}

func (s *StubCragStore) GetForecast(name string) string {
	forecast := s.crags[name]
	return forecast
}

func (s *StubCragStore) addForecast(forecast string) {
	s.forecasts = append(s.forecasts, forecast)
}

func TestStoreForecast(t *testing.T) {
	store := StubCragStore{
		crags:     map[string]string{},
		forecasts: []string{},
	}
	server := &CragServer{&store}

	t.Run("return accepted on POST", func(t *testing.T) {
		reportedForecast := "dry"
		request := newPostForecast(reportedForecast)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusAccepted)

		if len(store.forecasts) != 1 {
			t.Errorf("got %d calls to forecasts want %d", len(store.forecasts), 1)
		}

		if store.forecasts[0] != reportedForecast {
			t.Errorf("did not store correct forecast got %q want %q", store.forecasts[0], reportedForecast)
		}

	})
}

func newPostForecast(forecast string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("crags/stanage/%s", forecast), nil)
	return req

}
