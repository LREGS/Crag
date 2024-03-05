package store

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lregs/Crag/env"
)

type SqlStore struct {
	masterX *sql.DB
	Stores  SqlStoreStores
}

type SqlStoreStores struct {
	Crag CragStore
}

func NewSqlStore() (*SqlStore, error) {
	store := &SqlStore{}
	err := store.initConnection()
	if err != nil {
		return nil, errors.New("error getting connection")
	}

	return store, nil
}

func (ss SqlStore) initConnection() error {
	DBURL, err := env.DBString()
	if err != nil {
		fmt.Printf("error establishing db connection %s", err)
	}

	ss.masterX, err = sql.Open("postgres", DBURL)
	if err != nil {
		panic(err)
	}
	return err
}

func (ss *SqlStore) GetMasterX() *sql.DB {
	return ss.masterX
}
