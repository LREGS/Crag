package util

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

type TestDBConfig struct {
	Crag     bool
	Forecast bool
	Climb    bool
}

func ReturnTestDBConnection(config *TestDBConfig) (*sql.DB, error) {
	var DB *sql.DB

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
		DB, err = sql.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		return DB.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	err = CreateTables(DB)
	if err != nil {
		return nil, err
	}

	if config.Climb != false {
		_, err := DB.Exec("INSERT INTO Climb (Name, Grade, CragID) VALUES ($1, $2, $3)", "Climb A", "5.9", 1)
		if err != nil {
			return nil, err
		}
	}
	if config.Crag != false {
		_, err := DB.Exec("INSERT INTO Crag (Name, Latitude, Longitude) VALUES ($1, $2, $3)", "Example Crag 1", 40.7128, -74.0060)
		if err != nil {
			return nil, err
		}
	}
	if config.Crag != false {
		_, err := DB.Exec("INSERT INTO DBForecast (Time, ScreenTemperature, FeelsLikeTemp, WindSpeed, WindDirection, TotalPrecipAmount, ProbOfPrecipitation, Latitude, Longitude, CragId) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", "2024-03-16 08:00:00", 20.5, 18.2, 10.3, 120.5, 0.0, 0.0, 40.7128, -74.0060, 1)
		if err != nil {
			return nil, err
		}
	}

	return DB, nil
}

func CreateTables(DB *sql.DB) error {
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
		Longitude DOUBLE PRECISION,
		CragID INTEGER REFERENCES crag(Id)

	);`

	_, err := DB.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
