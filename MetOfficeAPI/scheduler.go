package main

import (
	"log"
	"time"
)

type Scheduler struct {
	timer *time.Timer
}

const intervalPeriod time.Duration = 1 * time.Hour

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


func (s *Scheduler) startSchedule(log *log.Logger, metAPI *MetOfficeAPI, store *MetStore, lastUpdate time.Time) {

	log.Println("checking reschedule")

	s.Update(lastUpdate)
	for {
		<-s.timer.C
		lastUpdate := ExecuteRefreshProcess(log, metAPI, store)
		s.Update(lastUpdate) 

	}
}
