package store

import "github.com/lregs/Crag/models"

type SqlCragStore struct {
	*SqlStore
}

func (cs *SqlCragStore) StoreCrag(crag models.Crag) error {
	query := `insert into crag(Name, Latitude, Longitude) values($1,$2,$3)`

	_, err := cs.SqlStore.masterX.Exec(query, crag.Name, crag.Latitude, crag.Longitude)
	if err != nil {
		return err
	}
	return nil

}
