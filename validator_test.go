package jsonschema

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type Sample struct {
	Name string `json:"name" jsonschema:"maxLength:5,pattern:^[0-8]+$"`
}

func Test_Validator_String(t *testing.T) {
	sample := &Sample{
		Name: "1234567890",
	}

	validator := NewValidator()
	err := validator.Validate(sample)
	assert.Error(t, err)
	assert.Equal(t, "maxLength:5(10)pattern:^[0-8]+$(1234567890)", err.Error())
}
