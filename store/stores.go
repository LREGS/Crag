package store

import "github.com/lregs/Crag/models"

type CragStore interface {
	StoreCrag(crag *models.Crag) (err error)
	//reminder that im returning a copy of the crag and not a pointer for better type safety?!
	GetCrag(Id int) (models.Crag, error)
	UpdateCragValue(name string, crag models.Crag) error
}
