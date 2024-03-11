package store

import "github.com/lregs/Crag/models"

type CragStore interface {
	StoreCrag(crag *models.Crag) (err error)
	GetCrag(Id int) error
}
