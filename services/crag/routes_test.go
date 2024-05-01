package crag

import (
	"encoding/json"
	"testing"

	"github.com/gorilla/mux"
	"github.com/lregs/Crag/models"
	"github.com/lregs/Crag/util"
	"github.com/stretchr/testify/assert"
)

// wtf is this mock
type MockCragStore struct {
	crags map[int]models.Crag
}

func (cs *MockCragStore) GetCrag(Id int) (models.Crag, error) {
	return models.Crag{Id: 1, Name: "milestone", Latitude: 41.01, Longitude: 41.11}, nil
}

func (cs *MockCragStore) DeleteCragByID(Id int) (models.Crag, error) {
	return models.Crag{Id: 1, Name: "milestone", Latitude: 41.01, Longitude: 41.11}, nil
}

func (cs *MockCragStore) UpdateCrag(crag models.Crag) (models.Crag, error) {
	return models.Crag{Id: 1, Name: "milestone", Latitude: 41.01, Longitude: 41.11}, nil
}

func (cs *MockCragStore) StoreCrag(crag models.CragPayload) (models.Crag, error) {
	// cs.crags = append(cs.crags, crag)
	return models.Crag{Id: 1, Name: "milestone", Latitude: 41.01, Longitude: 41.11}, nil
}

func newStore() *MockCragStore {
	return &MockCragStore{
		crags: map[int]models.Crag{
			1: {Id: 1, Name: "Stanage", Latitude: 40.01, Longitude: 40.11},
			2: {Id: 2, Name: "Milestone", Latitude: 41.01, Longitude: 41.11},
		},
	}
}

func TestGetCrag(t *testing.T) {
	handler := NewHandler(newStore())
	router := mux.NewRouter()
	router.PathPrefix("/crags/{key}").HandlerFunc(handler.GetById()).Methods("GET")

	t.Run("valid request", func(t *testing.T) {

		res, req := util.NewGetRequest("/crags/2")
		router.ServeHTTP(res, req)

		if !assert.Equal(t, 200, res.Code) {
			var resErr map[string]string
			if err := json.Unmarshal(res.Body.Bytes(), &resErr); err != nil {
				t.Fatalf("decoding error failed %s", err)
			}
			t.Fatalf("reqused failed %s", resErr["error"])
		}

		var crag models.Crag
		if err := json.Unmarshal(res.Body.Bytes(), &crag); err != nil {
			t.Fatalf("decode failed %s", err)
		}

		testData := models.Crag{Id: 1, Name: "milestone", Latitude: 41.01, Longitude: 41.11}
		assert.Equal(t, testData, crag)

	})

}

// func TestDeleteCragByID(t *testing.T) {

// 	store := &MockCragStore{
// 		crags: map[int]*models.Crag{
// 			1: {Id: 1, Name: "Stanage", Latitude: 40.01, Longitude: 40.11},
// 			2: {Id: 2, Name: "Milestone", Latitude: 41.01, Longitude: 41.11},
// 		},
// 	}

// 	handler := NewHandler(store)
// 	router := mux.NewRouter()

// 	router.PathPrefix("/crags/{key}").HandlerFunc(handler.handleDelCragById()).Methods("DELETE")

// 	testcases := []struct {
// 		name             string
// 		CragId           int
// 		ExpectedResponse string
// 	}{
// 		{
// 			name:             "Valid crag ID",
// 			CragId:           1,
// 			ExpectedResponse: `{"message":"Crag with id 1 deleted"}`,
// 			//do we want to be recieving what was deleted and completing data validation on it in the app but also in the test?
// 		},
// 		{
// 			name:             "Invalid crag ID",
// 			CragId:           3,
// 			ExpectedResponse: "",
// 		},
// 	}

// 	for _, testcase := range testcases {
// 		t.Run("Testing Delete Crag", func(t *testing.T) {

// 			request := NewDeleteRequest(testcase.CragId)
// 			response := httptest.NewRecorder()
// 			router.ServeHTTP(response, request)

// 			assertStatus(t, response.Code, http.StatusOK)

// 		})
// 	}

// }

// func TestPostCrag(t *testing.T) {
// 	store := &MockCragStore{crags: make(map[int]*models.Crag)}
// 	handler := NewHandler(store)
// 	//I dont know if we want a router or just an instance of server - but I think using router makes the tests more specific to the handler
// 	router := mux.NewRouter()
// 	router.PathPrefix("/crags").HandlerFunc(handler.handlePostCrag()).Methods("POST")

// 	//I actually think we want testCases to be their own struct with the testdata to be
// 	//in their own struct as well so that we can have more granularity in the tests
// 	testCases := []models.Crag{
// 		{Id: 1, Name: "Stanage", Latitude: 1.111, Longitude: 1.222},
// 		{Id: 1, Name: "Dank", Latitude: 1.111, Longitude: 1.222},
// 		{Id: 2, Name: "", Latitude: 1.111, Longitude: 1.222},
// 	}

// 	for _, tc := range testCases {
// 		t.Run("Testing POST Crag", func(t *testing.T) {

// 			reqBody, err := json.Marshal(tc)
// 			if err != nil {
// 				t.Fatalf("could not marhsall because of err %s", err)
// 			}

// 			request, err := newPostRequest(reqBody, "/crags")
// 			if err != nil {
// 				t.Fatalf("error getting new request: %s", err)
// 			}
// 			response := httptest.NewRecorder()
// 			router.ServeHTTP(response, request)

// 			if tc.Name == "Stanage" {
// 				assertStatus(t, response.Code, http.StatusOK) //200
// 			}

// 			if tc.Name == "Dank" {
// 				assertStatus(t, response.Code, http.StatusBadRequest) //409
// 			}

// 			if tc.Name == "" {
// 				assertStatus(t, response.Code, http.StatusBadRequest) //400
// 			}

// 		})
// 	}
// }

// func newGetCragRequest(Id int) *http.Request {
// 	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/crags/%d", Id), nil)
// 	return req
// }

// func assertStatus(t testing.TB, got, want int) {
// 	t.Helper()
// 	if got != want {
// 		t.Errorf("Did not get correct status, got %d, wanted %d", got, want)
// 	}
// }

// func assertResponseBody(t testing.TB, got, want string) {
// 	t.Helper()
// 	if got != want {
// 		t.Errorf("response is wrong, got %s, wanted %s", got, want)
// 	}
// }

// func NewDeleteRequest(Id int) *http.Request {
// 	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/crags/%d", Id), nil)
// 	return req
// }
// func newPostRequest(body []byte, url string) (*http.Request, error) {

// 	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
// 	if err != nil {
// 		return nil, err
// 	}

// 	req.Header.Set("Content-Type", "application/json")

// 	return req, nil

// }
