package store

import (
	"database/sql"

	"github.com/lregs/Crag/models"
)

//go:generate moq -out storess_test.go . CragStore

type Store interface {
	initConnect(*StoreConfig)
	GetMasterX() *sql.DB
	GetCragStore() CragStore
}

type CragStore interface {
	//Do I not want to be returning the store instance so it can be checked whether the correct data was stored
	//pls
	StoreCrag(crag *models.Crag) (err error)
	//reminder that im returning a copy of the crag and not a pointer for better type safety?! - should I be cause this has changed...
	//pls go back to returing a copy and not pointer?
	GetCrag(Id int) (*models.Crag, error)
	UpdateCragValue(crag models.Crag) error
	DeleteCragByID(Id int) error
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
