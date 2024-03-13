package store

import (
	"github.com/lregs/Crag/models"
)

type SqlCragStore struct {
	Store *SqlStore
}

func NewCragStore(sqlStore *SqlStore) *SqlCragStore {
	CS := &SqlCragStore{Store: sqlStore}
	return CS
}

func (cs *SqlCragStore) StoreCrag(crag *models.Crag) error {
	query := `insert into crag(Name, Latitude, Longitude) values($1,$2,$3)`

	_, err := cs.Store.masterX.Exec(query, crag.Name, crag.Latitude, crag.Longitude)
	if err != nil {
		return err
	}
	return nil

}

func (cs *SqlCragStore) GetCrag(Id int) (models.Crag, error) {
	var c models.Crag

	query := `select Id, Name, Latitude, Longitude from crag where id = $1`

	err := cs.Store.masterX.QueryRow(query, Id).Scan(
		&c.Id, &c.Name, &c.Latitude, &c.Longitude)

	return c, err
}

func (cs *SqlCragStore) UpdateCragValue(crag models.Crag) error {
	//takes a crag with the updates and strips the values that need updating
	query := "update crag set Name = $1 where Id = 1"
	_, err := cs.Store.masterX.Exec(query, "Milestone")
	if err != nil {
		return err
	}
	return nil
}
