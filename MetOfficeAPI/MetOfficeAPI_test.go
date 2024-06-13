package main

import (
	"encoding/json"
	"io"
	"os"
	"testing"
)

func GetTestData(t *testing.T) Forecast {
	jsonFile, err := os.Open("test/sampleData.json")
	if err != nil {
		t.Log(err)
	}
	defer jsonFile.Close()

	byteJson, err := io.ReadAll(jsonFile)
	if err != nil {
		t.Log(err)
	}

	var forecast Forecast

	if err := json.Unmarshal(byteJson, &forecast); err != nil {
		t.Log(err)
	}

	return forecast
}

func TestFindWindows(t *testing.T) {

	api := NewMetAPI(os.Getenv("apikey"))
	testData := GetTestData(t)

	t.Run("Testing Find Windows", func(t *testing.T) {

		l := NewLogger("test.txt")

		windows := api.FindWindows(l, testData)

		t.Error(windows)

	})

}
