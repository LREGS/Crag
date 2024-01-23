package httpclient

import (
	"net/http"
	"time"
)

func defaultClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
	}
}
