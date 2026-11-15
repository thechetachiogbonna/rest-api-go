package utils

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParseJsonRejectsUnknownFields(t *testing.T) {
	request := httptest.NewRequest("PUT", "/", strings.NewReader(`{"password":"secret"}`))
	payload := struct {
		FirstName string `json:"first_name"`
	}{}

	err := ParseJson(request, &payload)
	if err == nil {
		t.Fatal("expected unknown field to be rejected")
	}
}
