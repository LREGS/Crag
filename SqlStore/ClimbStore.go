package store

import (
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

func (cs *SqlClimbStore) StoreClimb(climb *models.Climb) (*models.Climb, error) {
	var storedClimb models.Climb

	err := cs.Store.masterX.QueryRow(StoreClimbQuery, climb.Name, climb.Grade, climb.CragID).Scan(&storedClimb.Id, &storedClimb.Name, &storedClimb.Grade, &storedClimb.CragID)
	if err != nil {
		return nil, err
	}
	return &storedClimb, nil
}

const GetClimbsAtCrag = `SELECT * FROM climb WHERE CragID = $1`

func (cs *SqlClimbStore) GetClimbsByCrag(CragId int) ([]*models.Climb, error) {

	rows, err := cs.Store.masterX.Query(GetClimbsAtCrag, CragId)
	if err != nil {
		return nil, err
	}

	var results []*models.Climb

	for rows.Next() {
		climb := &models.Climb{}
		err := rows.Scan(&climb.Id, &climb.Name, &climb.Grade, &climb.CragID)
		if err != nil {
			return nil, err
		}

		results = append(results, climb)

	}

	return results, err
}

const getAllClimbs = `SELECT * FROM CLIMB ORDER BY name`

func (cs *SqlClimbStore) GetAllClimbs() ([]*models.Climb, error) {

	var results []*models.Climb

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
		results = append(results, &climb)
	}

	return results, nil
}

const getClimbById = `SELECT * FROM climb WHERE Id = $1`

func (cs *SqlClimbStore) GetClimbById(Id int) (*models.Climb, error) {

	var climb models.Climb

	rows := cs.Store.masterX.QueryRow(getClimbById, Id)

	err := rows.Scan(&climb.Id, &climb.Name, &climb.Grade, &climb.CragID)
	if err != nil {
		return nil, err
	}

	return &climb, nil

}

const updateClimb = `update climb set Name = $1, Grade = $2, CragID = $3 WHERE Id = $4
RETURNING *`

func (cs *SqlClimbStore) UpdateClimb(climb *models.Climb) (*models.Climb, error) {
	var updatedClimb models.Climb

	rows := cs.Store.masterX.QueryRow(updateClimb, &climb.Name, &climb.Grade, &climb.CragID, &climb.Id)

	err := rows.Scan(&updatedClimb.Id, &updatedClimb.Name, &updatedClimb.Grade, &updatedClimb.CragID)
	if err != nil {
		return nil, err
	}

	return &updatedClimb, nil

}

const deleteClimb = `DELETE FROM climb WHERE id = $1`

func (cs *SqlClimbStore) DeleteClimb(Id int) error {

	_, err := cs.Store.masterX.Exec(deleteClimb, Id)
	if err != nil {
		return err
	}
	return nil
}

func (cs *SqlClimbStore) Validate(*models.Climb) error {
	return nil
}
