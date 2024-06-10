package met

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type scheduler struct {
	store *MetOfficeStore
	APIAccess *MetOfficeAPI
	timer *time.Timer
}

const intervalPeriod time.Duration = 1 * time.Hour

//this hasnt fixed it quiete because im still passing the payload which we DONT WANT. The whole point of this function is to
//gather whether we need to get the payload or not. So, the called needs to already have the lastUpdate time, and then
//it needs to call this function which will decide whether to trigger the store data function or not. Not actually trigger the store function data itself.

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

// im still unsure where I want to now initialise this, I know I need the scheduler to call 
// the store function, but Im not sure if im supposed to decouple this a little bit more 
// because atm its completely tied into the Met package, but is that idiomatic because of the whole
// package tells a story but I dont  think that means we need to call everything all in this scheduler 
// but we do want the scheduler to call the thing that calls everything maybe 
// or do we want the scheduler to passed the thing, by the thing that calls everything, including the scheduler.
// maybe thats the way.

func (s *scheduler) startSchedule(log *log.Logger, rdb *redis.Client) {

	MetOfficeAPI

	payload := 

	store := NewMetStore(rdb, log)

	lastUpdated, err := store.GetLastUpdate(log)
	if err == RedisEmpty {

	}

	s.Update(lastUpdate)
	for {
		<-s.timer.C
		if err := s.StoreForecastTotals(context.Background()); err != nil {
			log.Printf("failed to store data in scheduler")
		}
		time, err := time.Parse("2006-01-02T15:04Z07:00", payload.LastModelRunTime)
		if err != nil {
			log.Printf("failed parsing time during cache update")
		}
		scheduler.Update(time)

	}
}

func worker() {

}
