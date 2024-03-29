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

func (cs *SqlCragStore) GetCrag(Id int) (*models.Crag, error) {
	c := &models.Crag{}

	query := `select Id, Name, Latitude, Longitude from crag where id = $1`

	err := cs.Store.masterX.QueryRow(query, Id).Scan(
		&c.Id, &c.Name, &c.Latitude, &c.Longitude)

	return c, err
}

func (cs *SqlCragStore) UpdateCragValue(name string, crag models.Crag) error {
	//I'd think there has to be a way to make a query builder so I can be selective about which
	//field I want to update without having a set method for each field.
	//maybe do this in the future
	query := `
	update crag set 
	Name = $1, 
	Latitude = $2,
	Longitude = $3 
	where Name = $4`
	_, err := cs.Store.masterX.Exec(query, crag.Name, crag.Latitude, crag.Longitude, name)
	if err != nil {
		return err
	}
	return nil
}

func (cs *SqlCragStore) DeleteCragByID(Id int) error {
	query := `delete from crag where id = $1`
	_, err := cs.Store.masterX.Exec(query, Id)
	if err != nil {
		return err
	}
	return nil
}
