package main

import (
	"fmt"
	"testing"
	"time"
)

func TestBuild(t *testing.T) {

	t.Run("Testing Valid Payload", func(t *testing.T) {
		name := "stanage"
		b := NewForecastTable(name, sampleForecastTotals(), sampleWindows())
		b.build()
		fmt.Print(b.table.String())

	})

}

func sampleForecastTotals() map[string]*ForecastTotals {
	return map[string]*ForecastTotals{
		"day1": {
			HighestTemp:  1.5,
			LowestTemp:   15.2,
			AvgTemp:      22.3,
			AvgWindSpeed: 5.8,
			TotalPrecip:  12.7,
			Datapoints:   24,
		},
		"day2": {
			HighestTemp:  29.1,
			LowestTemp:   16.4,
			AvgTemp:      23.7,
			AvgWindSpeed: 6.2,
			TotalPrecip:  8.9,
			Datapoints:   24,
		},
	}
}

func sampleWindows() [][]time.Time {
	return [][]time.Time{
		{
			time.Date(2023, 11, 4, 14, 0, 0, 0, time.UTC),
			time.Date(2023, 11, 4, 15, 0, 0, 0, time.UTC),
		},
		{
			time.Date(2023, 11, 4, 17, 0, 0, 0, time.UTC),
			time.Date(2023, 11, 4, 19, 0, 0, 0, time.UTC),
		},
	}
}

func sampleWindow() Window {
	return Window{
		Time:                       []string{"2023-11-04T14:00:00Z", "2023-11-04T15:00:00Z"},
		AvgTemp:                    20.5,
		AvgWindSpeed:               6.5,
		WindDirection:              180,
		PrecipInLastXAmountofHours: 10,
	}
}

func sampleForecastPayload() ForecastPayload {
	return ForecastPayload{
		LastModelRunTime: time.Now().Format(time.RFC3339),
		Totals:           sampleForecastTotals(),
		Windows:          sampleWindows(),
	}
}
