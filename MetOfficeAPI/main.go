package main

import (
	"log"
	"net/http"
	"os"
	"sync"
	"text/template"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

// better di for log pls

//some of the names aren't very idiomatic - I know its a metX because this is the met package

func main() {
	// loger??????
	loger := NewLogger("log.txt")

	loger.Println("started")

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("couldn't load environment variables")
	}

	rc := redis.NewClient(&redis.Options{
		Addr:     "redis-13149.c85.us-east-1-2.ec2.redns.redis-cloud.com:13149",
		Password: os.Getenv("redis"),
		DB:       0,
	})

	// scheduler := Scheduler{}

	api := NewMetAPI(os.Getenv("apikey"), loger)
	store := NewMetStore(rc, loger)

	// at the moment the scheduler is starting from every hour from the time the app is started.
	// it needs to be an hour from when it was last updated

	// go scheduler.startSchedule(loger, api, store, time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC))
	go ScheduleMetOffice(loger, api, store)

	tmpl := template.Must(template.ParseFiles("./templates/main.html"))

	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data, err := store.GetTotals()
		if err != nil {
			log.Fatal("failed to get data %s", err)
		}

		if err := tmpl.ExecuteTemplate(w, "main", *data["17"]); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	if err := http.ListenAndServe(":6968", router); err != nil {
		panic(err)
	}

}

func Str2Time(timeString string) (time.Time, error) {
	parsedTime, err := time.Parse("2006-01-02T15:04Z07:00", timeString)
	if err != nil {
		log.Printf("failed parsing time during cache update")
		return time.Time{}, err
	}
	return parsedTime, nil
}

type forecastMutex struct {
	mu sync.Mutex
	//sorted by crag name and complete forecast
	forecastResults map[string]Forecast
}
