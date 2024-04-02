package climb

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
	"github.com/lregs/Crag/models"
	"github.com/lregs/Crag/util"
)

type MockClimbStore struct {
	climbs map[int]*models.Climb
}

func (s *MockClimbStore) Validate(climb *models.Climb) error {
	if climb.Id < 1 {
		return errors.New("Id is not valid")
	}
	if climb.Name == "" {
		return errors.New("Name not allowed to be empty")
	}
	// gradeRange := make(map[string]bool)
	// for i := 0; i <= 17; i++ {
	// 	gradeRange[fmt.Sprintf("v%d", i)] = true
	// }
	// if _, ok := gradeRange[climb.Grade]; !ok {
	// 	return errors.New("No French")
	// }
	//maybe we would want to validate this by actually checking if this cragID exists
	if climb.CragID < 1 {
		return errors.New("CragId is not valid")
	}
	return nil
}

func (s *MockClimbStore) StoreClimb(climb *models.Climb) (*models.Climb, error) {
	err := s.Validate(climb)
	if err != nil {
		return nil, err
	}
	s.climbs[1] = climb
	return s.climbs[1], nil
}

func (s *MockClimbStore) GetClimbsByCrag(CragId int) ([]*models.Climb, error) {
	if CragId == 0 {
		return nil, errors.New("No climbs")
	}
	res := []*models.Climb{}
	for _, crag := range s.climbs {
		res = append(res, crag)
	}
	return res, nil
}

// why are we not returning an error here?!
func (s *MockClimbStore) GetAllClimbs() []*models.Climb {
	res := []*models.Climb{}
	for _, crag := range s.climbs {
		res = append(res, crag)
	}
	return res
}

func (s MockClimbStore) GetClimbById(Id int) (*models.Climb, error) {
	if Id != 1 {
		return nil, errors.New("Invalid Id ")
	}
	return s.climbs[Id], nil
}

func (s MockClimbStore) UpdateClimb(climb *models.Climb) (*models.Climb, error) {
	//not exactly a mock...
	err := s.Validate(climb)
	if err != nil {
		return nil, errors.New("Climb could not be validated")
	}
	return climb, nil
}

// pretty sure we want to be returning an instance of the delete climb for data validation
func (s *MockClimbStore) DeleteClimb(Id int) error {

	_, ok := s.climbs[Id]
	if ok {
		delete(s.climbs, Id)
	} else {
		return errors.New("value does not exist in db")
	}

	delete(s.climbs, 1)
	return nil
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func TestPostClimb(t *testing.T) {
	//why arent we testing: incorrect model entirely, incorrect request type? must be more cases we can test too
	// Check if store is not nil
	store := &MockClimbStore{
		climbs: make(map[int]*models.Climb),
	}
	if store.climbs == nil {
		t.Fatalf("store.climbs is nil")
	}

	// Create a new handler with the store
	handler := NewHandler(store)
	router := mux.NewRouter()
	router.PathPrefix("/climb").HandlerFunc(handler.handlePostClimb()).Methods("POST")

	// Test cases
	testCases := []struct {
		Name             string
		model            *models.Climb
		expectedResponse *models.Climb
		expectedCode     int
	}{
		{
			Name: "Valid Climb",
			model: &models.Climb{
				Id:     1,
				Name:   "Harvey Oswald",
				Grade:  "v2",
				CragID: 1,
			},
			expectedResponse: &models.Climb{
				Id:     1,
				Name:   "Harvey Oswald",
				Grade:  "v2",
				CragID: 1,
			},
			expectedCode: 201,
		},
		{
			Name: "Invalid Climb Null values in Model",
			model: &models.Climb{
				Id:     0,
				Name:   "",
				Grade:  "",
				CragID: 0,
			},
			expectedResponse: nil,
			expectedCode:     400,
		},
	}

	// Loop over test cases
	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			// Marshal the model to JSON
			body, err := json.Marshal(test.model)
			if err != nil {
				t.Fatalf("Failed to marshal model: %s", err)
			}

			// Create a new POST request
			request, err := util.NewPostRequest(body, "/climb")
			if err != nil {
				t.Fatalf("Failed to create request: %s", err)
			}

			// Create a new response recorder
			response := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(response, request)

			// Assert the response status
			switch {
			case test.expectedCode == 400:
				util.AssertStatus(t, response.Code, http.StatusBadRequest)
			case test.expectedCode == 201:
				util.AssertStatus(t, response.Code, http.StatusCreated)
				responseBody, err := util.DecodeResponse(response.Body, &models.Climb{})
				if err != nil {
					t.Fatalf("Failed to decode response: %s", err)
				}
				// Check if the response body is equal to the expected response
				if !reflect.DeepEqual(test.expectedResponse, responseBody) {
					t.Fatalf("expected: %v, got: %v", test.expectedResponse, responseBody)
				}
			}

		})
	}
}

func TestGetClimbsByCrag(t *testing.T) {
	store := &MockClimbStore{
		climbs: map[int]*models.Climb{
			1: {
				Id:     1,
				Name:   "Harvey Oswald",
				Grade:  "v2",
				CragID: 1,
			},
			2: {
				Id:     2,
				Name:   "Harvey Oswald sds",
				Grade:  "v3",
				CragID: 1,
			},
		},
	}
	if store == nil || store.climbs == nil {
		t.Fatalf("store or store.climbs is nil")
	}

	// Cr   eate a new handler with the store
	handler := NewHandler(store)
	router := mux.NewRouter()
	router.PathPrefix("/climb/crag/{cragID}").HandlerFunc(handler.handleGetClimbsByCrag()).Methods("GET")

	testCases := []struct {
		Name              string
		CragId            int
		ExptectedResponse []*models.Climb
		ExpectedCode      int
	}{
		{
			Name:   "Valid Climb",
			CragId: 1,
			ExptectedResponse: []*models.Climb{
				&models.Climb{
					Id:     1,
					Name:   "Harvey Oswald",
					Grade:  "v2",
					CragID: 1,
				},
				&models.Climb{
					Id:     2,
					Name:   "Harvey Oswald sds",
					Grade:  "v3",
					CragID: 1,
				},
			},
			ExpectedCode: 200,
		},
		{
			Name:              "Invalid Climb",
			CragId:            0,
			ExptectedResponse: nil,
			ExpectedCode:      400,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			url := fmt.Sprintf("/climb/crag/%d", tc.CragId)
			request := util.NewGetRequest(url)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)

			t.Logf("Testing case %s, code = %d", tc.Name, response.Code)

			switch {
			case tc.ExpectedCode == 400:
				var errRes ErrorResponse
				err := json.NewDecoder(response.Body).Decode(&errRes)
				if err != nil {
					t.Fatalf("could not decode response %s", err)
				}
				t.Log(errRes)

				// util.AssertStatus(t, response.Code, http.StatusBadRequest)
			case tc.ExpectedCode == 200:
				// var errRes ErrorResponse
				// err := json.NewDecoder(response.Body).Decode(&errRes)
				// if err != nil {
				// 	t.Fatalf("could not decode response: %s", err)
				// }
				// t.Log(errRes)

				var climbsAtCrag []*models.Climb
				util.AssertStatus(t, response.Code, http.StatusOK)
				_, err := util.DecodeResponse(response.Body, &climbsAtCrag)
				if err != nil {
					t.Fatal(err)
				}
				for i, climb := range climbsAtCrag {
					if !reflect.DeepEqual(climb, tc.ExptectedResponse[i]) {
						t.Fatalf("returned response not what was expected. Got %v, want %v", climb, tc.ExptectedResponse[i])
					}
				}
			}
		})
	}
}

func TestGetAllClimbs(t *testing.T) {
	store := &MockClimbStore{
		climbs: map[int]*models.Climb{
			1: &models.Climb{
				Id:     1,
				Name:   "Harvey Oswald",
				Grade:  "v2",
				CragID: 1,
			},
		},
	}
	if store == nil || store.climbs == nil {
		t.Fatalf("store or store.climbs is nil")
	}

	handler := NewHandler(store)
	router := mux.NewRouter()
	router.PathPrefix("/climb/all").HandlerFunc(handler.HandleGetAllClimbs()).Methods("GET")

	//would we eventually want to get all crag info at the same time so we can return what crag these climbs belong too also instead of cragID"

	// need to add more test cases
	testCases := []struct {
		Name              string
		ExptectedResponse []*models.Climb
		ExpectedCode      int
	}{
		{
			Name: "Get",
			ExptectedResponse: []*models.Climb{
				&models.Climb{
					Id:     1,
					Name:   "Harvey Oswald",
					Grade:  "v2",
					CragID: 1,
				},
			},
			ExpectedCode: 200,
		},
	}
	t.Run(fmt.Sprintf("testing,%s", testCases[0].Name), func(t *testing.T) {
		url := "/climb/all"
		request := util.NewGetRequest(url)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		util.AssertStatus(t, response.Code, http.StatusOK)

		var climbs []*models.Climb

		_, err := util.DecodeResponse(response.Body, &climbs)
		if err != nil {
			t.Fatalf("Decoding Response failed because of err: %s", err)
		}
		if !reflect.DeepEqual(climbs[0].Name, "Harvey Oswald") {
			t.Fatalf("return value did not match expected")
		}

	})
}

func TestGetClimbById(t *testing.T) {
	store := &MockClimbStore{
		climbs: map[int]*models.Climb{
			1: &models.Climb{
				Id:     1,
				Name:   "Harvey Oswald",
				Grade:  "v2",
				CragID: 1,
			},
		},
	}
	if store == nil || store.climbs == nil {
		t.Fatalf("store or store.climbs is nil")
	}

	handler := NewHandler(store)
	router := mux.NewRouter()
	router.PathPrefix("/climb/{Id}").HandlerFunc(handler.HandleGetClimbById()).Methods("GET")

	testCases := []struct {
		Name              string
		Id                int
		InvId             string
		ExptectedResponse *models.Climb
		ExpectedCode      int
		ExpectedError     bool
	}{
		{
			Name:  "Correct Id",
			Id:    1,
			InvId: "",
			ExptectedResponse: &models.Climb{
				Id:     1,
				Name:   "Harvey Oswald",
				Grade:  "v2",
				CragID: 1,
			},
			ExpectedCode:  200,
			ExpectedError: false,
		},
		{
			Name:              "Invalid Id",
			Id:                100,
			InvId:             "",
			ExptectedResponse: nil,
			ExpectedCode:      400,
			ExpectedError:     true, //I guess we actually add the error when we know it

		},
		//In go it has to be an int, I guess the client could send the wrong type?
		// {
		// 	Name:              "Invalid ID type",
		// 	Id:                0,
		// 	InvId:             "One",
		// 	ExptectedResponse: nil,
		// 	ExpectedCode:      400,
		// 	ExpectedError:     true,
		// },
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("testing %s", tc.Name), func(t *testing.T) {
			// var url string
			// if tc.Id == 0 {
			// 	url = "/climb/one"
			// }
			url := fmt.Sprintf("/climb/%d", tc.Id)
			request := util.NewGetRequest(url)
			response := httptest.NewRecorder()

			router.ServeHTTP(response, request)

			util.AssertStatus(t, response.Code, tc.ExpectedCode)

			if tc.ExpectedCode == 200 {
				var climb *models.Climb
				_, err := util.DecodeResponse(response.Body, &climb)
				if err != nil {
					t.Fatalf("Could not decode response because of err: %s", err)
				}

				if !reflect.DeepEqual(climb, tc.ExptectedResponse) {
					t.Fatalf("responsed object did not meet ")
				}
			}

		})
	}
}

func TestUpdateClimbById(t *testing.T) {
	store := &MockClimbStore{
		climbs: map[int]*models.Climb{
			1: &models.Climb{
				Id:     1,
				Name:   "Harvey Oswald",
				Grade:  "v2",
				CragID: 1,
			},
		},
	}
	if store == nil || store.climbs == nil {
		t.Fatalf("store or store.climbs is nil")
	}

	t.Run("Testing Valid Update", func(t *testing.T) {
		handler := NewHandler(store)
		router := mux.NewRouter()
		router.PathPrefix("/climb/").HandlerFunc(handler.HandleUpdateClimb()).Methods("PUT")

		response := httptest.NewRecorder()

		updatedClimb := &models.Climb{
			Id:     1,
			Name:   "Harvey Oswald",
			Grade:  "v7",
			CragID: 1,
		}

		url := "/climb/"

		body, _ := json.Marshal(updatedClimb)

		request, err := util.NewPutRequest(body, url)
		if err != nil {
			t.Fatalf("failed generating push request becauseo of err : %s", err)
		}

		router.ServeHTTP(response, request)

		if response.Code == 400 {
			var putErr ErrorResponse
			err, _ := util.DecodeResponse(response.Body, putErr)
			t.Logf(err.Error)
		}

		util.AssertStatus(t, response.Code, http.StatusOK)

	})

	t.Run("Testing In-Valid Update", func(t *testing.T) {
		handler := NewHandler(store)
		router := mux.NewRouter()
		router.PathPrefix("/climb/").HandlerFunc(handler.HandleUpdateClimb()).Methods("PUT")

		response := httptest.NewRecorder()

		updatedClimb := &models.Climb{
			Id:     0,
			Name:   "",
			Grade:  "",
			CragID: 1,
		}

		url := "/climb/"

		body, _ := json.Marshal(updatedClimb)

		request, err := util.NewPutRequest(body, url)
		if err != nil {
			t.Fatalf("failed generating push request becauseo of err : %s", err)
		}

		router.ServeHTTP(response, request)

		//need to assert response
		//maybe I need to be asserting certain errors as well and making sure im getting the correct response code each time

		util.AssertStatus(t, response.Code, http.StatusBadRequest)

	})
}

func TestDelCrag(t *testing.T) {
	store := &MockClimbStore{
		climbs: map[int]*models.Climb{
			1: &models.Climb{
				Id:     1,
				Name:   "Harvey Oswald",
				Grade:  "v2",
				CragID: 1,
			},
		},
	}
	if store == nil || store.climbs == nil {
		t.Fatalf("store or store.climbs is nil")
	}

	t.Run("Testing Valid Delete", func(t *testing.T) {
		handler := NewHandler(store)
		router := mux.NewRouter()
		router.PathPrefix("/climb/{Id}").HandlerFunc(handler.HandleDeleteClimb()).Methods("DELETE")

		response := httptest.NewRecorder()

		url := "/climb/1"

		request, _ := http.NewRequest(http.MethodDelete, url, nil)

		router.ServeHTTP(response, request)

		util.AssertStatus(t, response.Code, http.StatusNoContent)
	})

	t.Run("Testing Invalid ID", func(t *testing.T) {
		handler := NewHandler(store)
		router := mux.NewRouter()
		router.PathPrefix("/climb/{Id}").HandlerFunc(handler.HandleDeleteClimb()).Methods("DELETE")

		response := httptest.NewRecorder()

		url := "/climb/3"

		request, _ := http.NewRequest(http.MethodDelete, url, nil)

		router.ServeHTTP(response, request)

		util.AssertStatus(t, response.Code, http.StatusBadRequest)
	})
}
