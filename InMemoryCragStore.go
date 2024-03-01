package main

func NewInMemoryCragStore() *InMemoryCragStore {
	return &InMemoryCragStore{[]string{}}
}

type InMemoryCragStore struct {
	Names []string
}

func (i *InMemoryCragStore) addCrag(name string) {
	i.Names = append(i.Names, name)
}
