package jsonschema

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidator_Validate_Int_Minimum(t *testing.T) {
	type Integer struct {
		Num int `jsonschema:"minimum:5"`
	}

	validator := NewValidator()

	// invalid
	n := Integer{
		Num: 4,
	}
	err := validator.Validate(n)
	assert.Error(t, err)
	assert.Equal(t, "minimum:5", err.Error())

	// valid
	n = Integer{
		Num: 5,
	}
	err = validator.Validate(n)
	assert.NoError(t, err)
}

func TestValidator_Validate_Int_Minimum_ExclusiveMinimum(t *testing.T) {
	type Integer struct {
		Num int `jsonschema:"minimum:5,exclusiveMinimum:true"`
	}

	validator := NewValidator()

	// invalid
	n := Integer{
		Num: 5,
	}
	err := validator.Validate(n)
	assert.Error(t, err)
	assert.Equal(t, "exclusiveMinimum:5", err.Error())

	// valid
	n = Integer{
		Num: 6,
	}
	err = validator.Validate(n)
	assert.NoError(t, err)
}

func TestValidator_Validate_Int_ExclusiveMinimum(t *testing.T) {
	type Integer struct {
		Num int `jsonschema:"exclusiveMinimum:5"`
	}

	validator := NewValidator()

	// invalid
	n := Integer{
		Num: 5,
	}
	err := validator.Validate(n)
	assert.Error(t, err)
	assert.Equal(t, "exclusiveMinimum:5", err.Error())

	// valid
	n = Integer{
		Num: 6,
	}
	err = validator.Validate(n)
	assert.NoError(t, err)
}

func TestValidator_Validate_Int_Maximum(t *testing.T) {
	type Integer struct {
		Num int `jsonschema:"maximum:5"`
	}

	validator := NewValidator()

	// invalid
	n := Integer{
		Num: 6,
	}
	err := validator.Validate(n)
	assert.Error(t, err)
	assert.Equal(t, "maximum:5", err.Error())

	// valid
	n = Integer{
		Num: 5,
	}
	err = validator.Validate(n)
	assert.NoError(t, err)
}

func TestValidator_Validate_Int_Maximum_ExclusiveMaximum(t *testing.T) {
	type Integer struct {
		Num int `jsonschema:"maximum:5,exclusiveMaximum:true"`
	}

	validator := NewValidator()

	// invalid
	n := Integer{
		Num: 5,
	}
	err := validator.Validate(n)
	assert.Error(t, err)
	assert.Equal(t, "exclusiveMaximum:5", err.Error())

	// valid
	n = Integer{
		Num: 4,
	}
	err = validator.Validate(n)
	assert.NoError(t, err)
}

func TestValidator_Validate_Int_ExclusiveMaximum(t *testing.T) {
	type Integer struct {
		Num int `jsonschema:"exclusiveMaximum:5"`
	}

	validator := NewValidator()

	// invalid
	n := Integer{
		Num: 5,
	}
	err := validator.Validate(n)
	assert.Error(t, err)
	assert.Equal(t, "exclusiveMaximum:5", err.Error())

	// valid
	n = Integer{
		Num: 4,
	}
	err = validator.Validate(n)
	assert.NoError(t, err)
}

func TestValidator_Validate_Int_MultipleOf(t *testing.T) {
	type Integer struct {
		Num int `jsonschema:"multipleOf:5"`
	}

	validator := NewValidator()

	// invalid
	n := Integer{
		Num: 4,
	}
	err := validator.Validate(n)
	assert.Error(t, err)
	assert.Equal(t, "multipleOf:5", err.Error())

	// valid
	n = Integer{
		Num: 15,
	}
	err = validator.Validate(n)
	assert.NoError(t, err)
}

func TestValidator_Validate_Float_MultipleOf(t *testing.T) {
	type Integer struct {
		Num float64 `jsonschema:"multipleOf:2.5"`
	}

	validator := NewValidator()

	// invalid
	n := Integer{
		Num: 4,
	}
	err := validator.Validate(n)
	assert.Error(t, err)
	assert.Equal(t, "multipleOf:2.5", err.Error())

	// valid
	n = Integer{
		Num: 7.5,
	}
	err = validator.Validate(n)
	assert.NoError(t, err)
}

func TestValidator_Validate_String_MaxLength(t *testing.T) {
	type String struct {
		Str string `jsonschema:"maxLength:5"`
	}

	validator := NewValidator()

	// invalid
	s := String{
		Str: "abcdef",
	}
	err := validator.Validate(s)
	assert.Error(t, err)
	assert.Equal(t, "maxLength:5(6)", err.Error())

	// valid
	s = String{
		Str: "abcde",
	}
	err = validator.Validate(s)
	assert.NoError(t, err)
}

func TestValidator_Validate_String_MinLength(t *testing.T) {
	type String struct {
		Str string `jsonschema:"minLength:3"`
	}

	validator := NewValidator()

	// invalid
	s := String{
		Str: "ab",
	}
	err := validator.Validate(s)
	assert.Error(t, err)
	assert.Equal(t, "minLength:3(2)", err.Error())

	// valid
	s = String{
		Str: "abc",
	}
	err = validator.Validate(s)
	assert.NoError(t, err)
}

func TestValidator_Validate_String_Pattern(t *testing.T) {
	type String struct {
		Str string `jsonschema:"pattern:[abc]+"`
	}

	validator := NewValidator()

	// invalid
	s := String{
		Str: "def",
	}
	err := validator.Validate(s)
	assert.Error(t, err)
	assert.Equal(t, "pattern:[abc]+(def)", err.Error())

	// valid
	s = String{
		Str: "abcd",
	}
	err = validator.Validate(s)
	assert.NoError(t, err)
}

func TestValidator_Validate_String_Format(t *testing.T) {
	type String struct {
		Str string `jsonschema:"format:ipv4"`
	}

	validator := NewValidator()

	// invalid
	s := String{
		Str: "abcd.adcd.adcd.abcd",
	}
	err := validator.Validate(s)
	assert.Error(t, err)
	assert.Equal(t, "", err.Error())

	// valid
	s = String{
		Str: "192.168.1.1",
	}
	err = validator.Validate(s)
	assert.NoError(t, err)
}

func TestValidator_Validate_Array_MaxItems(t *testing.T) {
	type Sample struct {
		Arr []int `jsonschema:"maxItems:3"`
	}

	validator := NewValidator()

	// invalid
	s := Sample{
		Arr: []int{1, 2, 3, 4},
	}
	err := validator.Validate(s)
	assert.Error(t, err)
	assert.Equal(t, "maxItems:3", err.Error())

	// valid
	s = Sample{
		Arr: []int{1, 2, 3},
	}
	err = validator.Validate(s)
	assert.NoError(t, err)
}

func TestValidator_Validate_Array_MinItems(t *testing.T) {
	type Sample struct {
		Arr []int `jsonschema:"minItems:3"`
	}

	validator := NewValidator()

	// invalid
	s := Sample{
		Arr: []int{1, 2},
	}
	err := validator.Validate(s)
	assert.Error(t, err)
	assert.Equal(t, "minItems:3", err.Error())

	// valid
	s = Sample{
		Arr: []int{1, 2, 1},
	}
	err = validator.Validate(s)
	assert.NoError(t, err)
}

func TestValidator_Validate_Array_UniqueItems(t *testing.T) {
	type Sample struct {
		Arr []int `jsonschema:"uniqueItems:true"`
	}

	validator := NewValidator()

	// invalid
	s := Sample{
		Arr: []int{1, 2, 1},
	}
	err := validator.Validate(s)
	assert.Error(t, err)
	assert.Equal(t, "uniqueItems:true", err.Error())

	// valid
	s = Sample{
		Arr: []int{1, 2, 3},
	}
	err = validator.Validate(s)
	assert.NoError(t, err)
}

func TestValidator_Validate_Map_MaxProperties(t *testing.T) {
	type Sample struct {
		Map map[string]int `jsonschema:"maxProperties:2"`
	}

	validator := NewValidator()

	// invalid
	s := Sample{
		Map: map[string]int{
			"a": 1,
			"b": 2,
			"c": 3,
		},
	}
	err := validator.Validate(s)
	assert.Error(t, err)
	assert.Equal(t, "maxProperties:2", err.Error())

	// valid
	s = Sample{
		Map: map[string]int{
			"a": 1,
			"b": 2,
		},
	}
	err = validator.Validate(s)
	assert.NoError(t, err)
}

func TestValidator_Validate_Map_MinProperties(t *testing.T) {
	type Sample struct {
		Map map[string]int `jsonschema:"minProperties:2"`
	}

	validator := NewValidator()

	// invalid
	s := Sample{
		Map: map[string]int{
			"a": 1,
		},
	}
	err := validator.Validate(s)
	assert.Error(t, err)
	assert.Equal(t, "minProperties:2", err.Error())

	// valid
	s = Sample{
		Map: map[string]int{
			"a": 1,
			"b": 2,
		},
	}
	err = validator.Validate(s)
	assert.NoError(t, err)
}

func TestValidator_Validate_Map_Required(t *testing.T) {
	type Sample struct {
		Map map[string]int `jsonschema:"required:[a,b,c]"`
	}

	validator := NewValidator()

	// invalid
	s := Sample{
		Map: map[string]int{
			"a": 1,
			"b": 2,
			"d": 3,
		},
	}
	err := validator.Validate(s)
	assert.Error(t, err)
	assert.Equal(t, "required", err.Error())

	// valid
	s = Sample{
		Map: map[string]int{
			"a": 1,
			"b": 2,
			"c": 3,
		},
	}
	err = validator.Validate(s)
	assert.NoError(t, err)
}

func TestValidator_Validate_Map_PatternProperties(t *testing.T) {
	type Sample struct {
		Map map[string]int `jsonschema:"patternProperties:^id-+"`
	}

	validator := NewValidator()

	// invalid
	s := Sample{
		Map: map[string]int{
			"id-a": 1,
			"id-b": 2,
			"d":    3,
		},
	}
	err := validator.Validate(s)
	assert.Error(t, err)
	assert.Equal(t, "patternProperties:^id-+", err.Error())

	// valid
	s = Sample{
		Map: map[string]int{
			"id-a": 1,
			"id-b": 2,
			"id-c": 3,
		},
	}
	err = validator.Validate(s)
	assert.NoError(t, err)
}
