package main

func NewInMemoryCragStore() *InMemoryCragStore {
	return &InMemoryCragStore{make(map[string][]string)}
}

type InMemoryCragStore struct {
	store map[string][]string
}

func (i *InMemoryCragStore) GetForecast(crag string) string {
	if len(i.store[crag]) == 0 {
		return ""
	}
	return i.store[crag][0]
}

func (i *InMemoryCragStore) addForecast(crag, forecast string) {
	i.store[crag] = append(i.store[crag], forecast)
}
