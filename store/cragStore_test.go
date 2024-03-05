package store

import (
	"testing"

	models "github.com/lregs/Crag/models"
	log "github.com/sirupsen/logrus"
)

func CreateStore(t *testing.T) CragStore {
	store, err := NewSqlStore()
	if err != nil {
		log.Fatalf("error creating store: %s", err)
	}
	return store.Stores.CragStore
}

func TestAddCrag(t *testing.T) {

	store := CreateStore

	t.Run("Testing add crag", func(t *testing.T) {
		//I dont have a type for climbs, forecast, or reports yet and we need to make
		//baby steps with out testing so they will just be null types atm
		crag := models.Crag{
			Id:        1,
			Name:      "Stanage",
			Latitude:  40.7128,
			Longitude: -74.0060,
		}

		err := store.StoreCrag(crag)
		if err != nil {
			log.Fatalf("was not about to store Crag because of err: %s", err)
		}

		var testData models.Crag
		testData.Id = 1

		query := "select name, latitude, longtitude from crag where id = $1"
		err = db.QueryRow(query, crag.Id).Scan(&testData.Name, &testData.Latitude, &testData.Longitude)
		if err != nil {
			log.Fatalf("wasn't able to retrieve data from db: %s", err)
		}

		if testData != crag {
			log.Fatalf("the returned data from the db does not match that of the inputted data")
		}

	})

}
