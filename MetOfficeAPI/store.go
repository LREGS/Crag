package main

import (
	"context"
	"encoding/json"
	"errors"
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

func (m *MetStore) GetAll() (map[string]ForecastPayload, error) {

	redisKeys := make([]string, len(crags))
	for _, crag := range crags {
		redisKeys = append(redisKeys, crag.Name)
	}

	vals, err := m.Rdb.MGet(context.Background(), redisKeys...).Result()
	if err != nil {
		m.Log.Print("failed getting all data from redis ", err)
		return nil, err
	}

	res := make(map[string]ForecastPayload)
	for i, val := range vals {
		if val != nil {
			var payload ForecastPayload
			if err := json.Unmarshal([]byte(val.(string)), &payload); err != nil {
				m.Log.Print("failed decoding redis response from get all ", err)
				continue
			}
			res[redisKeys[i]] = payload
		}
	}
	return res, nil
}

func (m *MetStore) GetKeys() []string {

	var cursor uint64
	var allKeys []string
	for {
		var err error
		var keys []string
		keys, cursor, err = m.Rdb.Scan(context.Background(), cursor, "*", 10).Result()
		if err != nil {
			m.Log.Print("failed getting keys", err)
		}

		allKeys = append(allKeys, keys...)

		if cursor == 0 {
			break
		}

	}

	return allKeys

}

func (m *MetStore) GetLastUpdate() time.Time {
	res, err := m.Rdb.Get(context.Background(), "LastUpdated").Result()
	// if err == redis.Nil {
	// 	log.Printf("last update failed: no entry exists %s", err)
	// 	// return time.Time{}, ErrorRedis
	// }
	if err != nil {
		log.Printf("failed getting last update from redis %s", err)
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
