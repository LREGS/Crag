package store

import (
	"reflect"
	"testing"

	models "github.com/lregs/Crag/models"
	log "github.com/sirupsen/logrus"
)

func CreateSqlStore(t *testing.T) *SqlStore {
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

	store := CreateSqlStore(t)
	// CragStore := store.Stores.CragStore

	t.Run("Testing add crag", func(t *testing.T) {
		//I dont have a type for climbs, forecast, or reports yet and we need to make
		//baby steps with out testing so they will just be null types atm
		crag := returnCrag()

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

		defer func() {
			dropTables(t)
		}()

	})

}

func TestGetCrag(t *testing.T) {
	MockStore := returnPrePopulatedMockStore(t, false, false)

	t.Run("Testing Get Crag", func(t *testing.T) {
		TestMock := &models.Crag{Id: 1}
		log.Infof("Creating Test Mock %+v", TestMock)
		Id := 1

		_, err := MockStore.Stores.CragStore.GetCrag(Id)
		if err != nil {
			t.Fatalf("could not store crag because of err: %s", err)
		}
	})

}

func TestUpdateCrag(t *testing.T) {
	t.Logf("Creating pre populated Mock Store")
	MockStore := returnPrePopulatedMockStore(t, false, false)

	t.Run("Testing Update Crag", func(t *testing.T) {

		crag := returnCrag()
		crag.Name = "Milestone"
		log.Infof("Crag name changed to %s", crag.Name)

		err := MockStore.Stores.CragStore.UpdateCragValue("Stanage", *crag)
		if err != nil {
			t.Fatalf("Update failed because of error: %s", err)
		}

		log.Infof("getting crag to verify update")
		currentCrag, err := MockStore.Stores.CragStore.GetCrag(1)
		if err != nil {
			t.Fatalf("Failed to get crag because of err: %s", err)
		}

		log.Infof("Crag name now %s, wanted: %s", currentCrag.Name, crag.Name)
		if currentCrag.Name != "Milestone" {
			t.Fatalf("The update name %s does not match %s", currentCrag.Name, crag.Name)
		}

	})
}
func TestDeleteCrag(t *testing.T) {
	MockStore := returnPrePopulatedMockStore(t, false, false)

	t.Run("Testing Delete Crag", func(t *testing.T) {
		id := 1

		err := MockStore.Stores.CragStore.DeleteCragByID(id)
		if err != nil {
			t.Fatalf("error deleting crag %s: ", err)
		}

		_, err = MockStore.Stores.CragStore.GetCrag(1)
		if err == nil {
			t.Fatalf("Crag still exists")
		}
	})
}

func returnCrag() *models.Crag {
	crag := &models.Crag{
		Id:        1,
		Name:      "Stanage",
		Latitude:  40.7128,
		Longitude: -74.0060,
	}
	return crag
}

func returnPrePopulatedMockStore(t *testing.T, climb bool, forecast bool) *SqlStore {

	//probably dont need options, or this function - we need to edit the docker test file
	//and update the db init to include these initial entries for the test database
	//and then make a more obvious const for each climb, crag, weatherreport so its easier to test

	store := CreateSqlStore(t)
	log.Infof("CreatedStore %+v", store)
	MockCrag := returnCrag()
	log.Infof("Created MockCrag %+v", MockCrag)
	err := store.Stores.CragStore.StoreCrag(MockCrag)
	if err != nil {
		t.Fatalf("Couldn't store crag because of this error: %s", err)
	}

	if climb != false {
		MockClimb := returnClimb()
		err = store.Stores.ClimbStore.StoreClimb(MockClimb)
		if err != nil {
			t.Fatalf("could not store climb because of this error: %s", err)
		}
	}

	if forecast != false {
		MockForecast := newForecast()
		_, err = store.Stores.ForecastStore.AddForecast(MockForecast)
		if err != nil {
			t.Fatalf("could not store forecast because of error: %s", err)
		}
	}

	return store

}

// query := "select name, latitude, longitude from crag where id = $1"
// err := db.QueryRow(query, MockCrag.Id).Scan(&TestMock.Name, &TestMock.Latitude, &TestMock.Longitude)
// if err != nil {
// 	log.Fatalf("could not select items from db because of error: %s", err)
// }
