package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Data  interface{}
	Error string
}

func Encode(w http.ResponseWriter, status int, v any) error {
	//why have I commented out the better code?!
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

	defer r.Body.Close()

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

// nake this take err string and err and just do the fmt.Errorf inside this function rather than in every function call
// but maybe we start handling errors in other ways
func WriteError(w http.ResponseWriter, status int, errStr string, err error) {
	w.WriteHeader(status)

	errRes := map[string]error{"Error": fmt.Errorf(errStr, err)}

	json.NewEncoder(w).Encode(errRes)
}

func WriteResponse(w http.ResponseWriter, status int, data any, err string) {
	response := &Response{
		Data:  data,
		Error: err,
	}

	w.Header().Set("Content Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)

}
