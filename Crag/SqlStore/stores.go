package store

import (
	"context"
	"database/sql"

	"github.com/lregs/CragWeather/Crag/models"
)

//go:generate moq -out storess_test.go . CragStore

//Storage layer for all the types. By the time the data has reached the store it will have already been validated **TODO**

type Store interface {
	initConnect(*StoreConfig)
	GetMasterX() *sql.DB
	GetCragStore() CragStore
}

type CragStore interface {
	StoreCrag(ctx context.Context, crag models.CragPayload) (models.Crag, error)
	GetCrag(ctx context.Context, Id int) (models.Crag, error)
	UpdateCrag(ctx context.Context, crag models.Crag) (models.Crag, error)
	DeleteCragByID(ctx context.Context, Id int) (models.Crag, error)
}

type ClimbStore interface {
	StoreClimb(ctx context.Context, climb models.ClimbPayload) (models.Climb, error)
	GetClimbsByCragId(ctx context.Context, CragId int) ([]models.Climb, error)
	GetAllClimbs(ctx context.Context) ([]models.Climb, error)
	GetClimbById(ctx context.Context, Id int) (models.Climb, error)
	UpdateClimb(ctx context.Context, climb models.Climb) (models.Climb, error)
	DeleteClimb(ctx context.Context, Id int) (models.Climb, error)
}
