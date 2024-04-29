package climb

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"testing"

	"github.com/gorilla/mux"
	"github.com/lregs/Crag/models"
	"github.com/lregs/Crag/util"
	"github.com/stretchr/testify/assert"
)

type MockClimbStore struct {
	climbs map[int]models.Climb
}

func (s *MockClimbStore) validatePayload(data models.ClimbPayload) error {
	if reflect.DeepEqual(models.ClimbPayload{}, data) {
		return errors.New("value is empty")
	}

	if data.Name == "" {
		return errors.New("climb must have name")
	}

	r, _ := regexp.Compile(`[6-9][abc]\+?$`)
	if !r.MatchString(data.Grade) {
		return errors.New("climb grade invalid ")
	}

	if data.CragID == 0 {
		return errors.New("invalid crag ID")
	}
	return nil
}

func (s *MockClimbStore) validateClimb(data models.Climb) error {
	if reflect.DeepEqual(models.Climb{}, data) {
		return errors.New("value is empty")
	}

	if data.Id == 0 {
		return errors.New("invalid id")
	}

	if data.Name == "" {
		return errors.New("climb must have name")
	}

	r, _ := regexp.Compile(`[6-9][abc]\+?$`)
	if !r.MatchString(data.Grade) {
		return errors.New("climb grade invalid ")
	}

	if data.CragID == 0 {
		return errors.New("invalid crag ID")
	}
	return nil
}

func (s *MockClimbStore) StoreClimb(climb models.ClimbPayload) (models.Climb, error) {

	if err := s.validatePayload(climb); err != nil {
		return models.Climb{}, fmt.Errorf("validatting payload failed %s", err)
	}

	s.climbs[1] = models.Climb{Id: 1, Name: climb.Name, Grade: climb.Grade, CragID: climb.CragID}
	return s.climbs[1], nil
}

func (s *MockClimbStore) GetClimbsByCragId(CragId int) ([]models.Climb, error) {
	return []models.Climb{{Id: 1, Name: "Harvey Oswald", Grade: "7a+", CragID: 1}}, nil
}

// why are we not returning an error here?!
func (s *MockClimbStore) GetAllClimbs() ([]models.Climb, error) {
	res := []models.Climb{}
	for _, crag := range s.climbs {
		res = append(res, crag)
	}
	return res, nil
}

func (s MockClimbStore) GetClimbById(Id int) (models.Climb, error) {

	if Id != 1 {
		return models.Climb{}, errors.New("Invalid Id ")
	}
	return s.climbs[Id], nil
}

func (s MockClimbStore) UpdateClimb(climb models.Climb) (models.Climb, error) {
	//not exactly a mock...
	err := s.validateClimb(climb)
	if err != nil {
		return models.Climb{}, errors.New("Climb could not be validated")
	}
	return climb, nil
}

func (s *MockClimbStore) DeleteClimb(Id int) (models.Climb, error) {

	_, ok := s.climbs[Id]
	if ok {
		delete(s.climbs, Id)
	} else {
		return models.Climb{}, errors.New("value does not exist in db")
	}

	delete(s.climbs, 1)
	return s.climbs[1], nil
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func TestPostClimb(t *testing.T) {

	store := returnPopulatedStore()
	handler := NewHandler(store)
	router := mux.NewRouter()
	router.PathPrefix("/climb").HandlerFunc(handler.Post()).Methods("POST")

	t.Run("Testing Valid Climb", func(t *testing.T) {
		payload := returnTestPayload()

		body, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("error marhsalling payload: %s", err)
		}

		res, req, err := util.NewPostRequest(body, "/climb")
		if err != nil {
			t.Fatalf("failed generating request@ %s", err)
		}

		router.ServeHTTP(res, req)

		assert.Equal(t, 201, res.Code)

		var resData models.Climb

		_, err = util.DecodeResponse(res.Body, &resData)
		if err != nil {
			t.Fatalf("decode failed: %s", err)
		}

		expectedVal := models.Climb{
			Id:     1,
			Name:   "the roof",
			Grade:  "7a",
			CragID: 1,
		}

		assert.Equal(t, expectedVal, resData)

	})

	t.Run("Testing Invalid Method", func(t *testing.T) {

		res, req := util.NewGetRequest("/climb")

		router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusMethodNotAllowed, res.Code)

	})
}

func TestGetById(t *testing.T) {

	store := returnPopulatedStore()
	handler := NewHandler(store)
	router := mux.NewRouter()
	router.PathPrefix("/climbs/{cragId}").HandlerFunc(handler.GetByCragId()).Methods("GET")

	t.Run("Testing Valid ID", func(t *testing.T) {

		res := httptest.NewRecorder()
		request, err := http.NewRequest("GET", "/climbs/2", nil)
		if err != nil {
			t.Fatalf("failed creating request %s", err)
		}
		router.ServeHTTP(res, request)

		assert.Equal(t, http.StatusOK, res.Code)

		var data []models.Climb
		if err := json.Unmarshal(res.Body.Bytes(), &data); err != nil {
			t.Fatalf("unmarshal failed, %s", err)
		}

		assert.Equal(t, returnTestClimb(), data[0])

	})

	t.Run("Testing Invalid URL", func(t *testing.T) {

		res := httptest.NewRecorder()
		request, err := http.NewRequest("GET", "/climbs/a", nil)
		if err != nil {
			t.Fatalf("failed creating request %s", err)
		}

		router.ServeHTTP(res, request)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		var response map[string]string
		if err := json.Unmarshal(res.Body.Bytes(), &response); err != nil {
			t.Fatalf("decoding failed %s", err)
		}
		errorPattern := "^could\\ not\\ convert\\ key.*"

		matched, err := regexp.MatchString(errorPattern, response["Error"])
		if err != nil {
			t.Fatalf("regex failed, %s", err)
		}

		if !matched {
			t.Fatalf("res %s did not match regex", response["Error"])
		}

	})
}

// func TestGetClimbsByCrag(t *testing.T) {
// 	store := &MockClimbStore{
// 		climbs: map[int]*models.Climb{
// 			1: {
// 				Id:     1,
// 				Name:   "Harvey Oswald",
// 				Grade:  "v2",
// 				CragID: 1,
// 			},
// 			2: {
// 				Id:     2,
// 				Name:   "Harvey Oswald sds",
// 				Grade:  "v3",
// 				CragID: 1,
// 			},
// 		},
// 	}
// 	if store == nil || store.climbs == nil {
// 		t.Fatalf("store or store.climbs is nil")
// 	}

// 	// Cr   eate a new handler with the store
// 	handler := NewHandler(store)
// 	router := mux.NewRouter()
// 	router.PathPrefix("/climb/crag/{cragID}").HandlerFunc(handler.handleGetClimbsByCrag()).Methods("GET")

// 	testCases := []struct {
// 		Name              string
// 		CragId            int
// 		ExptectedResponse []*models.Climb
// 		ExpectedCode      int
// 	}{
// 		{
// 			Name:   "Valid Climb",
// 			CragId: 1,
// 			ExptectedResponse: []*models.Climb{
// 				&models.Climb{
// 					Id:     1,
// 					Name:   "Harvey Oswald",
// 					Grade:  "v2",
// 					CragID: 1,
// 				},
// 				&models.Climb{
// 					Id:     2,
// 					Name:   "Harvey Oswald sds",
// 					Grade:  "v3",
// 					CragID: 1,
// 				},
// 			},
// 			ExpectedCode: 200,
// 		},
// 		{
// 			Name:              "Invalid Climb",
// 			CragId:            0,
// 			ExptectedResponse: nil,
// 			ExpectedCode:      400,
// 		},
// 	}
// 	for _, tc := range testCases {
// 		t.Run(tc.Name, func(t *testing.T) {
// 			url := fmt.Sprintf("/climb/crag/%d", tc.CragId)
// 			request := util.NewGetRequest(url)
// 			response := httptest.NewRecorder()
// 			router.ServeHTTP(response, request)

// 			t.Logf("Testing case %s, code = %d", tc.Name, response.Code)

// 			switch {
// 			case tc.ExpectedCode == 400:
// 				var errRes ErrorResponse
// 				err := json.NewDecoder(response.Body).Decode(&errRes)
// 				if err != nil {
// 					t.Fatalf("could not decode response %s", err)
// 				}
// 				t.Log(errRes)

// 				// util.AssertStatus(t, response.Code, http.StatusBadRequest)
// 			case tc.ExpectedCode == 200:
// 				// var errRes ErrorResponse
// 				// err := json.NewDecoder(response.Body).Decode(&errRes)
// 				// if err != nil {
// 				// 	t.Fatalf("could not decode response: %s", err)
// 				// }
// 				// t.Log(errRes)

// 				var climbsAtCrag []*models.Climb
// 				util.AssertStatus(t, response.Code, http.StatusOK)
// 				_, err := util.DecodeResponse(response.Body, &climbsAtCrag)
// 				if err != nil {
// 					t.Fatal(err)
// 				}
// 				for i, climb := range climbsAtCrag {
// 					if !reflect.DeepEqual(climb, tc.ExptectedResponse[i]) {
// 						t.Fatalf("returned response not what was expected. Got %v, want %v", climb, tc.ExptectedResponse[i])
// 					}
// 				}
// 			}
// 		})
// 	}
// }

// func TestGetAllClimbs(t *testing.T) {
// 	store := &MockClimbStore{
// 		climbs: map[int]*models.Climb{
// 			1: &models.Climb{
// 				Id:     1,
// 				Name:   "Harvey Oswald",
// 				Grade:  "v2",
// 				CragID: 1,
// 			},
// 		},
// 	}
// 	if store == nil || store.climbs == nil {
// 		t.Fatalf("store or store.climbs is nil")
// 	}

// 	handler := NewHandler(store)
// 	router := mux.NewRouter()
// 	router.PathPrefix("/climb/all").HandlerFunc(handler.HandleGetAllClimbs()).Methods("GET")

// 	//would we eventually want to get all crag info at the same time so we can return what crag these climbs belong too also instead of cragID"

// 	// need to add more test cases
// 	testCases := []struct {
// 		Name              string
// 		ExptectedResponse []*models.Climb
// 		ExpectedCode      int
// 	}{
// 		{
// 			Name: "Get",
// 			ExptectedResponse: []*models.Climb{
// 				&models.Climb{
// 					Id:     1,
// 					Name:   "Harvey Oswald",
// 					Grade:  "v2",
// 					CragID: 1,
// 				},
// 			},
// 			ExpectedCode: 200,
// 		},
// 	}
// 	t.Run(fmt.Sprintf("testing,%s", testCases[0].Name), func(t *testing.T) {
// 		url := "/climb/all"
// 		request := util.NewGetRequest(url)
// 		response := httptest.NewRecorder()
// 		router.ServeHTTP(response, request)

// 		util.AssertStatus(t, response.Code, http.StatusOK)

// 		var climbs []*models.Climb

// 		_, err := util.DecodeResponse(response.Body, &climbs)
// 		if err != nil {
// 			t.Fatalf("Decoding Response failed because of err: %s", err)
// 		}
// 		if !reflect.DeepEqual(climbs[0].Name, "Harvey Oswald") {
// 			t.Fatalf("return value did not match expected")
// 		}

// 	})
// }

// func TestGetClimbById(t *testing.T) {
// 	store := &MockClimbStore{
// 		climbs: map[int]*models.Climb{
// 			1: &models.Climb{
// 				Id:     1,
// 				Name:   "Harvey Oswald",
// 				Grade:  "v2",
// 				CragID: 1,
// 			},
// 		},
// 	}
// 	if store == nil || store.climbs == nil {
// 		t.Fatalf("store or store.climbs is nil")
// 	}

// 	handler := NewHandler(store)
// 	router := mux.NewRouter()
// 	router.PathPrefix("/climb/{Id}").HandlerFunc(handler.HandleGetClimbById()).Methods("GET")

// 	testCases := []struct {
// 		Name              string
// 		Id                int
// 		InvId             string
// 		ExptectedResponse *models.Climb
// 		ExpectedCode      int
// 		ExpectedError     bool
// 	}{
// 		{
// 			Name:  "Correct Id",
// 			Id:    1,
// 			InvId: "",
// 			ExptectedResponse: &models.Climb{
// 				Id:     1,
// 				Name:   "Harvey Oswald",
// 				Grade:  "v2",
// 				CragID: 1,
// 			},
// 			ExpectedCode:  200,
// 			ExpectedError: false,
// 		},
// 		{
// 			Name:              "Invalid Id",
// 			Id:                100,
// 			InvId:             "",
// 			ExptectedResponse: nil,
// 			ExpectedCode:      400,
// 			ExpectedError:     true, //I guess we actually add the error when we know it

// 		},
// 		//In go it has to be an int, I guess the client could send the wrong type?
// 		// {
// 		// 	Name:              "Invalid ID type",
// 		// 	Id:                0,
// 		// 	InvId:             "One",
// 		// 	ExptectedResponse: nil,
// 		// 	ExpectedCode:      400,
// 		// 	ExpectedError:     true,
// 		// },
// 	}

// 	for _, tc := range testCases {
// 		t.Run(fmt.Sprintf("testing %s", tc.Name), func(t *testing.T) {
// 			// var url string
// 			// if tc.Id == 0 {
// 			// 	url = "/climb/one"
// 			// }
// 			url := fmt.Sprintf("/climb/%d", tc.Id)
// 			request := util.NewGetRequest(url)
// 			response := httptest.NewRecorder()

// 			router.ServeHTTP(response, request)

// 			util.AssertStatus(t, response.Code, tc.ExpectedCode)

// 			if tc.ExpectedCode == 200 {
// 				var climb *models.Climb
// 				_, err := util.DecodeResponse(response.Body, &climb)
// 				if err != nil {
// 					t.Fatalf("Could not decode response because of err: %s", err)
// 				}

// 				if !reflect.DeepEqual(climb, tc.ExptectedResponse) {
// 					t.Fatalf("responsed object did not meet ")
// 				}
// 			}

// 		})
// 	}
// }

// func TestUpdateClimbById(t *testing.T) {
// 	store := &MockClimbStore{
// 		climbs: map[int]*models.Climb{
// 			1: &models.Climb{
// 				Id:     1,
// 				Name:   "Harvey Oswald",
// 				Grade:  "v2",
// 				CragID: 1,
// 			},
// 		},
// 	}
// 	if store == nil || store.climbs == nil {
// 		t.Fatalf("store or store.climbs is nil")
// 	}

// 	t.Run("Testing Valid Update", func(t *testing.T) {
// 		handler := NewHandler(store)
// 		router := mux.NewRouter()
// 		router.PathPrefix("/climb/").HandlerFunc(handler.HandleUpdateClimb()).Methods("PUT")

// 		response := httptest.NewRecorder()

// 		updatedClimb := &models.Climb{
// 			Id:     1,
// 			Name:   "Harvey Oswald",
// 			Grade:  "v7",
// 			CragID: 1,
// 		}

// 		url := "/climb/"

// 		body, _ := json.Marshal(updatedClimb)

// 		request, err := util.NewPutRequest(body, url)
// 		if err != nil {
// 			t.Fatalf("failed generating push request becauseo of err : %s", err)
// 		}

// 		router.ServeHTTP(response, request)

// 		if response.Code == 400 {
// 			var putErr ErrorResponse
// 			err, _ := util.DecodeResponse(response.Body, putErr)
// 			t.Logf(err.Error)
// 		}

// 		util.AssertStatus(t, response.Code, http.StatusOK)

// 	})

// 	t.Run("Testing In-Valid Update", func(t *testing.T) {
// 		handler := NewHandler(store)
// 		router := mux.NewRouter()
// 		router.PathPrefix("/climb/").HandlerFunc(handler.HandleUpdateClimb()).Methods("PUT")

// 		response := httptest.NewRecorder()

// 		updatedClimb := &models.Climb{
// 			Id:     0,
// 			Name:   "",
// 			Grade:  "",
// 			CragID: 1,
// 		}

// 		url := "/climb/"

// 		body, _ := json.Marshal(updatedClimb)

// 		request, err := util.NewPutRequest(body, url)
// 		if err != nil {
// 			t.Fatalf("failed generating push request becauseo of err : %s", err)
// 		}

// 		router.ServeHTTP(response, request)

// 		//need to assert response
// 		//maybe I need to be asserting certain errors as well and making sure im getting the correct response code each time

// 		util.AssertStatus(t, response.Code, http.StatusBadRequest)

// 	})
// }

// func TestDelCrag(t *testing.T) {
// 	store := MockClimbStore{
// 		climbs: map[int]models.Climb{
// 			1: models.Climb{
// 				Id:     1,
// 				Name:   "Harvey Oswald",
// 				Grade:  "v2",
// 				CragID: 1,
// 			},
// 		},
// 	}
// 	if store == nil || store.climbs == nil {
// 		t.Fatalf("store or store.climbs is nil")
// 	}

// 	t.Run("Testing Valid Delete", func(t *testing.T) {
// 		handler := NewHandler(store)
// 		router := mux.NewRouter()
// 		router.PathPrefix("/climb/{Id}").HandlerFunc(handler.HandleDeleteClimb()).Methods("DELETE")

// 		response := httptest.NewRecorder()

// 		url := "/climb/1"

// 		request, _ := http.NewRequest(http.MethodDelete, url, nil)

// 		router.ServeHTTP(response, request)

// 		util.AssertStatus(t, response.Code, http.StatusNoContent)
// 	})

// 	t.Run("Testing Invalid ID", func(t *testing.T) {
// 		handler := NewHandler(store)
// 		router := mux.NewRouter()
// 		router.PathPrefix("/climb/{Id}").HandlerFunc(handler.HandleDeleteClimb()).Methods("DELETE")

// 		response := httptest.NewRecorder()

// 		url := "/climb/3"

// 		request, _ := http.NewRequest(http.MethodDelete, url, nil)

// 		router.ServeHTTP(response, request)

// 		util.AssertStatus(t, response.Code, http.StatusBadRequest)
// 	})
// }

func returnPopulatedStore() *MockClimbStore {
	return &MockClimbStore{
		climbs: map[int]models.Climb{
			1: models.Climb{
				Id:     1,
				Name:   "Harvey Oswald",
				Grade:  "7a+",
				CragID: 1,
			},
			2: models.Climb{
				Id:     2,
				Name:   "Slopers",
				Grade:  "7a",
				CragID: 1,
			},
		},
	}
}

type testServer struct {
	store   *MockClimbStore
	handler *Handler
	router  *mux.Router
}

func newServer() *testServer {

	return &testServer{
		handler: NewHandler(returnPopulatedStore()),
		router:  mux.NewRouter(),
	}
}

func returnTestClimb() models.Climb {
	return models.Climb{
		Id:     1,
		Name:   "Harvey Oswald",
		Grade:  "7a+",
		CragID: 1,
	}
}

func returnTestPayload() models.ClimbPayload {
	return models.ClimbPayload{
		Name:   "the roof",
		Grade:  "7a",
		CragID: 1,
	}
}
