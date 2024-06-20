package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

//some of the names aren't very idiomatic - I know its a metX because this is the met package

func main() {

	loger := NewLogger("log.txt")

	loger.Println("started")

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
	store := NewMetStore(rc, loger)

	// lastUpdate, err := store.GetLastUpdate(loger)
	// if err != nil {
	// 	log.Println(err)
	// }
	// switch err {
	// case ErrorRedis:
	// 	//because api updates every hour the scheduler will shedule an immediate trigger to fetch the data and store it
	// 	go scheduler.startSchedule(loger, api, store, (time.Now().Add(-3 * time.Hour)))

	// 	if err != nil {
	// 		log.Println(err)
	// 	}

	// }

	// at the moment the scheduler is starting from every hour from the time the app is started.
	// it needs to be an hour from when it was last updated

	go scheduler.startSchedule(loger, api, store, time.Now())

	tmpl := template.Must(template.ParseFiles("./templates/main.html"))

	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data, err := store.GetForecastTotals()
		if err != nil {
			log.Fatal("failed to get data")
		}

		// d := []ForecastTotals{}
		// for _, key := range data {

		// }

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

func ExecuteRefreshProcess(log *log.Logger, api *MetOfficeAPI, store *MetStore) time.Time {

	//	need to now figure out how to properly store multiple at a time

	crags, err := ReadFile("crags.txt")

	f, err := api.GetForecast([]float64{53.12266792026611, -3.9965825915253648})
	if err != nil {
		log.Printf("couldn't fetch met data %s", err)
	}

	p, err := api.GetPayload(log, f)
	if err != nil {
		log.Printf("failed creating payload %s", err)
	}

	if err := store.ForecastTotals(context.Background(), p); err != nil {
		log.Printf("failed storing forecast totals, %s", err)
	}

	t, err := store.GetLastUpdate(log)
	if err != nil {
		log.Printf("failed refresh: couldn't get last update value")
	}

	return t

}
