package httpclient

import (
	"net/http"
	"time"
)

func DefaultClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
	}
}
