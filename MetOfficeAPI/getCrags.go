package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func readFile(n string) ([][]string, error) {

	file, err := os.Open(n)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	//	this should be [name], ["lat, long"] which will need to be split and cast itself later
	res := [][]string{}

	for scanner.Scan() {
		splitLine := strings.Split(scanner.Text(), " ")

		coords := strings.Split(splitLine[1], ",")

		res := append(res, []string{splitLine[0], coords[0], coords[1]})
	}

	if err = scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return res, nil
}

func worker() {

}
