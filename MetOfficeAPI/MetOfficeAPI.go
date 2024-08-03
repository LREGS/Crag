package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	rate "golang.org/x/time/rate"
)

// type MetOfficeAPI interface{

// }

const baseHourlyURL string = "https://data.hub.api.metoffice.gov.uk/sitespecific/v0/point/hourly?"

type api interface {
	Call(url string, headers http.Header) (*http.Response, error)
}

type Api struct {
	client http.Client
}

func NewApi() *Api {
	return &Api{client: http.Client{}}
}

func (a *Api) Call(url string, header http.Header) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header = header

	res, err := a.client.Do(req)
	if err != nil {
		return nil, err

	}

	return res, nil

}

type MetOfficeAPI struct {
	log     *log.Logger
	BaseURL string
	// I dont understand having this and not just returning from env?
	APIKey      string
	client      api
	store       Store
	ratelimiter *rate.Limiter
}

func NewMetAPI(apikey string, log *log.Logger) *MetOfficeAPI {
	return &MetOfficeAPI{
		log:         log,
		BaseURL:     baseHourlyURL,
		APIKey:      apikey,
		client:      NewApi(),
		ratelimiter: rate.NewLimiter(rate.Every(time.Minute/100), 100),
	}
}

// Returns the hourly forecast for a 72hour period from the met office data hub api, hourly
func (mAPI *MetOfficeAPI) GetForecast(url string) (Forecast, error) {

	var forecast Forecast

	res, err := mAPI.client.Call(url, mAPI.GetHeaders())
	if err != nil {
		return Forecast{}, nil
	}

	if res.StatusCode != 200 {
		return forecast, fmt.Errorf("code %d", res.StatusCode)
	}

	if err := json.NewDecoder(res.Body).Decode(&forecast); err != nil {
		return forecast, err
	}

	return forecast, nil

}

func (mAPI *MetOfficeAPI) CreateURL(coords []float64) string {
	return fmt.Sprintf("%slatitude=%f&longitude=%f", mAPI.BaseURL, coords[0], coords[1])
}

func (mAPI *MetOfficeAPI) GetHeaders() http.Header {
	return http.Header{

		"apikey": {mAPI.APIKey},
		"accept": {"application/json"},
	}
}

func (mAPI *MetOfficeAPI) GetPayload(forecast Forecast) (ForecastPayload, error) {

	//1 hourly spot has three days worth of data. This function provides the totals for
	//all three days

	data := forecast.Features[0].Properties.TimeSeries

	if len(data) == 0 {
		return ForecastPayload{}, errors.New("Forecast provided is empty")
	}

	totals := map[string]*ForecastTotals{}

	//we can try and calculate a running average using the mean calculator and then adding that value as the avg instead of current
	//clumsyness

	for _, val := range data {

		_, ok := totals[val.Time[8:10]]
		if !ok {
			totals[val.Time[8:10]] = &ForecastTotals{}
			entry := totals[val.Time[8:10]]
			entry.AvgTemp += val.ScreenTemperature
			entry.AvgWindSpeed += val.WindSpeed10m
			entry.HighestTemp = val.MaxScreenAirTemp
			entry.LowestTemp = val.MinScreenAirTemp
			//this is just assiging it as the first given value
			entry.WindDirection = val.WindDirectionFrom10m
			entry.TotalPrecip += val.TotalPrecipAmount
		} else {
			entry := totals[val.Time[8:10]]
			entry.AvgTemp += val.ScreenTemperature
			entry.AvgWindSpeed += val.WindSpeed10m
			entry.TotalPrecip += val.TotalPrecipAmount
		}

	}

	for day, entry := range totals {
		entry.AvgTemp = float64((entry.AvgTemp / float64((len(data) / 3))))
		entry.AvgWindSpeed = float64((entry.AvgWindSpeed / float64((len(data) / 3))))
		totals[day] = entry

	}

	return ForecastPayload{LastModelRunTime: forecast.Features[0].Properties.ModelRunDate, ForecastTotals: totals, Windows: mAPI.FindWindows(forecast)}, nil

}

func (mAPI *MetOfficeAPI) FindWindows(forecast Forecast) [][]time.Time {

	//probably should be returning and handling our errors in this

	var startOfWindow string
	var endOfWindow string

	windows := [][]time.Time{}
	for i := 0; i < len(forecast.Features[0].Properties.TimeSeries); i++ {

		currForecast := forecast.Features[0].Properties.TimeSeries[i]

		if currForecast.TotalPrecipAmount != 0.00 {
			if startOfWindow != "" {
				endOfWindow = currForecast.Time
				startTime, err := Str2Time(startOfWindow)
				if err != nil {
					log.Printf("faield converting start string %s", err)
					continue
				}
				endTime, err := Str2Time(endOfWindow)
				if err != nil {
					log.Printf("faield converting end string %s", err)
					continue
				}
				windows = append(windows, []time.Time{startTime, endTime})
				//	reset pointer
				// log.Println(startOfWindow, endOfWindow)
				startOfWindow = ""

			}

			continue
		}

		if startOfWindow == "" {
			startOfWindow = currForecast.Time
		}

		if i == (len(forecast.Features[0].Properties.TimeSeries) - 1) {
			if currForecast.TotalPrecipAmount != 0.00 && startOfWindow != "" {
				endOfWindow = forecast.Features[0].Properties.TimeSeries[i-1].Time
				startTime, err := Str2Time(startOfWindow)
				if err != nil {
					log.Printf("faield converting start string %s", err)
					continue
				}
				endTime, err := Str2Time(endOfWindow)
				if err != nil {
					log.Printf("faield converting end string %s", err)
					continue
				}
				windows = append(windows, []time.Time{startTime, endTime})
			}

			if startOfWindow != "" {
				endOfWindow = forecast.Features[0].Properties.TimeSeries[i-1].Time
				startTime, err := Str2Time(startOfWindow)
				if err != nil {
					log.Printf("faield converting start string %s", err)
					continue
				}
				endTime, err := Str2Time(endOfWindow)
				if err != nil {
					log.Printf("faield converting end string %s", err)
					continue
				}
				windows = append(windows, []time.Time{startTime, endTime})
			}

		}

		// log.Println(startOfWindow, endOfWindow)

	}

	return windows
}

// func (mAPI *MetOfficeAPI) ScheduleMetOffice(lastUpdate time.Time) {

// 	// if err := store.Flush(); err != nil {
// 	// 	log.Printf("couldn't update cache because of error whilst flushing %s", err)
// 	// }

// 	for {

// 		// not handling if there was no last update from redis
// 		// t := mAPI.store.GetLastUpdate()
// 		// if err != nil {
// 		// 	// do we need to be sending errors to a channel or something

// 		// 	// yes when this doesnt parse correctly we're calling the api infinitely
// 		// 	mAPI.log.Println(err)
// 		// }

// 		// log.Printf("last updated %s", t.String())

// 		if time.Since(lastUpdate) > time.Hour {
// 			mAPI.log.Printf("updating now")
// 			mAPI.updater(mAPI.log, mAPI.store)
// 		} else {
// 			log.Print("will update later")
// 			c := time.Tick(time.Duration(60-(time.Now().Minute())) * time.Minute)
// 			for next := range c {
// 				mAPI.log.Print(next)
// 				mAPI.updater(mAPI.log, mAPI.store)
// 			}
// 		}

// 	}

// 	// err := store.SetLastUpdatedNow()
// 	// if err != nil {
// 	// 	log.Printf("failed setting last updated %s", err)
// 	// }

// 	// return time.Now()

// }

// func (mAPI *MetOfficeAPI) updater(log *log.Logger, store Store) {
// 	wg := sync.WaitGroup{}

// 	// really this should be gained through di and gotten from the sql db somewhere
// 	for _, crag := range crags {
// 		wg.Add(1)
// 		go func(crag Crag) {

// 			log.Printf("go %s route started", crag.Name)
// 			// this isnt how we want to be adding the url into the code but we dont have data in our
// 			// crag db at the moment
// 			// this will require adding the htmx form to add data into that database
// 			// and then also creating a notification system of kind where when a crag is addeed to the
// 			// tracker that we add it to the forecasts that we want to get too
// 			f, err := mAPI.GetForecast(mAPI.CreateURL([]float64{crag.Latitude, crag.Longitude}))
// 			if err != nil {
// 				return
// 			}
// 			p, err := mAPI.GetPayload(log, f)
// 			if err != nil {
// 				log.Printf("failed creating payload %s", err)
// 			}
// 			if err := store.Add(context.Background(), crag.Name, p); err != nil {
// 				log.Printf("failed storing forecast totals, %s", err)
// 			}
// 			log.Printf("go %s route done", crag.Name)
// 			wg.Done()
// 		}(crag)
// 	}
// 	wg.Wait()
// }

// TODO: Add track previous potential windows, to somehow examine likelihood of when certain place could be dry based on previous weather

// temp solution
// var crags = []Crag{
// 	{"cromlech", 53.08977582752912, -4.0494354521953895},
// 	{"beddgelert", 53.01401346937128, -4.1086367318613055},
// 	{"gwynant", 53.04567339439013, -4.021447439922229},
// 	{"blaenau", 52.99729599359651, -3.9578734953238475},
// 	{"crafnant", 52.99729599359651, -3.9578734953238475},
// 	{"cwellyn", 53.07568570139747, -4.148701296939546},
// 	{"orme", 53.33236585445307, -3.8311890286450865},
// 	{"Penmaenbach", 53.285, -3.8684},
// 	{"ysgo", 52.80614677538971, -4.656639551730091},
// 	{"tremadoch", 52.94008535336955, -4.140997768369204},
// 	{"rhiwGoch", 53.09199013529737, -3.803795346023221},
// 	{"portland", 50.545900401402854, -2.438814867485551},
// 	{"cuckooRock", 50.545900401402854, -2.438814867485551},
// 	{"MountSionEast", 50.545900401402854, -2.438814867485551},
// 	{"Froggatt", 53.2942103060766, -1.6201285054945418},
// }
