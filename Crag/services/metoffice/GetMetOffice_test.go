package met

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/lregs/Crag/logger"
	"github.com/redis/go-redis/v9"
)

func MarshallTestData(t *testing.T) Forecast {
	jsonFile, err := os.Open("sampleData.json")
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

// TODO: Remake they're using up the api calls
// func TestGetForecast(t *testing.T) {

// 	coords := []float64{50.374422, -4.153563}

// 	f, err := GetForecast(coords)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	t.Log(f)

// 	if len(f.Features) == 0 {
// 		t.Fatal("forecast empty")
// 	}
// }

// func TestGetPayload(t *testing.T) {
// 	//at least we now we're getting a payload but not testing that its correct? DO I need too?

// 	t.Run("Testing Get Payload", func(t *testing.T) {

// 		log := NewLogger("dummy.txt")

// 		payload, err := GetPayload(log, []float64{53.120607133644576, -3.9983421531498133})
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		fmt.Print(payload)

// 		if len(payload) > 1 != true {
// 			t.Fatalf("payload empty")
// 		}

// 	})

// }

// func TestRedisPayload(t *testing.T) {
// 	data := MarshallTestData(t)
// 	t.Run("Testing Redis Payload", func(t *testing.T) {
// 		log := NewLogger("dummy.txt")

// 		p, err := TotalsByDay(log, data)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		keys := make([]string, len(p))

// 		i := 0
// 		for k := range p {
// 			keys[i] = k
// 			i++
// 		}

// 		t.Error(keys)

// 		for _, v := range p {
// 			t.Error(v)
// 		}

// 		t.Error(time.Now())

// 	})
// }

func TestStoreData(t *testing.T) {
	log := logger.NewLogger("MetOfficeTestLog.txt")
	data := MarshallTestData(t)
	payload, err := GetPayload(log, data)
	if err != nil {
		t.Fatalf("err getting payload %s", err)
	}

	rc := redis.NewClient(&redis.Options{
		Addr:     "redis-19441.c233.eu-west-1-1.ec2.redns.redis-cloud.com:19441",
		Password: "N9jHgekt2GxfqkHpQtNHL7jmwUCkq3zA",
		DB:       0,
	})

	t.Run("Testing Store Redis Data", func(t *testing.T) {

		if err := StoreData(log, context.Background(), rc, payload); err != nil {
			t.Errorf("err storing %s", err)
		}

		s, err := rc.Get(context.Background(), "LastUpdated").Result()
		if err != nil {
			t.Errorf("err getting data, %s", err)
		}

		// sconv, _ := strconv.ParseInt(s, 10, 64)

		// time := time.Unix(sconv, 0)

		t.Error(s)

	})

}
