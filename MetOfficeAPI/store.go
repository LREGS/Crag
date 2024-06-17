package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type MetStore struct {
	Log       *log.Logger
	Rdb       *redis.Client //Redis Database
	scheduler *Scheduler
}

func NewMetStore(rdb *redis.Client, log *log.Logger) *MetStore {
	return &MetStore{Log: log, Rdb: rdb}
}

func (m *MetStore) ForecastTotals(ctx context.Context, payload ForecastPayload) error {

	// clear cache before every store as we only care about the last hours results
	if err := m.Flush(); err != nil {
		log.Printf("couldn't update cache because of error whilst flushing %s", err)
		return err
	}

	//	is it better to split the storage like this from a single payload or should it be independant
	data, err := json.Marshal(payload.ForecastTotals)
	if err != nil {
		log.Printf("failed marshalling %s", err)
		return err
	}

	windows, err := json.Marshal(payload.Windows)
	if err != nil {
		return err
	}

	// this is adding the time into redis like "2024-06-13 21:00:00.31195387 +0100 BST m=+2814.586676661" and
	// its causing the scheduler to fail because it cannot parse a time in this format
	if err := m.Rdb.Set(ctx, "LastUpdated", time.Now().String(), 0).Err(); err != nil {
		log.Printf("error storing last updated %s", err)
		return err
	}

	if err := m.Rdb.Set(ctx, "totals", data, 0).Err(); err != nil {
		log.Printf("error storing totals %s", err)
		return err
	}

	if err := m.Rdb.Set(ctx, "windows", windows, 0).Err(); err != nil {
		log.Printf("failed storing windows %s", err)
		return err
	}

	return nil

}

func (m *MetStore) Flush() error {
	err := m.Rdb.FlushDB(context.Background()).Err()
	if err != nil {
		return err
	}

	return nil
}

var ErrorRedis = errors.New("redis empty, cannot get last updated")

func (m *MetStore) GetForecastTotals() (map[string]*ForecastTotals, error) {

	var totals map[string]*ForecastTotals

	res, err := m.Rdb.Get(context.Background(), "totals").Bytes()
	if err != nil {
		return totals, err
	}

	if err := json.Unmarshal(res, &totals); err != nil {
		return totals, err
	}

	return totals, nil
}

func (m *MetStore) GetLastUpdate(log *log.Logger) (time.Time, error) {
	res, err := m.Rdb.Get(context.Background(), "LastUpdated").Result()
	if err == redis.Nil {
		log.Printf("last update failed: no entry exists %s", err)
		return time.Time{}, ErrorRedis
	}
	if err != nil {
		log.Printf("failed getting last update from redis %s", err)
		return time.Time{}, err
	}

	parsedTime, err := time.Parse("2006-01-02T15:04Z07:00", res)
	if err != nil {
		return time.Time{}, err
	}

	return parsedTime, nil
}

// The scheduler is not part of the Store, the same way as the api access isn't part of the store.
// The scheduler is instead the orchestrator, which will bring together the methods of the api access
// store methods within the met office package

// func (m *MetStore) InitScheduler(lastUpdate time.Time) {

// 	m.Scheduler.Update(lastUpdate)
// 	for {
// 		<-m.Scheduler.timer.C
// 		if err := m.StoreForecastTotals(context.Background()); err != nil {
// 			log.Printf("failed to store data in scheduler")
// 		}
// 		time, err := time.Parse("2006-01-02T15:04Z07:00", payload.LastModelRunTime)
// 		if err != nil {
// 			log.Printf("failed parsing time during cache update")
// 		}
// 		m.Scheduler.Update(time)

// 	}
// }
