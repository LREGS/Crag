package store

import (
	"database/sql"

	"github.com/lregs/Crag/models"
)

//go:generate moq -out storess_test.go . CragStore

//Storage layer for all the types. By the time the data has reached the store it will have already been validated **TODO**

type Store interface {
	initConnect(*StoreConfig)
	GetMasterX() *sql.DB
	GetCragStore() CragStore
}

type CragStore interface {
	StoreCrag(crag models.CragPayload) (models.Crag, error)
	GetCrag(Id int) (models.Crag, error)
	UpdateCrag(crag models.Crag) (models.Crag, error)
	DeleteCragByID(Id int) error
	Validate(models.CragPayload) error
}

type ClimbStore interface {
	StoreClimb(climb *models.Climb) (*models.Climb, error)
	GetClimbsByCrag(CragId int) ([]*models.Climb, error)
	GetAllClimbs() ([]*models.Climb, error)
	GetClimbById(Id int) (*models.Climb, error)
	UpdateClimb(climb *models.Climb) (*models.Climb, error)
	DeleteClimb(Id int) error
	Validate(climb *models.Climb) error
}

type ForecastStore interface {
	AddForecast(*models.DBForecastPayload) (models.DBForecast, error)
	GetForecastByCragId(CragId int) ([]models.DBForecast, error)
	//im not sure we need this unles we want to seperate it into days in the store/with the query??
	GetAllForecastsByCragId() (map[int][]models.DBForecast, error)
	DeleteForecastById(Id int) error
	Validate(*models.DBForecast) error
}
