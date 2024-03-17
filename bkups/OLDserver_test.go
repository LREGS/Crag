package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// var url = "http://localhost:6969/crags/"

type StubCragStore struct {
	Names []string
}

func (s *StubCragStore) addCrag(name string) {
	s.Names = append(s.Names, name)
}

func (s *StubCragStore) getNames() []string {
	return s.Names
}

func TestStoreCrag(t *testing.T) {
	StubStore := StubCragStore{Names: []string{}}
	t.Run("Records on Post", func(t *testing.T) {
		crag := "stanage"
		handler := NewServer(&StubStore)

		request := newPostCragRequst(crag)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

		if len(StubStore.Names) == 0 {
			t.Fatalf("there has been no successful entries into the store")
		}

		if StubStore.Names[len(StubStore.Names)-1] != crag {
			t.Fatalf("the last entry into names was %s but was meant to be %s", StubStore.Names[len(StubStore.Names)-1], crag)
		}

	})

}
func newPostCragRequst(name string) *http.Request {
	req, _ := http.NewRequest("POST", fmt.Sprintf("/crags/%s", name), nil)
	return req
}

func TestGetCrag(t *testing.T) {
	StubStore := StubCragStore{Names: []string{"stanage"}}
	t.Run("Returns Crag named Stanage", func(t *testing.T) {
		crag := "stanage"
		handler := NewServer(&StubStore)

		request := newGetCragRequest(crag)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "stanage")
	})
}
func newGetCragRequest(name string) *http.Request {
	req, _ := http.NewRequest("GET", fmt.Sprintf("/crags/%s", name), nil)
	return req
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body incorrect. Got %s, wanted %s", got, want)
	}

}

func assertStatus(t testing.TB, got, want int) {
	if got != want {
		t.Fatalf("expected %d but got %d", want, got)
	}

}