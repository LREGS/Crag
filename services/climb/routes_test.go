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

	return returnTestClimb(), nil
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

func TestGetByCragId(t *testing.T) {

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

	t.Run("Invalid Method", func(t *testing.T) {

		res := httptest.NewRecorder()
		request, err := http.NewRequest("POST", "/climbs/2", nil)
		if err != nil {
			t.Fatalf("failed making request")
		}

		router.ServeHTTP(res, request)

		assert.Equal(t, 405, res.Code)

	})
}

func TestGetAll(t *testing.T) {
	store := returnPopulatedStore()
	handler := NewHandler(store)
	router := mux.NewRouter()
	router.PathPrefix("/climbs/{cragId}").HandlerFunc(handler.GetAll()).Methods("GET")

	t.Run("Valid Request", func(t *testing.T) {
		res, req := util.NewGetRequest("/climbs/1")
		router.ServeHTTP(res, req)

		if !assert.Equal(t, 200, res.Code) {
			var errorResponse map[string]string
			if err := json.Unmarshal(res.Body.Bytes(), &errorResponse); err != nil {
				t.Fatalf("decoding error failed %s", err)
			}
			t.Fatalf("request failed: %s", errorResponse["error"])
		}

		var resData []models.Climb
		if err := json.Unmarshal(res.Body.Bytes(), &resData); err != nil {
			t.Fatalf("decoding failed  %s", err)
		}

		testData := []models.Climb{
			{Id: 1,
				Name:   "Harvey Oswald",
				Grade:  "7a+",
				CragID: 1,
			},
			{
				Id:     2,
				Name:   "Slopers",
				Grade:  "7a",
				CragID: 1,
			},
		}
		assert.Equal(t, testData, resData)

	})

	t.Run("Invalid Request", func(t *testing.T) {

		res := httptest.NewRecorder()
		req, err := http.NewRequest("DELETE", "/climbs/1", nil)
		if err != nil {
			t.Fatalf("failed creating request: %s", err)
		}

		router.ServeHTTP(res, req)

		assert.Equal(t, 405, res.Code)

	})

}

func TestGetClimb(t *testing.T) {

	store := returnPopulatedStore()
	handler := NewHandler(store)
	router := mux.NewRouter()

	router.PathPrefix("/{Id}").HandlerFunc(handler.GetById()).Methods("GET")

	t.Run("Valid Request", func(t *testing.T) {
		res, req := util.NewGetRequest("/1")
		router.ServeHTTP(res, req)

		if !assert.Equal(t, 200, res.Code) {
			var errResponse map[string]string
			if err := json.Unmarshal(res.Body.Bytes(), &errResponse); err != nil {
				t.Fatalf("decoding failed: %s", err)
			}
		}

		var resData models.Climb
		if err := json.Unmarshal(res.Body.Bytes(), &resData); err != nil {
			t.Fatalf("deocding failed %s", err)
		}
		assert.Equal(t, returnTestClimb(), resData)

	})

	t.Run("Invalid Request", func(t *testing.T) {
		res := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/1", nil)
		if err != nil {
			t.Fatalf("creating request failed %s", err)
		}

		router.ServeHTTP(res, req)

		assert.Equal(t, 405, res.Code)

	})

}

func TestUpdateClimbById(t *testing.T) {
	store := returnPopulatedStore()
	handler := NewHandler(store)
	router := mux.NewRouter()
	router.PathPrefix("/climb").HandlerFunc(handler.Update()).Methods("PUT")

	t.Run("Testing Valid Update", func(t *testing.T) {

		updatedClimb := models.Climb{
			Id:     1,
			Name:   "Harvey Oswald",
			Grade:  "7a+",
			CragID: 1,
		}

		body, err := json.Marshal(updatedClimb)
		if err != nil {
			t.Fatalf("encoding err %s", err)
		}

		res, req, err := util.NewPutRequest(body, "/climb")
		if err != nil {
			t.Fatalf("failed generating push request becauseo of err : %s", err)
		}

		router.ServeHTTP(res, req)

		if !assert.Equal(t, 200, res.Code) {
			var errRes map[string]string
			if err := json.Unmarshal(res.Body.Bytes(), &errRes); err != nil {
				t.Fatalf("decoding failed %s", err)
			}
			t.Fatalf("code was %d, wanted 200. Error: %s", res.Code, errRes["Error"])

		}

		var resData models.Climb
		if err := json.Unmarshal(res.Body.Bytes(), &resData); err != nil {
			t.Fatalf("decoding failed %s", err)
		}

		assert.Equal(t, updatedClimb, resData)

	})

	// t.Run("Testing In-Valid Update", func(t *testing.T) {
	// 	handler := NewHandler(store)
	// 	router := mux.NewRouter()
	// 	router.PathPrefix("/climb/").HandlerFunc(handler.HandleUpdateClimb()).Methods("PUT")

	// 	response := httptest.NewRecorder()

	// 	updatedClimb := &models.Climb{
	// 		Id:     0,
	// 		Name:   "",
	// 		Grade:  "",
	// 		CragID: 1,
	// 	}

	// 	url := "/climb/"

	// 	body, _ := json.Marshal(updatedClimb)

	// 	request, err := util.NewPutRequest(body, url)
	// 	if err != nil {
	// 		t.Fatalf("failed generating push request becauseo of err : %s", err)
	// 	}

	// 	router.ServeHTTP(response, request)

	// 	//need to assert response
	// 	//maybe I need to be asserting certain errors as well and making sure im getting the correct response code each time

	// 	util.AssertStatus(t, response.Code, http.StatusBadRequest)

	// })
}

func TestDelCrag(t *testing.T) {
	store := returnPopulatedStore()
	handler := NewHandler(store)
	router := mux.NewRouter()
	router.PathPrefix("/climb/{Id}").HandlerFunc(handler.Delete()).Methods("DELETE")

	t.Run("Testing Valid Delete", func(t *testing.T) {

		res := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodDelete, "/climb/1", nil)

		router.ServeHTTP(res, req)

		testData := returnTestClimb()
		var resData models.Climb

		if err := json.Unmarshal(res.Body.Bytes(), &resData); err != nil {
			t.Fatalf("decoding failed %s", err)
		}

		assert.Equal(t, testData, resData)

	})

}

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
