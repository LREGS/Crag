package store

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	log "github.com/sirupsen/logrus"
)

var db *sql.DB

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "11",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=user_name",
			"POSTGRES_DB=dbname",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://user_name:secret@%s/dbname?sslmode=disable", hostAndPort)

	log.Println("Connecting to database on url: ", databaseUrl)

	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		db, err = sql.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	//Run tests
	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestDBCreation(t *testing.T) {
	query := `-- Drop tables if they exist
	DROP TABLE IF EXISTS forecast;
	DROP TABLE IF EXISTS report;
	DROP TABLE IF EXISTS climb;
	DROP TABLE IF EXISTS crag;
	
	-- Create tables
	CREATE TABLE crag (
		Id SERIAL PRIMARY KEY, 
		Name TEXT, 
		Latitude DOUBLE PRECISION,
		Longitude DOUBLE PRECISION
	);
	
	CREATE TABLE climb (
		Id SERIAL PRIMARY KEY,
		Name VARCHAR(255),
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
		Time VARCHAR(255),
		ScreenTemperature DOUBLE PRECISION,
		FeelsLikeTemp DOUBLE PRECISION, 
		WindSpeed DOUBLE PRECISION,
		WindDirection DOUBLE PRECISION,
		totalPrecipitation DOUBLE PRECISION,
		ProbofPrecipitation INT,
		Latitude DOUBLE PRECISION,
		Longitude DOUBLE PRECISION
	);`

	db.Exec(query)

	assertTables(t)

}

func assertTables(t *testing.T) {
	tables := map[string]bool{
		"crag":     true,
		"climb":    true,
		"report":   true,
		"forecast": true,
	}

	t.Run("Testing if tables correctly added", func(t *testing.T) {

		Query := `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public' AND table_type = 'BASE TABLE' `

		rows, err := db.Query(Query)
		if err != nil {
			log.Fatalf("wasnt able to make queary because of %s:", err)
		}

		for rows.Next() {
			var currentTable string
			if err := rows.Scan(&currentTable); err != nil {
				log.Fatal(err)
			}
			_, ok := tables[currentTable]

			if ok {
				continue
			} else {
				log.Fatalf("%s has been incorrectly added to the db, or doesnt exist in the test set", currentTable)
			}
		}
	})
}
