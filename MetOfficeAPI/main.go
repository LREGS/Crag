package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"sync"
	"text/template"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/semaphore"
)

func main() {
	log := NewLogger("log.txt")

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("couldn't load environment variables")
	}

	rc := redis.NewClient(&redis.Options{
		Addr:     "redis-13149.c85.us-east-1-2.ec2.redns.redis-cloud.com:13149",
		Password: os.Getenv("redis"),
		DB:       0,
	})

	api := NewMetAPI(os.Getenv("apikey"), log)
	store := NewMetStore(rc, log)

	ctx := context.Background()
	errs := make(chan error, 1)
	go UpdateForecasts(ctx, store.GetLastUpdate(), api, store, errs)

	select {
	case err := <-errs:
		log.Print(err)
	}

	tmpl := template.Must(template.ParseFiles("./templates/main.html"))

	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data, err := store.Get()
		if err != nil {
			log.Fatalf("failed to get data %s", err)
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

func UpdateForecasts(ctx context.Context, lastUpdate time.Time, api *MetOfficeAPI, store *MetStore, errs chan<- error) error {

	log.Print("updating")
	var sem = semaphore.NewWeighted(int64(api.ratelimiter.Limit() * 0.8))
	wg := sync.WaitGroup{}

	for {
		if time.Since(lastUpdate) > time.Hour {
			for _, crag := range crags {
				log.Print(crag)
				wg.Add(1)
				if err := sem.Acquire(ctx, 1); err != nil {
					log.Print(err)
					errs <- err
				}
				go func() {
					log.Print(api.CreateURL([]float64{crag.Latitude, crag.Longitude}))
					// DO i want to store these or make them as I am?
					f, err := api.GetForecast(api.CreateURL([]float64{crag.Latitude, crag.Longitude}))
					log.Print(f)
					if err != nil {
						log.Print(err)
						errs <- err

					}
					// log.Print("adding payload")
					// p, err := api.GetPayload(f.Features[0].Properties.ModelRunDate, f.Features[0].Properties.TimeSeries)
					// if err != nil {
					// 	errs <- err
					// }

					if len(f.Features[0].Properties.TimeSeries) == 0 {
						log.Print("forecast empty")
						errs <- errors.New("forecast time series data is empty")

					}
					log.Print("adding to store")
					if err := store.Add(ctx, crag.Name,
						ForecastPayload{LastModelRunTime: f.Features[0].Properties.ModelRunDate, ForecastTotals: api.CalculateTotals(f.Features[0].Properties.TimeSeries)}); err != nil {
						log.Print(err)
						errs <- err

					}
					wg.Done()
				}()
			}
		} else {
			c := time.Tick(time.Duration(60-(time.Now().Minute())) * time.Minute)
			log.Print(c)
			for _ = range c {
				continue
			}
		}

	}

}

var crags = []Crag{
	{"cromlech", 53.08977582752912, -4.0494354521953895},
	{"beddgelert", 53.01401346937128, -4.1086367318613055},
	{"gwynant", 53.04567339439013, -4.021447439922229},
	{"blaenau", 52.99729599359651, -3.9578734953238475},
	{"crafnant", 52.99729599359651, -3.9578734953238475},
	{"cwellyn", 53.07568570139747, -4.148701296939546},
	{"orme", 53.33236585445307, -3.8311890286450865},
	{"Penmaenbach", 53.285, -3.8684},
	{"ysgo", 52.80614677538971, -4.656639551730091},
	{"tremadoch", 52.94008535336955, -4.140997768369204},
	{"rhiwGoch", 53.09199013529737, -3.803795346023221},
	{"portland", 50.545900401402854, -2.438814867485551},
	{"cuckooRock", 50.545900401402854, -2.438814867485551},
	{"MountSionEast", 50.545900401402854, -2.438814867485551},
	{"Froggatt", 53.2942103060766, -1.6201285054945418},
}
