package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lregs/Crag/models"
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

func WriteError(w http.ResponseWriter, status int, err error) {
	Encode(w, status, map[string]string{"error": err.Error()})
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

func RWriteResponse(w http.ResponseWriter, status int, data models.Response) {

	w.Header().Set("Content Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)

}

func DResponse(body *bytes.Buffer, r models.Response) (models.Response, error) {
	//maybe only needs to return an error becuase im getting passed a pointer to the value so am editing hte original value
	if body == nil {
		return nil, fmt.Errorf("decode failed: Body is nil")
	}

	d := json.NewDecoder(body)
	if err := d.Decode(r); err != nil {
		return nil, fmt.Errorf("error decoding response %s", err)
	}
	return r, nil
}
