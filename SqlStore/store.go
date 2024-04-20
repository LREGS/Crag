package store

import (
	"database/sql"
)

type SqlStore struct {
	masterX *sql.DB
	Stores  SqlStoreStores
}

type SqlStoreStores struct {
	CragStore CragStore
	//this should be the interface not a pointer to the concrete type but cba to change this second
	ClimbStore    *SqlClimbStore
	ForecastStore *SqlForecastStore
}

type StoreConfig struct {
	DbConnection *sql.DB
}

func NewSqlStore(c *StoreConfig) (*SqlStore, error) {
	store := &SqlStore{}
	store.initConnect(c)
	// err := store.initConnection(c)
	// if err != nil {
	// 	return nil, errors.New("error getting connection")
	// }

	store.Stores.CragStore = NewCragStore(store)
	store.Stores.ClimbStore = NewClimbStore(store)
	store.Stores.ForecastStore = NewForecastStore(store)

	return store, nil
}

func (ss *SqlStore) initConnect(c *StoreConfig) {

	ss.masterX = c.DbConnection

	// DBURL, err := env.DBString()
	// if err != nil {
	// 	fmt.Printf("error establishing db connection %s", err)
	// 	return err
	// }

	// ss.masterX, err = sql.Open("postgres", DBURL)
	// if err != nil {
	// 	return err
	// }
	// return nil
}

func (ss *SqlStore) GetMasterX() *sql.DB {
	return ss.masterX
}

func (ss *SqlStore) GetCragStore() CragStore {
	return ss.Stores.CragStore
}

func (ss *SqlStore) Insert(query string, params ...any) *sql.Row {

	row := ss.masterX.QueryRow(query, params)
	return row
}
