package store

import (
	"database/sql"
)

type SqlStore struct {
	masterX *sql.DB
	Stores  SqlStoreStores
}

type SqlStoreStores struct {
	CragStore  CragStore
	ClimbStore *SqlClimbStore
}

type StoreConfig struct {
	dbConnection *sql.DB
}

func NewSqlStore(c *StoreConfig) (*SqlStore, error) {
	store := &SqlStore{}
	store.initConnection(c)
	// err := store.initConnection(c)
	// if err != nil {
	// 	return nil, errors.New("error getting connection")
	// }

	store.Stores.CragStore = NewCragStore(store)
	store.Stores.ClimbStore = NewClimbStore(store)

	return store, nil
}

func (ss *SqlStore) initConnection(c *StoreConfig) {

	ss.masterX = c.dbConnection

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
