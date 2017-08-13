package jsonschema

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidator_Validate_String(t *testing.T) {
	type String struct {
		Str string `json:"str" jsonschema:"maxLength:5,pattern:^[a-z]+$"`
	}
	s := String{
		Str: "1234567890",
	}

	validator := NewValidator()
	err := validator.Validate(s)
	assert.Error(t, err)
	assert.Equal(t, "maxLength:5(10)pattern:^[a-z]+$(1234567890)", err.Error())
}

func TestValidator_Validate_Int(t *testing.T) {
	type Integer struct {
		Num int `json:"num" jsonschema:"minimum:5"`
	}
	n := Integer{
		Num: 4,
	}

	validator := NewValidator()
	err := validator.Validate(n)
	assert.Error(t, err)
	assert.Equal(t, "minimum:5", err.Error())
}
