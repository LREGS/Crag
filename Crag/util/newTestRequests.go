package util

import (
	"bytes"
	"net/http"
	"net/http/httptest"
)

func NewPostRequest(body []byte, url string) (*httptest.ResponseRecorder, *http.Request, error) {

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	httptest.NewRecorder()

	return httptest.NewRecorder(), req, nil

}

func NewGetRequest(url string) (*httptest.ResponseRecorder, *http.Request) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	return httptest.NewRecorder(), req
}

func NewPutRequest(body []byte, url string) (*httptest.ResponseRecorder, *http.Request, error) {
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return httptest.NewRecorder(), req, nil
}
