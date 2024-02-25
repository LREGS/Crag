package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGETCrags(t *testing.T) {
	t.Run("returns forecast of a crag", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/crags/forecast", nil)
		response := httptest.NewRecorder()

		CragServer(response, request)

		got := response.Body.String()
		want := "cold"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
