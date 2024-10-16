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
	log.Print(store.GetLastUpdate())
	go UpdateForecasts(ctx, store.GetLastUpdate(), api, store, errs)

	// not sure if we want to handle our error this way but before we were blocking the whole app waiting
	// for error that maybe never occured
	go func() {
		for {
			select {
			case err := <-errs:
				if err != nil {
					log.Print("Received error:", err)
				}
			}
		}
	}()

	tmpl := template.Must(template.ParseFiles("./templates/main.html"))

	router := http.NewServeMux()

	// rat a tat tat
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Print("handling")
		data, err := store.Get()
		log.Print(data)
		if err != nil {
			log.Print("failed to get data", err)
		}

		if err := tmpl.ExecuteTemplate(w, "main", data.ForecastTotals["15"]); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	if err := http.ListenAndServe(":8181", router); err != nil {
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

// I dont think this is correct because it needs to be updating on the hour, not an hour from when the prog was run
func UpdateForecasts(ctx context.Context, lastUpdate time.Time, api MetAPI, store *MetStore, errs chan<- error) error {

	log.Print("Starting forecast update process")

	var sem = semaphore.NewWeighted(int64(api.GetRateLimit() * 0.8))
	wg := sync.WaitGroup{}

	// This loop runs indefinitely, updating once per hour
	for {
		// Check for context cancellation to stop the loop
		if err := ctx.Err(); err != nil {
			return err
		}

		// Only proceed if it's time for the next update (once per hour)
		if time.Since(lastUpdate) >= time.Hour {
			log.Print("Updating forecasts for all crags")

			for _, crag := range crags {
				wg.Add(1)
				// Acquire a semaphore slot to control the rate limit
				if err := sem.Acquire(ctx, 1); err != nil {
					log.Print(err)
					errs <- err
					wg.Done()
					continue
				}

				// Process each crag concurrently using goroutines
				go func(crag Crag) {
					defer sem.Release(1) // Ensure semaphore is released
					defer wg.Done()      // Ensure wg.Done() is called

					url := api.CreateURL([]float64{crag.Latitude, crag.Longitude})
					log.Print("Fetching forecast for: ", crag.Name, url)

					// Call the API to get the forecast
					f, err := api.GetForecast(url)
					if err != nil {
						log.Print("Error fetching forecast for", crag.Name, ":", err)
						errs <- err
						return
					}

					timeSeries := f.Features[0].Properties.TimeSeries
					if len(timeSeries) == 0 {
						log.Print("Forecast empty for", crag.Name)
						errs <- errors.New("forecast time series data is empty for " + crag.Name)
						return
					}

					// Add the forecast data to the store
					log.Print("Storing forecast for", crag.Name)
					if err := store.Add(ctx, crag.Name,
						ForecastPayload{
							LastModelRunTime: f.Features[0].Properties.ModelRunDate,
							ForecastTotals:   api.CalculateTotals(timeSeries),
						}); err != nil {
						log.Print("Error storing forecast for", crag.Name, ":", err)
						errs <- err
					}
				}(crag)
			}

			// Wait for all crags to be updated before proceeding
			wg.Wait()

			// Update lastUpdate to the current time
			lastUpdate = time.Now()

			// Log that the batch update is complete
			log.Print("Batch update complete. Waiting for the next hour.")

			// Wait for the next hour
			select {
			case <-time.After(time.Hour):
				// Continue to the next batch after an hour
				log.Print("Starting new batch after one hour.")
			case <-ctx.Done():
				// If context is canceled, exit the loop gracefully
				return ctx.Err()
			}
		} else {
			// If less than an hour has passed, wait for the next full hour
			sleepDuration := time.Hour - time.Since(lastUpdate)
			log.Printf("Waiting for the next update window in %v", sleepDuration)

			select {
			case <-time.After(sleepDuration):
				// Proceed after waiting for the remaining time to complete the hour
				log.Print("Resuming after waiting.")
			case <-ctx.Done():
				// If context is canceled, exit the loop gracefully
				return ctx.Err()
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
