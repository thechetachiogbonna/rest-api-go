package utils

import (
	"errors"
	"testing"
)

func TestFormatValidationErrorUsesJSONNamesForAnyStruct(t *testing.T) {
	payload := struct {
		DisplayName string `json:"display_name" validate:"required"`
		Code        string `json:"code" validate:"oneof=active inactive"`
	}{}

	err := Validate.Struct(payload)
	if err == nil {
		t.Fatal("expected validation error")
	}

	formatted := FormatValidationError(err)
	if got := formatted["display_name"]; got != "This field is required." {
		t.Errorf("display_name error = %q, want required message", got)
	}
	if got := formatted["code"]; got != "Failed on constraint: oneof" {
		t.Errorf("code error = %q, want generic constraint message", got)
	}
}

func TestFormatValidationErrorHandlesNonValidationError(t *testing.T) {
	err := errors.New("invalid request")
	formatted := FormatValidationError(err)

	if got := formatted["error"]; got != err.Error() {
		t.Errorf("error message = %q, want %q", got, err.Error())
	}
}
