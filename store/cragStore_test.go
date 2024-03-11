package store

import (
	"reflect"
	"testing"

	models "github.com/lregs/Crag/models"
	log "github.com/sirupsen/logrus"
)

func CreateCragStore(t *testing.T) *SqlStore {
	store, err := NewSqlStore(&StoreConfig{dbConnection: db})
	if err != nil {
		t.Fatalf("error creating store: %s", err)
	}
	log.Info("Creating tables")
	err = CreateTables(t)
	if err != nil {
		t.Fatalf("could not create tables because of %s", err)
	}
	log.Infof("Tables created")
	return store
}

func TestAddCrag(t *testing.T) {

	store := CreateCragStore(t)
	// CragStore := store.Stores.CragStore

	t.Run("Testing add crag", func(t *testing.T) {
		//I dont have a type for climbs, forecast, or reports yet and we need to make
		//baby steps with out testing so they will just be null types atm
		crag := &models.Crag{
			Id:        1,
			Name:      "Stanage",
			Latitude:  40.7128,
			Longitude: -74.0060,
		}

		log.Infof("store: %+v", store)
		log.Infof("Cragstore = %+v", store.Stores.CragStore)

		err := store.Stores.CragStore.StoreCrag(crag)
		if err != nil {
			log.Fatalf("was not able to store Crag because of err: %s", err)
		}

		testData := &models.Crag{Id: 1}

		query := "select name, latitude, longitude from crag where id = $1"
		err = db.QueryRow(query, crag.Id).Scan(&testData.Name, &testData.Latitude, &testData.Longitude)
		if err != nil {
			log.Fatalf("wasn't able to retrieve data from db: %s", err)
		}

		if !reflect.DeepEqual(testData, crag) {
			log.Fatalf("the returned data from the db does not match that of the inputted data")
		}

	})

}
