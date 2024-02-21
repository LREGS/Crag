package Store

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

type SqlStore struct {
	//rrCounter int 64
	//srCounter int 64

	masterX *sql.DB
	//Replica
	//Context
	//RWMutex

	Stores SqlStoreStores
}

type SqlStoreStores struct {
	Forecast ForecastStore
	Crag     CragStore
}

func New() (*SqlStore, error) {
	var err error
	store := &SqlStore{}

	err = store.initConnection()
	if err != nil {
		return nil, errors.New("error starting connection")
	}

	// inits the Forecast store
	store.stores.Forecast = NewSqlForecastStore(store)

	return store, nil
}

func (ss *SqlStore) initConnection() error {
	// DbUsername := os.Getenv("DB_USERNAME")
	// DbPassword := os.Getenv("DB_PASSWORD")
	var err error

	ss.masterX, err = sql.Open("postgres", fmt.Sprintf("user=william dbname=crag password=1 sslmode=disable"))
	if err != nil {
		panic(err)
	}
	return err

}

func (ss *SqlStore) GetMasterX() *sql.DB {
	return ss.masterX
}

func (ss *SqlStore) Forecast() ForecastStore {
	return ss.stores.Forecast
}
