package main_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	Crag "github.com/lregs/Crag"
)

var url = "http://localhost:8080/crags"

type StubCragStore struct {
	Names []string
}

func (s *StubCragStore) addCrag(crag string) {
	s.Names = append(s.Names, crag)
}

func TestStoreCrag(t *testing.T) {
	StubStore := &StubCragStore{Names: []string{}}
	t.Run("Records on Post", func(t *testing.T) {
		crag := "stanage"
		handler := Crag.NewServer(StubStore)

		request := newPostCragRequst(crag)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

		if len(StubStore.Names) == 0 {
			t.Fatalf("there has been no successful entries into the store")
		}

		if StubStore.Names[len(StubStore.Names)-1] != crag {
			t.Fatalf("the last entry into names was %s but was meant to be %s", StubStore.Names[len(StubStore.names)-1], crag)
		}

	})

}

func newPostCragRequst(name string) *http.Request {
	req, _ := http.NewRequest("POST", fmt.Sprintf("https://localhost:6969/crags/"+name), nil)
	return req
}

func assertStatus(t testing.TB, got, want int) {

}
