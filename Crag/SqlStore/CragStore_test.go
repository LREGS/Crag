package store

//this isnt a mock store, this is the store we're testing lol

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"

	models "github.com/lregs/Crag/models"
	log "github.com/sirupsen/logrus"
)

var initDb = `DROP TABLE IF EXISTS forecast;
DROP TABLE IF EXISTS report;
DROP TABLE IF EXISTS climb;
DROP TABLE IF EXISTS crag;

-- Create tables
CREATE TABLE crag (
	Id SERIAL PRIMARY KEY, 
	Name TEXT UNIQUE, 
	Latitude DOUBLE PRECISION,
	Longitude DOUBLE PRECISION
);

CREATE TABLE climb (
	Id SERIAL PRIMARY KEY,
	Name VARCHAR(255) UNIQUE,
	Grade VARCHAR(255),
	CragID INTEGER REFERENCES crag(Id)
);

CREATE TABLE report (
	Id SERIAL PRIMARY KEY, 
	Content VARCHAR(255),
	Author VARCHAR(255),
	CragID INTEGER REFERENCES crag(Id)
);

CREATE TABLE forecast (
	Id SERIAL PRIMARY KEY, 
	Time VARCHAR(255) UNIQUE,
	ScreenTemperature DOUBLE PRECISION,
	FeelsLikeTemp DOUBLE PRECISION, 
	WindSpeed DOUBLE PRECISION,
	WindDirection DOUBLE PRECISION,
	totalPrecipitation DOUBLE PRECISION,
	ProbofPrecipitation INT,
	Latitude DOUBLE PRECISION,
	Longitude DOUBLE PRECISION
);`

func CreateSqlStore(t *testing.T) *SqlStore {

	conn, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	_, err = conn.Exec(context.Background(), initDb)
	if err != nil {
		panic(err)
	}

	store, err := NewSqlStore(&StoreConfig{DbConnection: conn})
	if err != nil {
		t.Fatalf("error creating store: %s", err)
	}

	// log.Info("Creating tables")
	// err = CreateTables(t)
	// if err != nil {
	// 	t.Fatalf("could not create tables because of %s", err)
	// }
	// log.Infof("Tables created")
	return store
}

func TestAddCrag(t *testing.T) {

	CragStore := NewCragStore(CreateSqlStore(t))
	// CragStore := store.Stores.CragStore

	t.Run("Testing add crag", func(t *testing.T) {
		crag := models.CragPayload{Name: "crag", Latitude: 0.13435463, Longitude: 0.5435436}

		d, err := CragStore.StoreCrag(context.Background(), crag)
		if err != nil {
			log.Fatalf("was not able to store Crag because of err: %s", err)
		}

		assert.Equal(t, crag.Name, d.Name)

		defer func() {
			dropTables(t)
		}()

	})

	t.Run("testing adding invalid data", func(t *testing.T) {
		invalidCrag := models.CragPayload{
			Name:      "",
			Latitude:  0.0,
			Longitude: 0.0,
		}

		_, err := CragStore.StoreCrag(context.Background(), invalidCrag)
		if err == nil {
			t.Fatal("no error returned when passing invalid crag to store")
		}

	})

}

func TestGetCrag(t *testing.T) {
	MockStore := returnPrePopulatedMockStore(t, false, false)

	t.Run("Testing Get Crag", func(t *testing.T) {
		Id := 1

		data, err := MockStore.Stores.CragStore.GetCrag(context.Background(), Id)
		if err != nil {
			t.Fatalf("could not store crag because of err: %s", err)
		}

		testData := returnCrag()

		assert.Equal(t, data.Name, testData.Name)
		assert.Equal(t, data.Latitude, testData.Latitude)
		assert.Equal(t, data.Longitude, testData.Longitude)

	})

	t.Run("Testing Invalid Crag Id", func(t *testing.T) {
		Id := -99

		_, err := MockStore.Stores.CragStore.GetCrag(context.Background(), Id)
		if err == nil {
			t.Fatal("store accepted invalid id")
		}
	})

}

func TestUpateCrag(t *testing.T) {
	MockStore := returnPrePopulatedMockStore(t, false, false)

	t.Run("Update Crag", func(t *testing.T) {

		crag := models.Crag{Id: 1, Name: "dank", Latitude: 1.1, Longitude: 2.2}

		data, err := MockStore.Stores.CragStore.UpdateCrag(context.Background(), crag)
		if err != nil {
			t.Fatalf("failed storing crag: %s", err)
		}
		assert.Equal(t, crag.Name, data.Name)

	})

	t.Run("Send Invalid Crag", func(t *testing.T) {
		crag := models.Crag{}

		_, err := MockStore.Stores.CragStore.UpdateCrag(context.Background(), crag)
		if err == nil {
			t.Fatalf("store accepted invalid crag")
		}
	})
}

func TestDeleteCrag(t *testing.T) {
	MockStore := returnPrePopulatedMockStore(t, false, false)

	t.Run("Testing Delete Crag", func(t *testing.T) {
		id := 1

		_, err := MockStore.Stores.CragStore.DeleteCragByID(context.Background(), id)
		if err != nil {
			t.Fatalf("error deleting crag %s: ", err)
		}

		_, err = MockStore.Stores.CragStore.GetCrag(context.Background(), 1)
		if err == nil {
			t.Fatalf("Crag still exists")
		}
	})
}

func returnCrag() models.CragPayload {
	crag := models.CragPayload{
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

	MockCrag := returnCrag()
	_, err := store.Stores.CragStore.StoreCrag(context.Background(), MockCrag)
	if err != nil {
		t.Fatalf("Couldn't store crag because of this error: %s", err)
	}

	if climb != false {
		MockClimb := returnPayload()
		_, err = store.Stores.ClimbStore.StoreClimb(context.Background(), MockClimb)
		if err != nil {
			t.Fatalf("could not store climb because of this error: %s", err)
		}
	}

	if forecast != false {
		MockForecast := testPayload()
		_, err = store.Stores.ForecastStore.StoreForecast(context.Background(), MockForecast)
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
