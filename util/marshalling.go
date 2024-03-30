package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func Encode(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	// if err := json.NewEncoder(w).Encode(v); err != nil {
	// 	return fmt.Errorf("encode json: %w", err)
	// }
	// return nil
	return json.NewEncoder(w).Encode(v)
}

func Decode(r *http.Request, v any) error {
	// // var v T
	// if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
	// 	return v, fmt.Errorf("decode json: %w", err)
	// }
	// return v, nil
	if r.Body == nil {
		return fmt.Errorf("missing body")
	}
	return json.NewDecoder(r.Body).Decode(v)
}

func DecodeResponse[T any](body *bytes.Buffer, v T) (T, error) {
	if body == nil {
		return v, fmt.Errorf("decode json: body is nil")
	}
	dec := json.NewDecoder(body)
	if err := dec.Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}

func WriteError(w http.ResponseWriter, status int, err error) {
	Encode(w, status, map[string]string{"error": err.Error()})
}
