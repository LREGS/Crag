package store

type SqlClimbStore struct {
	Store *SqlStore
}

func NewClimbStore(sqlStore *SqlStore) *SqlClimbStore {
	store := &SqlClimbStore{sqlStore}
	return store
}
