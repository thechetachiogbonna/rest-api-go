package utils

import (
	"encoding/json"
	"errors"
	"net/http"
)

func ParseJson(r *http.Request, v any) error {
	if r.Body == nil {
		return errors.New("Missing request body.")
	}

	return json.NewDecoder(r.Body).Decode(v)
}

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err any) {
	er, ok := err.(error)
	if !ok {
		WriteJson(w, status, err)
		return
	}

	WriteJson(w, status, FormatValidationError(er))
}
