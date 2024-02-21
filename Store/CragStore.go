package Store

import "workspaces/github.com/lregs/Crag/data"

type SqlCragStore struct {
	*SqlStore
}

func NewSqlCragStore(sqlStore *SqlStore) CragStore {
	cs := &SqlCragStore{
		SqlStore: sqlStore,
	}
	return cs
}

func (cs SqlCragStore) StoreCrag(crag data.Crag) (err error) {

	Query := `
	insert into crag(
	Name 
	Latitude
	Longtitude)
	values(
		$1, $2, $3, $4
	)
	returning id`

	err = cs.SqlStore.masterX.QueryRow(Query, crag.Name, crag.Latitude, crag.Longitude).Scan(&crag.Id)
	return err
}

func (cs SqlCragStore) GetCrag(id int) (crag data.Crag, err error) {
	crag = data.Crag{}
	crag.Climbs = []data.Climb{}
	crag.Reports = []data.Report{}

	err = cs.SqlStore.masterX.QueryRow("select id, Name, Latitude, Longitude from crag where id = $1", id).Scan(&crag.Id, &crag.Name, &crag.Latitude, &crag.Longitude)
	if err != nil {
		return crag, err
	}

	reportRows, err := cs.SqlStore.masterX.Query("sllect Id, Content, Author from Report where CragID = $1", id)
	if err != nil {
		return crag, err
	}
	climbRows, err := cs.SqlStore.masterX.Query("select Id, Name, Grade, from climb where CragID = $1", id)
	if err != nil {
		return crag, err
	}

	for reportRows.Next() {
		report := data.Report{Crag: &crag}

		reportErr := reportRows.Scan(&report.Id, &report.Author, &report.Content)
		if reportErr != nil {
			return crag, reportErr
		}

		crag.Reports = append(crag.Reports, report)
	}
	reportRows.Close()

	for climbRows.Next() {
		climbs := data.Climb{Crag: &crag}

		climbsErr := climbRows.Scan(&climbs.Id, &climbs.Name, &climbs.Grade)
		if climbsErr != nil {
			return crag, climbsErr
		}
		crag.Climbs = append(crag.Climbs, climbs)
	}
	climbRows.Close()

	return crag, nil

}

func (cs SqlCragStore) UpdateCrag(crag data.Crag) error {
	_, err := cs.SqlStore.masterX.Exec("update crag set Name = $2, Latitude = $3, Longitude = $4 where id = $1", crag.Id, crag.Name, crag.Latitude, crag.Longitude)
	return err
}

func (cs SqlCragStore) DeleteCrag(id int) error {
	_, err := cs.SqlStore.masterX.Exec("delete from crag where id = $1", id)
	return err
}
