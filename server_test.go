package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGETCrags(t *testing.T) {

	server := &CragServer{}

	t.Run("returns forecast of a stange", func(t *testing.T) {
		request := newGetForecastRequest("stanage")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertBodyResponse(t, response.Body.String(), "cold")

	})

	t.Run("returns forecast of milestone", func(t *testing.T) {
		request := newGetForecastRequest("milestone")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertBodyResponse(t, response.Body.String(), "dry")

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
