package jsonschema

import (
	"github.com/stretchr/testify/assert"
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
	assert.True(t, ret.Valid())
	assert.NoError(t, err)
}
