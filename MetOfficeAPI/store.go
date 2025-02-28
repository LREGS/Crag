package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Store interface {
	Add(ctx context.Context, name string, payload ForecastPayload) error
	Get() (map[string]*ForecastTotals, error)
	GetLastUpdate() time.Time
}

type MetStore struct {
	Log *log.Logger
	Rdb *redis.Client //Redis Database
}

func NewMetStore(rdb *redis.Client, log *log.Logger) *MetStore {
	return &MetStore{Log: log, Rdb: rdb}
}

func (m *MetStore) Add(ctx context.Context, key string, payload ForecastPayload) error {

	// if err := m.Flush(); err != nil {
	// 	log.Printf("couldn't update cache because of error whilst flushing %s", err)
	// 	return err
	// }

	p, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	if err := m.Rdb.Set(ctx, key, p, 0).Err(); err != nil {
		log.Printf("failed storing payload %s", err)
		return err
	}

	if err = m.SetLastUpdatedNow(); err != nil {
		log.Printf("failed setting LastUpdated %s", err)
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

func (m *MetStore) Get(key string) (ForecastPayload, error) {

	var totals ForecastPayload

	res, err := m.Rdb.Get(context.Background(), key).Bytes()
	if err != nil {
		return totals, err
	}

	if err := json.Unmarshal(res, &totals); err != nil {
		return totals, err
	}

	return totals, nil
}

func (m *MetStore) MuiltiForecasts(ctx context.Context, keys []string) ([]ForecastPayload, error) {
	cmd := m.Rdb.MGet(ctx, keys...)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	results := cmd.Val()

	forecasts := make([]ForecastPayload, len(results))

	for i, result := range results {

		var forecast ForecastPayload
		if err := json.Unmarshal(result.([]byte), &forecast); err != nil {
			return nil, fmt.Errorf("failed to unmarshal %d %s", i, err)
		}

		forecasts[i] = forecast

	}

	return forecasts, nil
}

func (m *MetStore) GetLastUpdate() time.Time {
	// this always fails if the db is restarted and there is no lastUpdated because it will always return nil
	res, err := m.Rdb.Get(context.Background(), "LastUpdated").Result()
	if err == redis.Nil {
		log.Printf("last update failed: no entry exists %s", err)
		// return time.Time{}, ErrorRedis
	}

	if err != nil {
		log.Printf("failed getting last update from redis %s", err.Error())
		return time.Now().Add(-2 * time.Hour)
	}

	parsedTime, err := time.Parse("2006-01-02T15:04Z07:00", res)
	if err != nil {
		log.Printf("failed parsing time %s", err)
		// return time.Time{}, err
		return time.Time{}
	}

	return parsedTime
}

func (m *MetStore) SetLastUpdatedNow() error {
	if err := m.Rdb.Set(context.Background(), "LastUpdated", time.Now().Format("2006-01-02T15:04Z07:00"), 0).Err(); err != nil {
		log.Printf("error storing last updated %s", err)
		return err
	}
	return nil
}
