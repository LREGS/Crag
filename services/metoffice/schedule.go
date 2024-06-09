package met

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type scheduler struct {
	timer *time.Timer
}

const intervalPeriod time.Duration = 1 * time.Hour

//this hasnt fixed it quiete because im still passing the payload which we DONT WANT. The whole point of this function is to
//gather whether we need to get the payload or not. So, the called needs to already have the lastUpdate time, and then
//it needs to call this function which will decide whether to trigger the store data function or not. Not actually trigger the store function data itself.

func UpdateCacheScheduler(log *log.Logger, ctx context.Context, rdb *redis.Client, payload ForecastPayload, lastUpdate time.Time) {
	scheduler := &scheduler{}
	scheduler.Update(lastUpdate)
	for {
		<-scheduler.timer.C
		if err := StoreData(log, ctx, rdb, payload); err != nil {
			log.Printf("failed to store data in scheduler")
		}
		time, err := time.Parse("2006-01-02T15:04Z07:00", payload.LastModelRunTime)
		if err != nil {
			log.Printf("failed parsing time during cache update")
		}
		scheduler.Update(time)

	}
}

func (s *scheduler) Update(lastUpdated time.Time) {

	var nextTick time.Time

	if time.Since(lastUpdated).Hours() > 1.00 {
		nextTick = time.Now()
	} else {
		nextTick = lastUpdated.Truncate(time.Hour).Add(time.Hour)
	}

	if !nextTick.After(time.Now()) {
		nextTick = nextTick.Add(intervalPeriod)
	}

	diff := time.Until(nextTick)
	if s.timer == nil {
		s.timer = time.NewTimer(diff)
	} else {
		s.timer.Reset(diff)
	}

}
