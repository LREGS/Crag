package data

import (
	"database/sql"
	"errors"
	helpers "workspaces/github.com/lregs/Crag/helper"

	_ "github.com/lib/pq"
)

type Crag struct {
	Id        int
	Name      string
	Latitude  float64
	Longitude float64
	Climbs    []Climb  //many to one relationship not db field
	Reports   []Report //many to one relationship not db field
	Forecast  []Forecast
}

type Climb struct {
	Id    int
	Name  string
	Grade string
	Crag  *Crag
}

type Report struct {
	Id      int
	Content string
	Author  string
	Crag    *Crag
}

func (crag *Crag) Create(db *sql.DB) (err error) {
	return
}

func GetCrag(id int, db *sql.DB) (crag Crag, err error) {
	crag = Crag{}
	crag.Climbs = []Climb{}
	crag.Reports = []Report{}
	err = db.QueryRow("select id, Name, Latitude, Longitude from crag where id = $1", id).Scan(&crag.Id, &crag.Name, &crag.Latitude, &crag.Longitude)

	reportRows, err := db.Query("select Id, Content, Author from Report where CragID = $1", id)
	climbRows, err := db.Query("select Id, Name, Grade from climb where CragID = $1", id)
	helpers.CheckError(err)

	for reportRows.Next() {
		report := Report{Crag: &crag}

		reportErr := reportRows.Scan(&report.Id, &report.Author, &report.Content)
		helpers.CheckError(reportErr)
		crag.Reports = append(crag.Reports, report)

	}
	reportRows.Close()

	for climbRows.Next() {
		climbs := Climb{Crag: &crag}

		climbsErr := climbRows.Scan(&climbs.Id, &climbs.Name, &climbs.Grade)
		helpers.CheckError(climbsErr)
		crag.Climbs = append(crag.Climbs, climbs)
	}
	climbRows.Close()

	if crag.Id == 0 {
		return crag, errors.New("This Crag Doesn't Exist")
	}

	return crag, nil
}

func (crag *Crag) UpdateCrag(db *sql.DB) (err error) {
	_, err = db.Exec("update crag set Name = $2, Latitude = $3, Longitude = $4 where id = $1", crag.Id, crag.Name, crag.Latitude, crag.Longitude)
	return
}

func (crag *Crag) DeleteCrag(db *sql.DB) (err error) {
	_, err = db.Exec("delete from crag where id = $1", crag.Id)
	return
}

func (Report *Report) Create(db *sql.DB) (err error) {
	err = db.QueryRow("insert into report (Content, Author, cragID) vales($1, $2, $3) returning id", Report.Content, Report.Author, Report.Crag.Id).Scan(&Report.Id)
	return
}

func (Climb *Climb) Create(db *sql.DB) (err error) {
	err = db.QueryRow("insert into climb (Name, Grade, cragID) vales($1, $2, $3) returning id", Climb.Name, Climb.Grade, Climb.Crag.Id).Scan(&Climb.Id)
	return

}
