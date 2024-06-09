package met

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type TimeSeriesData struct {
	Time                      string  `json:"time"`
	ScreenTemperature         float64 `json:"screenTemperature"`
	MaxScreenAirTemp          float64 `json:"maxScreenAirTemp"`
	MinScreenAirTemp          float64 `json:"minScreenAirTemp"`
	ScreenDewPointTemperature float64 `json:"screenDewPointTemperature"`
	FeelsLikeTemperature      float64 `json:"feelsLikeTemperature"`
	WindSpeed10m              float64 `json:"windSpeed10m"`
	WindDirectionFrom10m      int     `json:"windDirectionFrom10m"`
	WindGustSpeed10m          float64 `json:"windGustSpeed10m"`
	Max10mWindGust            float64 `json:"max10mWindGust"`
	Visibility                int     `json:"visibility"`
	ScreenRelativeHumidity    float64 `json:"screenRelativeHumidity"`
	Mslp                      int     `json:"mslp"`
	UvIndex                   int     `json:"uvIndex"`
	SignificantWeatherCode    int     `json:"significantWeatherCode"`
	PrecipitationRate         float64 `json:"precipitationRate"`
	TotalPrecipAmount         float64 `json:"totalPrecipAmount"`
	TotalSnowAmount           float64 `json:"totalSnowAmount"`
	ProbOfPrecipitation       int     `json:"probOfPrecipitation"`
}

type Feature struct {
	Type        string     `json:"type"`
	Geometry    Geometry   `json:"geometry"`
	Properties  Properties `json:"properties"`
	Coordinates []float64  `json:"coordinates"`
}

type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type Properties struct {
	RequestPointDistance float64          `json:"requestPointDistance"`
	ModelRunDate         string           `json:"modelRunDate"`
	TimeSeries           []TimeSeriesData `json:"timeSeries"`
}

type Forecast struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

// Selects the data we want from the response

type DBForecast struct {
	Id                  int     `json:"id"`
	Time                string  `json:"time"`
	ScreenTemperature   float64 `json:"screenTemperature"`
	FeelsLikeTemp       float64 `json:"feelsLikeTemp"`
	WindSpeed           float64 `json:"windSpeed"`
	WindDirection       int     `json:"windDirection"`
	TotalPrecipAmount   float64 `json:"totalPrecipAmount"`
	ProbOfPrecipitation int     `json:"probOfPrecipitation"`
	Latitude            float64 `json:"latitude"`
	Longitude           float64 `json:"longitude"`
}

type DBForecastPayload struct {
	Time                string  `json:"time"`
	ScreenTemperature   float64 `json:"screenTemperature"`
	FeelsLikeTemp       float64 `json:"feelsLikeTemp"`
	WindSpeed           float64 `json:"windSpeed"`
	WindDirection       int     `json:"windDirection"`
	TotalPrecipAmount   float64 `json:"totalPrecipAmount"`
	ProbOfPrecipitation int     `json:"probOfPrecipitation"`
	Latitude            float64 `json:"latitude"`
	Longitude           float64 `json:"longitude"`
	CragId              int     `json:"cragId"`
}

//do I need to have a struct that has the methods or just the functions I dont know

// returns the forecast for a crag based on its stored coords
func GetForecast(client http.Client, coords []float64) (Forecast, error) {
	var forecast Forecast

	//this should be recieving a client so im not making a new one with every request plls

	url := fmt.Sprintf("https://data.hub.api.metoffice.gov.uk/sitespecific/v0/point/hourly?latitude=%f&longitude=%f", coords[0], coords[1])

	//need this back online

	// if err := godotenv.Load(); err != nil {
	// 	return forecast, err
	// }

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return forecast, err
	}

	req.Header = http.Header{

		"apikey": {"eyJ4NXQiOiJOak16WWpreVlUZGlZVGM0TUdSalpEaGtaV1psWWpjME5UTXhORFV4TlRZM1ptRTRZV1JrWWc9PSIsImtpZCI6ImdhdGV3YXlfY2VydGlmaWNhdGVfYWxpYXMiLCJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJzdWIiOiJ3aWxsaWFtLWN1bGx1bUBob3RtYWlsLmNvLnVrQGNhcmJvbi5zdXBlciIsImFwcGxpY2F0aW9uIjp7Im93bmVyIjoid2lsbGlhbS1jdWxsdW1AaG90bWFpbC5jby51ayIsInRpZXJRdW90YVR5cGUiOm51bGwsInRpZXIiOiJVbmxpbWl0ZWQiLCJuYW1lIjoic2l0ZV9zcGVjaWZpYy1hYjE3ZTkyMy1kODI2LTQ5ZDQtYWZhMC01ODY3ZTQxODMwNzciLCJpZCI6NDIwNCwidXVpZCI6IjA4NzBjYjI5LWIyNDYtNDE3OS05OWQzLTE1ZDg5Njc5MGE0MSJ9LCJpc3MiOiJodHRwczpcL1wvYXBpLW1hbmFnZXIuYXBpLW1hbmFnZW1lbnQubWV0b2ZmaWNlLmNsb3VkOjQ0M1wvb2F1dGgyXC90b2tlbiIsInRpZXJJbmZvIjp7IndkaF9zaXRlX3NwZWNpZmljX2ZyZWUiOnsidGllclF1b3RhVHlwZSI6InJlcXVlc3RDb3VudCIsImdyYXBoUUxNYXhDb21wbGV4aXR5IjowLCJncmFwaFFMTWF4RGVwdGgiOjAsInN0b3BPblF1b3RhUmVhY2giOnRydWUsInNwaWtlQXJyZXN0TGltaXQiOjAsInNwaWtlQXJyZXN0VW5pdCI6InNlYyJ9fSwia2V5dHlwZSI6IlBST0RVQ1RJT04iLCJzdWJzY3JpYmVkQVBJcyI6W3sic3Vic2NyaWJlclRlbmFudERvbWFpbiI6ImNhcmJvbi5zdXBlciIsIm5hbWUiOiJTaXRlU3BlY2lmaWNGb3JlY2FzdCIsImNvbnRleHQiOiJcL3NpdGVzcGVjaWZpY1wvdjAiLCJwdWJsaXNoZXIiOiJKYWd1YXJfQ0kiLCJ2ZXJzaW9uIjoidjAiLCJzdWJzY3JpcHRpb25UaWVyIjoid2RoX3NpdGVfc3BlY2lmaWNfZnJlZSJ9XSwidG9rZW5fdHlwZSI6ImFwaUtleSIsImlhdCI6MTcxNTcxMjE4MywianRpIjoiYjFkMDZjZDctZDViNy00OWFkLThhM2YtNzhjMDRkZjZjM2ZhIn0=.PzEmpYP8PqjABpjJN8z4LbfILUgALEybfjJIJ2IrV9gwWV9GoLTGjqFrajJ0QkSJID80HuOlDLp0psYZLgSVe-l1DpAe6FtwYoBP6TaTN8PdiLle5m7JFlIR-sYd_iXDHUpAnjWrNh1u_Ofz8bPcQZ8F5szg9DbZQL_umsw-ST5L01tl3PEmqBkZiZ20rCMTxK9OMpoukfX2iPX4US48sIny6XVExLAXvZXt4uFqnChInFJXkIALksndiUm8OL_sDwbraPZKd1MUgII-SBtbJcp-nWCI3J9oNvKrK53HlqOWPTTcrsLnZkcJLmhQbvXXEgqyRnJM5Usa297EveFWYA=="},
		"accept": {"application/json"},
	}

	res, err := client.Do(req)
	if err != nil {
		return forecast, err

	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return forecast, fmt.Errorf("code %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return forecast, err

	}

	err = json.Unmarshal(body, &forecast)
	if err != nil {
		return forecast, err
	}

	// defer res.Body.Close()
	// err = json.NewDecoder(res.Body).Decode(&forecast)
	// if err != nil {
	// 	return forecast, err
	// }

	return forecast, nil

}

// type metOfficeHeaders struct {
// 	ApiKey string `json:"apikey"`
// 	Accept string `json:"accept"`
// }

// func getHeaders() (metOfficeHeaders, error) {
// 	env, err := util.GetEnv([]string{"apikey"})
// 	if err != nil {
// 		return metOfficeHeaders{}, nil
// 	}

// 	return metOfficeHeaders{apikey: env[0], Accept: "application/json"}, nil
// }

// func GetPayload(log *log.Logger, coords []float64) ([]DBForecast, error) {

// 	//if get forecast fails we get an index out of range error because of the timeSeries
// 	//im not sure why the error is obviously being returned as nil but tis annoying

// 	forecast, err := GetForecast(coords)
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	if len(forecast.Features) == 0 {
// 		return nil, errors.New("empty forecast")
// 	}

// 	timeSeries := forecast.Features[0].Properties.TimeSeries

// 	//not sure I need [][]interface{} as this was for sql copy
// 	payload := make([]DBForecast, len(timeSeries))

// 	for i := 0; i < (len(timeSeries) - 1); i++ {

// 		payload[i] = DBForecast{
// 			Id:                  i + 1, //Id
// 			Time:                timeSeries[i].Time,
// 			ScreenTemperature:   timeSeries[i].ScreenTemperature,
// 			FeelsLikeTemp:       timeSeries[i].FeelsLikeTemperature,
// 			WindSpeed:           timeSeries[i].WindSpeed10m,
// 			WindDirection:       timeSeries[i].WindDirectionFrom10m,
// 			TotalPrecipAmount:   timeSeries[i].TotalPrecipAmount,
// 			ProbOfPrecipitation: timeSeries[i].ProbOfPrecipitation,
// 			Latitude:            forecast.Features[0].Geometry.Coordinates[0],
// 			Longitude:           forecast.Features[0].Geometry.Coordinates[1],
// 		}

// 	}

// 	return payload, nil

// }

// this doesnt show prob of precipitation because that needs to be hourly, not totals

//These stats will provide header // outline stats for each crag but the full hourly will also be available on inspection I guess for all crags
//but a condensed leaderboard style report will also be available

// window is a weather window and the average weather conditions within the weather window are shown.
// I'm not sure how useful the whole scale averages are for multi day windows but I guess the full hourly forecast will still be presented somewhere alongside the windows

// This will represent a single window in time detailed by the time string.
// each window is a minimum of 2 hours long
type Window struct {
	Time          []string
	AvgTemp       float64
	AvgWindSpeed  float64
	WindDirection int
}

// TODO: Change to dailyTotals
type ForecastTotals struct {
	HighestTemp   float64
	LowestTemp    float64
	AvgTemp       float64
	AvgWindSpeed  float64
	WindDirection int
	TotalPrecip   float64
	// Windows       [][]int
}

func (f ForecastTotals) MarshalBinary() (data []byte, err error) {
	return json.Marshal(f)
}

func UnmarshalBinary(data []byte) error {
	return nil
}

// type MeanCalculator struct {
// 	count int
// 	mean float64
// }

//  func (mc *MeanCalculator) Add(value float64){
// 	mc.count++
// 	delta := value - mc.mean
// 	mc.mean += delta / float64(mc.count)
//  }

type ForecastPayload struct {
	LastModelRunTime string
	ForecastTotals   map[string]*ForecastTotals
}

func GetPayload(log *log.Logger, forecast Forecast) (ForecastPayload, error) {

	//1 hourly spot has three days worth of data. This function provides the totals for
	//all three days

	data := forecast.Features[0].Properties.TimeSeries

	if len(data) == 0 {
		return ForecastPayload{}, errors.New("Forecast provided is empty")
	}

	// modelStartDate, err := strconv.Atoi(forecast.Features[0].Properties.ModelRunDate[8:10])
	// if err != nil {
	// 	return nil, err
	// }

	totals := map[string]*ForecastTotals{}

	//TODO: whats the better way than just reallocating a lot of the values here already
	//put them into days and then total them? That will require two passes of the data
	//but wont be continually reallocating the avgtemp etc

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

	return ForecastPayload{LastModelRunTime: forecast.Features[0].Properties.ModelRunDate, ForecastTotals: totals}, nil

}

func StoreData(log *log.Logger, ctx context.Context, rdb *redis.Client, payload ForecastPayload) error {

	data, err := json.Marshal(payload.ForecastTotals)
	if err != nil {
		log.Printf("failed marshalling %s", err)
		return err
	}

	//marks the time the db was last updated as that we can have a consistent key across both.
	//init redis and update redis will both check the timestamp first to make sure the data actually needs to be updated
	//some way needs to be resolved where we can check if its the same date but less than an hour previous

	exists, err := rdb.Exists(context.Background(), "LastUpdated").Result()
	if err != nil {
		log.Printf("failed checking key for last update %s", err)
		return err
	}

	if exists != 0 {
		b, err := CheckLastUpdated(rdb, payload.LastModelRunTime)
		if err != nil {
			log.Printf("error checking update status %s", err)
			return err
		}
		if b {
			return errors.New("no update required")
		}

	}

	if err := rdb.Set(ctx, "LastUpdated", payload.LastModelRunTime, 0).Err(); err != nil {
		log.Printf("error storing last updated %s", err)
		return err
	}

	if err := rdb.Set(ctx, "totals", data, 0).Err(); err != nil {
		log.Printf("error storing totals %s", err)
		return err
	}

	rd, err := rdb.Get(ctx, "totals").Result()
	if err != nil {
		log.Printf("error getting totals, %s", err)
		return err
	}

	log.Printf("data stored %s", rd)

	return nil
}

func CheckLastUpdated(rdb *redis.Client, LastRunTime string) (bool, error) {

	res, err := rdb.Get(context.Background(), "LastUpdated").Result()
	if err != nil {
		return false, err
	}

	if res != LastRunTime {
		return true, nil
	}

	return false, nil

}

func GetLastUpdateTime(rdb *redis.Client) (time.Time, error) {
	res, err := rdb.Get(context.Background(), "LastUpdated").Result()
	if err != nil {
		return time.Time{}, err
	}

	parsedTime, err := time.Parse("2006-01-02T15:04Z07:00", res)
	if err != nil {
		return time.Time{}, err
	}

	return parsedTime, nil

}
