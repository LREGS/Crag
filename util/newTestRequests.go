package util

import (
	"bytes"
	"net/http"
)

func NewPostRequest(body []byte, url string) (*http.Request, error) {

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return req, nil

}
