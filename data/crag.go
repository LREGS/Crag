package data

type Crag struct {
	Id       int
	Name     string
	location []float64
}

type Climb struct {
	Id     int
	Name   string
	Grade  string
	CragID int
}

type Report struct {
	Id      int
	Content string
	Author  string
	CragID  int
}
