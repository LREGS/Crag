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

func (m *MetStore) Totals(ctx context.Context, name string, payload ForecastPayload) error {

	// if err := m.Flush(); err != nil {
	// 	log.Printf("couldn't update cache because of error whilst flushing %s", err)
	// 	return err
	// }

	p, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	if err := m.Rdb.Set(ctx, name, p, 0).Err(); err != nil {
		log.Printf("failed storing payload %s", err)
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

func (m *MetStore) GetTotals() (map[string]*ForecastTotals, error) {

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

func (m *MetStore) SetLastUpdatedNow() error {
	if err := m.Rdb.Set(context.Background(), "LastUpdated", time.Now().String(), 0).Err(); err != nil {
		log.Printf("error storing last updated %s", err)
		return err
	}
	return nil
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
