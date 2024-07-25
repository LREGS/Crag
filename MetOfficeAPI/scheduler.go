package main

import (
	"context"
	"log"
	"sync"
	"time"
)

// im wondering if this whole thing should just be a go routine in a neverending loop
// that sends data to a channel.
// When the loop starts it should check last update, create a ticker based on this
// and then call the update function . - or maybe this is handled by watching the channel

type Scheduler struct {
	timer *time.Timer
}

const intervalPeriod time.Duration = 1 * time.Hour

func (s *Scheduler) Update(log *log.Logger, lastUpdated time.Time) {

	var nextTick time.Time

	log.Println("Starting scheduler")

	if time.Since(lastUpdated).Hours() > 1.00 {
		log.Println("starting update immediately")
		nextTick = time.Now()
	} else {
		nextTick = lastUpdated.Truncate(time.Hour).Add(time.Hour)
	}

	if !nextTick.After(time.Now()) {
		nextTick = nextTick.Add(intervalPeriod)
	}

	diff := time.Until(nextTick)
	log.Println(diff)
	if s.timer == nil {
		s.timer = time.NewTimer(0)
		return
	} else {
		s.timer.Reset(diff)
	}

}

func (s *Scheduler) startSchedule(log *log.Logger, metAPI *MetOfficeAPI, store *MetStore, lastUpdate time.Time) {

	log.Println("checking reschedule")

	s.Update(log, lastUpdate)
	for {
		<-s.timer.C
		lastUpdate := s.ExecuteRefreshProcess(log, metAPI, store)
		s.Update(log, lastUpdate)

	}
}

func (s *Scheduler) ExecuteRefreshProcess(log *log.Logger, api *MetOfficeAPI, store *MetStore) time.Time {

	// if err := store.Flush(); err != nil {
	// 	log.Printf("couldn't update cache because of error whilst flushing %s", err)
	// }

	log.Println("starting refresh")

	crags := []Crag{
		{"cromlech", 53.08977582752912, -4.0494354521953895},
		{"beddgelert", 53.01401346937128, -4.1086367318613055},
		{"gwynant", 53.04567339439013, -4.021447439922229},
		{"blaenau", 52.99729599359651, -3.9578734953238475},
		{"crafnant", 52.99729599359651, -3.9578734953238475},
		{"cwellyn", 53.07568570139747, -4.148701296939546},
		{"orme", 53.33236585445307, -3.8311890286450865},
		{"Penmaenbach", 53.285, -3.8684},
		{"ysgo", 52.80614677538971, -4.656639551730091},
		{"tremadoch", 52.94008535336955, -4.140997768369204},
		{"rhiwGoch", 53.09199013529737, -3.803795346023221},
		{"portland", 50.545900401402854, -2.438814867485551},
		{"cuckooRock", 50.545900401402854, -2.438814867485551},
		{"MountSionEast", 50.545900401402854, -2.438814867485551},
		{"Froggatt", 53.2942103060766, -1.6201285054945418},
	}

	wg := sync.WaitGroup{}

	log.Print(crags)

	for _, crag := range crags {
		wg.Add(1)
		go func(crag Crag) {

			log.Printf("go %s route started", crag.Name)
			f, err := api.GetForecast(api.CreateURL([]float64{crag.Latitude, crag.Longitude}))
			if err != nil {
				return
			}
			p, err := api.GetPayload(log, f)
			if err != nil {
				log.Printf("failed creating payload %s", err)
			}

			if err := store.Totals(context.Background(), crag.Name, p); err != nil {
				log.Printf("failed storing forecast totals, %s", err)
			}
			log.Printf("go %s route done", crag.Name)
			wg.Done()
		}(crag)
	}
	wg.Wait()

	// this is adding the time into redis like "2024-06-13 21:00:00.31195387 +0100 BST m=+2814.586676661" and
	// its causing the scheduler to fail because it cannot parse a time in this format

	err := store.SetLastUpdatedNow()
	if err != nil {
		log.Printf("failed setting last updated %s", err)
	}

	return time.Now()

}

var crags = []Crag{
	{"cromlech", 53.08977582752912, -4.0494354521953895},
	{"beddgelert", 53.01401346937128, -4.1086367318613055},
	{"gwynant", 53.04567339439013, -4.021447439922229},
	{"blaenau", 52.99729599359651, -3.9578734953238475},
	{"crafnant", 52.99729599359651, -3.9578734953238475},
	{"cwellyn", 53.07568570139747, -4.148701296939546},
	{"orme", 53.33236585445307, -3.8311890286450865},
	{"Penmaenbach", 53.285, -3.8684},
	{"ysgo", 52.80614677538971, -4.656639551730091},
	{"tremadoch", 52.94008535336955, -4.140997768369204},
	{"rhiwGoch", 53.09199013529737, -3.803795346023221},
	{"portland", 50.545900401402854, -2.438814867485551},
	{"cuckooRock", 50.545900401402854, -2.438814867485551},
	{"MountSionEast", 50.545900401402854, -2.438814867485551},
	{"Froggatt", 53.2942103060766, -1.6201285054945418},
}

func ScheduleMetOffice(log *log.Logger, api *MetOfficeAPI, store *MetStore) time.Time {

	// if err := store.Flush(); err != nil {
	// 	log.Printf("couldn't update cache because of error whilst flushing %s", err)
	// }

	for {

		t, err := store.GetLastUpdate(log)
		if err != nil {
			// do we need to be sending errors to a channel or something

			// yes when this doesnt parse correctly we're calling the api infinitely
			log.Println(err)
		}

		log.Printf("last updated %s", t.String())

		if time.Since(t) > time.Hour {
			log.Printf("updating now")
			updater(log, api, store)
		} else {
			log.Print("will update later")
			c := time.Tick(time.Duration(60-(time.Now().Minute())) * time.Minute)
			for next := range c {
				log.Print(next)
				updater(log, api, store)
			}
		}

	}

	// this is adding the time into redis like "2024-06-13 21:00:00.31195387 +0100 BST m=+2814.586676661" and
	// its causing the scheduler to fail because it cannot parse a time in this format

	// err := store.SetLastUpdatedNow()
	// if err != nil {
	// 	log.Printf("failed setting last updated %s", err)
	// }

	// return time.Now()

}

func updater(log *log.Logger, api *MetOfficeAPI, store *MetStore) {
	wg := sync.WaitGroup{}

	// really this should be gained through di and gotten from the sql db somewhere
	for _, crag := range crags {
		wg.Add(1)
		go func(crag Crag) {

			log.Printf("go %s route started", crag.Name)
			f, err := api.GetForecast(api.CreateURL([]float64{crag.Latitude, crag.Longitude}))
			if err != nil {
				return
			}
			p, err := api.GetPayload(log, f)
			if err != nil {
				log.Printf("failed creating payload %s", err)
			}

			// this shouldn't be in the same function

			if err := store.Totals(context.Background(), crag.Name, p); err != nil {
				log.Printf("failed storing forecast totals, %s", err)
			}
			log.Printf("go %s route done", crag.Name)
			wg.Done()
		}(crag)
	}
	wg.Wait()
}
