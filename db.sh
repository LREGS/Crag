#usr/bin/bash

if [ -f .env ]; then 
    source .env 
else
    echo "Error: .env file not found"
    exit 1 
fi 

if createdb -U william -O william crags; then 
    echo "Database crags created"
else
    echo "failed to create db"
    exit 1
fi 

psql -U william -d crags <<EOF
-- Drop tables if they exist
DROP TABLE IF EXISTS forecast;
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
)

COPY forecast FROM 'forecast.csv' WITH (format csv, header);
EOF