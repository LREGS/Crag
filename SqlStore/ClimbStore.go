package store

import (
	"errors"
	"reflect"
	"regexp"

	"github.com/lregs/Crag/models"
)

type SqlClimbStore struct {
	Store *SqlStore
}

func NewClimbStore(sqlStore *SqlStore) *SqlClimbStore {
	store := &SqlClimbStore{sqlStore}
	return store
}

const StoreClimbQuery = `insert into climb(Name, Grade, CragID) VALUES($1,$2,$3)RETURNING *`

func (cs *SqlClimbStore) StoreClimb(climb models.ClimbPayload) (models.Climb, error) {
	var storedClimb models.Climb

	err := cs.validatePayload(climb)
	if err != nil {
		return storedClimb, err
	}
	err = cs.Store.masterX.QueryRow(StoreClimbQuery, climb.Name, climb.Grade, climb.CragID).Scan(&storedClimb.Id, &storedClimb.Name, &storedClimb.Grade, &storedClimb.CragID)
	if err != nil {
		return storedClimb, err
	}
	return storedClimb, nil
}

const GetClimbsAtCrag = `SELECT * FROM climb WHERE CragID = $1`

func (cs *SqlClimbStore) GetClimbsByCragId(CragId int) ([]models.Climb, error) {
	//returns all climbs by their associated crag
	rows, err := cs.Store.masterX.Query(GetClimbsAtCrag, CragId)
	if err != nil {
		return nil, err
	}

	var results []models.Climb

	for rows.Next() {
		climb := models.Climb{}
		err := rows.Scan(&climb.Id, &climb.Name, &climb.Grade, &climb.CragID)
		if err != nil {
			return nil, err
		}

		results = append(results, climb)

	}

	return results, nil
}

const getAllClimbs = `SELECT * FROM CLIMB ORDER BY name`

func (cs *SqlClimbStore) GetAllClimbs() ([]models.Climb, error) {

	var results []models.Climb

	rows, err := cs.Store.masterX.Query(getAllClimbs)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var climb models.Climb
		err := rows.Scan(&climb.Id, &climb.Name, &climb.Grade, &climb.CragID)
		if err != nil {
			return nil, err
		}
		results = append(results, climb)
	}

	return results, nil
}

const getClimbById = `SELECT * FROM climb WHERE Id = $1`

func (cs *SqlClimbStore) GetClimbById(Id int) (models.Climb, error) {

	var climb models.Climb

	rows := cs.Store.masterX.QueryRow(getClimbById, Id)

	err := rows.Scan(&climb.Id, &climb.Name, &climb.Grade, &climb.CragID)
	if err != nil {
		return climb, err
	}

	return climb, nil

}

const updateClimb = `update climb set Name = $1, Grade = $2, CragID = $3 WHERE Id = $4
RETURNING *`

func (cs *SqlClimbStore) UpdateClimb(climb models.Climb) (models.Climb, error) {

	var updatedClimb models.Climb

	err := cs.validateClimb(climb)
	if err != nil {
		return updatedClimb, err
	}

	err = cs.Store.masterX.QueryRow(updateClimb, &climb.Name, &climb.Grade, &climb.CragID, &climb.Id).Scan(&updatedClimb.Id, &updatedClimb.Name, &updatedClimb.Grade, &updatedClimb.CragID)
	if err != nil {
		return updatedClimb, err
	}

	return updatedClimb, nil

}

const deleteClimb = `DELETE FROM climb WHERE id = $1 RETURNING *`

func (cs *SqlClimbStore) DeleteClimb(Id int) (models.Climb, error) {

	var deletedClimb models.Climb

	err := cs.Store.masterX.QueryRow(deleteClimb, Id).Scan(&deletedClimb.Id, &deletedClimb.Name, &deletedClimb.Grade, &deletedClimb.CragID)
	if err != nil {
		return deletedClimb, err
	}
	return deletedClimb, nil
}

func (cs *SqlClimbStore) validatePayload(data models.ClimbPayload) error {
	if reflect.DeepEqual(models.ClimbPayload{}, data) {
		return errors.New("value is empty")
	}

	if data.Name == "" {
		return errors.New("climb must have name")
	}

	r, _ := regexp.Compile(`[6-9][abc]\+?$`)
	if !r.MatchString(data.Grade) {
		return errors.New("climb grade invalid ")
	}

	if data.CragID == 0 {
		return errors.New("invalid crag ID")
	}
	return nil
}

func (cs *SqlClimbStore) validateClimb(data models.Climb) error {
	if reflect.DeepEqual(models.Climb{}, data) {
		return errors.New("value is empty")
	}

	if data.Id == 0 {
		return errors.New("invalid id")
	}

	if data.Name == "" {
		return errors.New("climb must have name")
	}

	r, _ := regexp.Compile(`[6-9][abc]\+?$`)
	if !r.MatchString(data.Grade) {
		return errors.New("climb grade invalid ")
	}

	if data.CragID == 0 {
		return errors.New("invalid crag ID")
	}
	return nil
}
