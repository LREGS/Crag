package store

import (
	"errors"

	"github.com/lregs/Crag/models"
)

type SqlCragStore struct {
	Store *SqlStore
}

func NewCragStore(sqlStore *SqlStore) *SqlCragStore {
	CS := &SqlCragStore{Store: sqlStore}
	return CS
}

const storeCrag = `insert into crag(Name, Latitude, Longitude) values($1,$2,$3) RETURNING *`

func (cs *SqlCragStore) StoreCrag(crag models.CragPayload) (models.Crag, error) {

	var storedCrag models.Crag

	err := cs.Validate(crag)
	if err != nil {
		return storedCrag, err
	}

	err = cs.Store.masterX.QueryRow(storeCrag, crag.Name, crag.Latitude, crag.Longitude).Scan(&storedCrag.Id, &storedCrag.Name, &storedCrag.Latitude, &storedCrag.Longitude)
	if err != nil {
		return storedCrag, nil
	}

	return storedCrag, nil

}

const getCrag = `select Id, Name, Latitude, Longitude from crag where id = $1`

func (cs *SqlCragStore) GetCrag(Id int) (models.Crag, error) {
	var storedCrag models.Crag

	err := cs.Store.masterX.QueryRow(getCrag, Id).Scan(
		&storedCrag.Id, &storedCrag.Name, &storedCrag.Latitude, &storedCrag.Longitude)

	if err != nil {
		return storedCrag, err
	}

	return storedCrag, nil
}

const updateCrag = `update crag set Name = $1, Latitude = $2, Longitude = $3 where Id = $4 RETURNING *`

func (cs *SqlCragStore) UpdateCrag(crag models.Crag) (models.Crag, error) {

	var updatedCrag models.Crag

	err := cs.Store.masterX.QueryRow(updateCrag, crag.Name, crag.Latitude, crag.Longitude, crag.Id).Scan(&updatedCrag.Id, &updatedCrag.Name, &updatedCrag.Latitude, &updatedCrag.Longitude)

	if err != nil {
		return updatedCrag, err
	}
	return updatedCrag, nil
}

func (cs *SqlCragStore) DeleteCragByID(Id int) error {
	query := `delete from crag where id = $1`
	_, err := cs.Store.masterX.Exec(query, Id)
	if err != nil {
		return err
	}
	return nil
}

func (cs *SqlCragStore) Validate(payload models.CragPayload) error {
	//niave validation
	if payload.Name == "" {
		return errors.New("cannot use empty name")
	}
	if payload.Latitude < -90 || payload.Latitude > 90 {
		return errors.New("invalid latitude")
	}
	if payload.Longitude < -180 || payload.Longitude > 180 {
		return errors.New("invalid longitude")
	}
	return nil
}
