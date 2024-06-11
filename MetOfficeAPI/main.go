package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

//some of the names aren't very idiomatic - I know its a metX because this is the met package

func main() {

	log := NewLogger("log.txt")

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("couldn't load environment variables")
	}

	rc := redis.NewClient(&redis.Options{
		Addr:     "redis-19441.c233.eu-west-1-1.ec2.redns.redis-cloud.com:19441",
		Password: os.Getenv("redis"),
		DB:       0,
	})

	scheduler := Scheduler{}

	api := NewMetAPI(os.Getenv("apikey"))
	store := NewMetStore(rc, log)

	lastUpdate, err := store.GetLastUpdate(log)
	switch err {
	case ErrorRedis:
		//because api updates every hour the scheduler will shedule an immediate trigger to fetch the data and store it
		go scheduler.startSchedule(log, api, store, (time.Now().Add(-3 * time.Hour)))

		if err != nil {
			log.Println(err)
		}

	}
	go scheduler.startSchedule(log, api, store, lastUpdate)
}

func Str2Time(timeString string) (time.Time, error) {
	parsedTime, err := time.Parse("2006-01-02T15:04Z07:00", timeString)
	if err != nil {
		log.Printf("failed parsing time during cache update")
		return time.Time{}, err
	}
	return parsedTime, nil
}

func ExecuteRefreshProcess(log *log.Logger, api *MetOfficeAPI, store *MetStore) {

	f, err := api.GetForecast([]float64{53.12266792026611, -3.9965825915253648})
	if err != nil {
		log.Printf("couldn't fetch met data %s", err)
	}

	p, err := api.GetPayload(log, f)
	if err != nil {
		log.Printf("failed creating payload %s", err)
	}

	if err := store.StoreForecastTotals(context.Background(), p); err != nil {
		log.Printf("failed storing forecast totals, %s", err)
	}

}
