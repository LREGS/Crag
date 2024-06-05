package met

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"testing"

	"github.com/lregs/Crag/models"
)

func MarshallTestData(t *testing.T) models.Forecast {
	jsonFile, err := os.Open("sampleData.json")
	if err != nil {
		t.Log(err)
	}
	defer jsonFile.Close()

	byteJson, err := io.ReadAll(jsonFile)
	if err != nil {
		t.Log(err)
	}

	var forecast models.Forecast

	if err := json.Unmarshal(byteJson, &forecast); err != nil {
		t.Log(err)
	}

	return forecast

}

func TestGetForecast(t *testing.T) {

	coords := []float64{50.374422, -4.153563}

	f, err := GetForecast(coords)
	if err != nil {
		t.Fatal(err)
	}

	if len(f.Features) == 0 {
		t.Fatal("forecast empty")
	}
}

func TestGetPayload(t *testing.T) {
	//at least we now we're getting a payload but not testing that its correct? DO I need too?

	t.Run("Testing Get Payload", func(t *testing.T) {

		log := NewLogger("dummy.txt")

		payload, err := GetPayload(log, []float64{53.120607133644576, -3.9983421531498133})
		if err != nil {
			t.Fatal(err)
		}

		fmt.Print(payload)

		if len(payload) > 1 != true {
			t.Fatalf("payload empty")
		}

	})

}

func TestRedisPayload(t *testing.T) {
	data := MarshallTestData(t)
	t.Run("Testing Redis Payload", func(t *testing.T) {
		log := NewLogger("dummy.txt")

		p, err := GetRedisPayload(log, data)
		if err != nil {
			t.Fatal(err)
		}
		keys := make([]string, len(p))

		i := 0
		for k := range p {
			keys[i] = k
			i++
		}

		t.Error(p)

	})
}

func NewLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("bad file")
	}
	return log.New(logfile, "[main]", log.Ldate|log.Ltime|log.Lshortfile)
}
