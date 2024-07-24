package main

import (
	"net/http"
	"testing"
)

func TestGetAvailableReleases(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name     string
		code     int
		body     string
		forecast Forecast
	}{
		{
			name: "bad code",
			code: http.StatusInternalServerError,
		},
		{
			name: "available forecast", 
			code: http.StatusOK,
			body: `[
			
			
			]`
		},
}
