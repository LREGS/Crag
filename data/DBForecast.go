package data

import (
	_ "github.com/lib/pq"
)

type DBForecast struct {
	Id                  int
	crag                *Crag
	Time                string
	ScreenTemperature   float64
	FeelsLikeTemp       float64
	WindSpeed           float64
	WindDirection       float64
	totalPrecipAmount   float64
	ProbOfPrecipitation float64
	Latitude            float64
	Longitude           float64
}
