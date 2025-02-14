package store

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type SqlStore struct {
	masterX *pgxpool.Pool
	Stores  SqlStoreStores
}

type SqlStoreStores struct {
	CragStore CragStore
	//this should be the interface not a pointer to the concrete type but cba to change this second
	ClimbStore ClimbStore
}

type StoreConfig struct {
	DbConnection *pgxpool.Pool
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

	return store, nil
}

func (ss *SqlStore) initConnect(c *StoreConfig) {

	ss.masterX = c.DbConnection

}

func (ss *SqlStore) GetMasterX() *pgxpool.Pool {
	return ss.masterX
}

func (ss *SqlStore) GetCragStore() CragStore {
	return ss.Stores.CragStore
}
