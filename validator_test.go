package jsonschema

import (
	"testing"
)

type Sample struct {
	Name string `json:"name" jsonschema:"maxLength:20"`
}

func TestNewValidator(t *testing.T) {
	sample := &Sample{
		Name: "my name",
	}

	validator := NewValidator()
	ret, err := validator.Validate(sample)
	if ret.Valid() != true {
		t.Error("invalid")
	}
	if err != nil {
		t.Error("err")
	}
}
