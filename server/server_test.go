package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lregs/Crag/models"
)

type MockCragStore struct {
	crags map[int]*models.Crag
}

func (cs *MockCragStore) GetCrag(Id int) (*models.Crag, error) {
	crag, ok := cs.crags[Id]
	if !ok {
		err := fmt.Sprintf("Crag with Id %d not found", Id)
		return nil, errors.New(err)
	}
	return crag, nil

}

func (cs *MockCragStore) DeleteCragByID(Id int) error {
	_, ok := cs.crags[Id]
	if !ok {
		err := errors.New("No crag with idfound")
		return err
	}
	delete(cs.crags, Id)
	return nil
}

func (cs *MockCragStore) UpdateCragValue(name string, crag models.Crag) error {
	return nil
}

func (cs *MockCragStore) StoreCrag(crag *models.Crag) (err error) {
	// cs.crags = append(cs.crags, crag)

	if crag == nil {
		return errors.New("Crag is empty")
	}

	if crag.Name == "" {
		return errors.New("Name field is empty")
	}
	_, ok := cs.crags[crag.Id]
	if !ok {
		cs.crags[crag.Id] = crag
		return nil
	}
	return errors.New("Crag already exists")
}

func TestGetCrag(t *testing.T) {

	store := &MockCragStore{
		crags: map[int]*models.Crag{
			1: {Id: 1, Name: "Stanage", Latitude: 40.01, Longitude: 40.11},
			2: {Id: 2, Name: "Milestone", Latitude: 41.01, Longitude: 41.11},
		},
	}
	srv := NewServer(store)

	testCases := []struct {
		name             string
		cragID           int
		expectedResponse string
	}{
		{
			name:             "Valid crag ID",
			cragID:           1,
			expectedResponse: `{"Id":1,"Name":"Stanage","Latitude":40.01,"Longitude":40.11}`,
		},
		{
			name:             "Non-existent crag ID",
			cragID:           3,
			expectedResponse: "", // Empty string as we expect an error
		},
	}

	for _, testcase := range testCases {
		t.Run("Testing Crag Get Request", func(t *testing.T) {
			t.Logf("testing %d", testcase.cragID)
			response := httptest.NewRecorder()
			request := newGetCragRequest(testcase.cragID)
			srv.ServeHTTP(response, request)

			if testcase.expectedResponse != "" {
				assertStatus(t, response.Code, http.StatusOK)
				assertResponseBody(t, strings.TrimSpace(response.Body.String()), testcase.expectedResponse)

			} else {
				assertStatus(t, response.Code, http.StatusNotFound)
			}

		})
	}

}

func TestPostCrag(t *testing.T) {
	store := &MockCragStore{crags: make(map[int]*models.Crag)}
	srv := NewServer(store)

	testCases := []models.Crag{
		{Id: 1, Name: "Stanage", Latitude: 1.111, Longitude: 1.222},
		{Id: 1, Name: "Dank", Latitude: 1.111, Longitude: 1.222},
		{Id: 2, Name: "", Latitude: 1.111, Longitude: 1.222},
	}

	for _, tc := range testCases {
		t.Run("Testing POST Crag", func(t *testing.T) {

			reqBody, err := json.Marshal(tc)
			if err != nil {
				t.Fatalf("could not marhsall because of err %s", err)
			}

			request, err := newPostRequest(reqBody)
			if err != nil {
				t.Fatalf("error getting new request: %s", err)
			}
			response := httptest.NewRecorder()
			srv.ServeHTTP(response, request)

			if tc.Name == "Stanage" {
				assertStatus(t, response.Code, http.StatusOK) //200
			}

			if tc.Name == "Dank" {
				assertStatus(t, response.Code, http.StatusBadRequest) //409
			}

			if tc.Name == "" {
				assertStatus(t, response.Code, http.StatusBadRequest) //400
			}

		})
	}

}

func TestDeleteCragByID(t *testing.T) {
	store := &MockCragStore{crags: map[int]*models.Crag{
		1: {Id: 1, Name: "Stanage", Latitude: 40.0101, Longitude: 40.1011},
	}}
	srv := NewServer(store)
	testCases := []int{1, 2}

	for _, tc := range testCases {
		t.Run("Testing Delete By Id", func(t *testing.T) {
			request := NewDeleteRequest(tc)
			response := httptest.NewRecorder()
			srv.ServeHTTP(response, request)

			assertStatus(t, response.Code, http.StatusOK)
		})
	}

}

func NewDeleteRequest(Id int) *http.Request {
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/crags/%d", Id), nil)
	return req
}.I

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("Did not get correct status, got %d, wanted %d", got, want)
	}
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response is wrong, got %s, wanted %s", got, want)
	}
}

func newGetCragRequest(Id int) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/crags/%d", Id), nil)
	return req
}

func newPostRequest(body []byte) (*http.Request, error) {

	//im sure im supposed to marshall in a different way but im not sure this doesnt seem right
	req, err := http.NewRequest(http.MethodPost, "/crags", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return req, nil

}
