package main

import (
	"log"
	"time"
)

type Scheduler struct {
	timer *time.Timer
}

const intervalPeriod time.Duration = 1 * time.Hour

//this hasnt fixed it quiete because im still passing the payload which we DONT WANT. The whole point of this function is to
//gather whether we need to get the payload or not. So, the called needs to already have the lastUpdate time, and then
//it needs to call this function which will decide whether to trigger the store data function or not. Not actually trigger the store function data itself.

func (s *Scheduler) Update(lastUpdated time.Time) {

	var nextTick time.Time

	log.Println("Starting scheduler")

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
		return
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

func (s *Scheduler) startSchedule(log *log.Logger, metAPI *MetOfficeAPI, store *MetStore, lastUpdate time.Time) {

	log.Println("checking reschedule")

	s.Update(lastUpdate)
	for {
		<-s.timer.C
		lastUpdate := ExecuteRefreshProcess(log, metAPI, store)
		s.Update(lastUpdate) //placeholder for the last update time from the api

	}
}
