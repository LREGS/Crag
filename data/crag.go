package data

type Crag struct {
	Id       int
	Name     string
	location []float64
	Climbs   []Climb
	Forecast Forecast
}

type Climb struct {
	Id   int
	Name string
	Crag string
}

type Report struct {
	Id      int
	Content string
	Author  string
}
