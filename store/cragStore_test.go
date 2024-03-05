package store

import (
	"testing"

	log "github.com/sirupsen/logrus"
)

func CreateStore(t *testing.T) CragStore {
	store, err := NewSqlStore()
	if err != nil {
		log.Fatalf("error creating store: %s", err)
	}
	return store.Stores.CragStore
}
