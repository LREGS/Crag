package server

import (
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

func (cs *MockCragStore) UpdateCragValuej(name string, crag models.Crag) error {
	return nil
}
func (cs *MockCragStore) DeleteCragByID(Id int) error {
	return nil
}
func (cs *MockCragStore) UpdateCragValue(name string, crag models.Crag) error {
	return nil
}
func (cs *MockCragStore) StoreCrag(crag *models.Crag) (err error) {
	// cs.crags = append(cs.crags, crag)
	return nil
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

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("Did not get correct staus, got %d, wanted %d", got, want)
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
