package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
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
	// I probbaly just want to store the urls with the crag information that way I only need to create them once
	var sb strings.Builder
	sb.WriteString(mAPI.BaseURL)
	sb.WriteString("latitude=")
	sb.WriteString(strconv.FormatFloat(coords[0], 'f', -1, 64))
	sb.WriteString("longitude=")
	sb.WriteString(strconv.FormatFloat(coords[1], 'f', -1, 64))
	return sb.String()
}

func (mAPI *MetOfficeAPI) GetHeaders() http.Header {
	return http.Header{

		"apikey": {mAPI.APIKey},
		"accept": {"application/json"},
	}
}

func (mAPI *MetOfficeAPI) CalculateTotals(data []TimeSeriesData) map[string]*ForecastTotals {

	var totals = make(map[string]*ForecastTotals)

	for _, val := range data {
		day := val.Time[8:10]
		entry, ok := totals[day]
		if !ok {
			entry = &ForecastTotals{
				HighestTemp: val.MaxScreenAirTemp,
				LowestTemp:  val.MinScreenAirTemp,
				Datapoints:  0,
			}
			totals[day] = entry
		}

		entry.AvgTemp += val.ScreenTemperature
		entry.AvgWindSpeed += val.WindSpeed10m
		entry.TotalPrecip += val.TotalPrecipAmount
		entry.Datapoints++
	}

	for _, entry := range totals {
		if entry.Datapoints > 0 {
			entry.AvgTemp /= float64(entry.Datapoints)
			entry.AvgWindSpeed /= float64(entry.Datapoints)

		}
	}
	return totals
}

func (mAPI *MetOfficeAPI) FindWindows(data []TimeSeriesData) [][]time.Time {

	//probably should be returning and handling our errors in this

	var startOfWindow string
	var endOfWindow string

	windows := [][]time.Time{}
	for i := 0; i < len(data); i++ {

		currForecast := data[i]

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

		if i == (len(data) - 1) {
			if currForecast.TotalPrecipAmount != 0.00 && startOfWindow != "" {
				endOfWindow = data[i-1].Time
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
				endOfWindow = data[i-1].Time
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
