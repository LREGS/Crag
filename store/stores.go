package store

import "github.com/lregs/Crag/models"

type CragStore interface {
	StoreCrag(crag *models.Crag) (err error)
	//reminder that im returning a copy of the crag and not a pointer for better type safety?!
	GetCrag(Id int) (models.Crag, error)
	UpdateCragValue(name string, crag models.Crag) error
	DeleteCragByID(Id int) error
}

type ClimbStore interface {
	StoreClimb(climb *models.Climb) (err error)
	GetClimbsByCrag(CragId int) ([]*models.Climb, error)
	GetAllClimbs() []*models.Climb
	GetClimbById(Id int) (*models.Climb, error)
	UpdateClimb(climb *models.Climb) (*models.Climb, error)
	DeleteClimb(Id int) error
}

type ForecastStore interface {
	AddForecast(models.DBForecast) (*models.DBForecast, error)
	GetForecastByCragId(CragId int) ([]models.DBForecast, error)
	GetAllForecasts() (map[int][]models.DBForecast, error)
	DeleteForecastById(Id int) error
}
