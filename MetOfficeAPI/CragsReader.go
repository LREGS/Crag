package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

func ReadCragsTxt(n string) ([]Crag, error) {

	file, err := os.Open(n)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	//	this should be [name, "lat, long"] which will need to be split and cast itself later
	res := []Crag{}

	for scanner.Scan() {

		name := strings.Split(scanner.Text(), " ")
		coords := strings.Split(name[1], ",")
		lat, err := strconv.ParseFloat(coords[0], 64)
		if err != nil {
			return []Crag{}, err
		}
		lon, err := strconv.ParseFloat(coords[1], 64)
		if err != nil {
			return []Crag{}, err
		}

	

		res = append(res, Crag{name: name[0], Latitude: lat, Longitude: lon,)
	}

	if err = scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return res, nil
}

func worker() {

}
