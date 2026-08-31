package utils

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
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

func WriteError(w http.ResponseWriter, status int, err error) {
	if _, ok := err.(validator.ValidationErrors); ok {
		formattedErrors := FormatValidationError(err)

		WriteJson(w, status, formattedErrors)
		return
	}

	WriteJson(w, status, map[string]string{"error": err.Error()})
}
