package main

type cragStore interface {
	addCrag(crag string)
}

func NewInMemoryCragStore() *InMemoryCragStore {
	return &InMemoryCragStore{[]string{}}
}

type InMemoryCragStore struct {
	Names []string
}

func (i *InMemoryCragStore) addCrag(name string) {
	i.Names = append(i.Names, name)
}
