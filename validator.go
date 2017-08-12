package jsonschema

import (
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"unicode"
	"unicode/utf8"
)

const (
	tagName = "jsonschema"
)

type ValidationResult struct {
	Name    string
	Meta    string
	Message string
	Causes  []*ValidationResult
}

func (v *ValidationResult) Valid() bool {
	return v.isEmpty()
}

func (v *ValidationResult) isEmpty() bool {
	return v.Name == "" && v.Meta == "" && v.Message == "" && len(v.Causes) == 0
}

func (v *ValidationResult) add(err *ValidationResult) {
	if err != nil {
		v.Causes = append(v.Causes, err)
	}
}

func newValidationResult() *ValidationResult {
	return &ValidationResult{
		Causes: []*ValidationResult{},
	}
}

// NewValidator -
func NewValidator() *Validator {
	return &Validator{
		formats: map[string]ValidateFunc{
			// Defined formats
			"date-time":     dataTime,
			"email":         email,
			"hostname":      hostname,
			"ipv4":          ipv4,
			"ipv6":          ipv6,
			"uri":           uri,
			"uri-reference": uriReference,
			"uri-template":  uriReference,
			"json-pointer":  jsonPointer,
		},
	}
}

// ValidateFunc -
type ValidateFunc func(data string) error

// Validator -
type Validator struct {
	formats map[string]ValidateFunc
}

// AddFormat -
func (v *Validator) AddFormat(key string, f ValidateFunc) error {
	if key == "" || f == nil {
		return errors.New("")
	}
	if _, ok := v.formats[key]; ok {
		return errors.New("")
	}
	v.formats[key] = f
	return nil
}

func (v *Validator) execFormat(key, value string) error {
	f, ok := v.formats[key]
	if !ok {
		return errors.New("")
	}
	return f(value)
}

// Validate -
func (v *Validator) Validate(data interface{}) (*ValidationResult, error) {
	rt, rv := reflect.TypeOf(data), reflect.ValueOf(data)

	if rv.Kind() == reflect.Ptr && !rv.IsNil() {
		return v.Validate(rv.Elem().Interface())
	}
	if rv.Kind() != reflect.Struct && rv.Kind() != reflect.Interface {
		return nil, errors.New("")
	}

	ret := newValidationResult()
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		name := field.Name
		if !unicode.IsUpper(rune(name[0])) {
			continue
		}

		tagValue := field.Tag.Get(tagName)
		if tagValue == "-" {
			continue
		}

		tag := v.parseTag(tagValue)

		value := rv.Field(i)
		if value.Kind() == reflect.Ptr && !value.IsNil() {
			value = value.Elem()
		}

		if ve := v.validate(value, name, tag); ve != nil && !ve.isEmpty() {
			ret.add(ve)
		}
	}

	return ret, nil
}

func (v *Validator) parseTag(meta string) *tag {
	tag := newTag(meta)
	r := newReader([]byte(meta))
	tag.read(r)
	return tag
}

func (v *Validator) validate(value reflect.Value, fieldName string, tag *tag) *ValidationResult {
	switch value.Kind() {
	case reflect.Struct:
		ret, err := v.Validate(value.Interface())
		if err != nil {
			panic(err)
		}
		return ret
	case reflect.Map:
		ret := newValidationResult()
		if tag != nil && (tag.minProperties != nil || tag.maxProperties != nil) {
			l := int64(value.Len())
			if tag.minProperties != nil && l < *tag.minProperties {
				ret.add(&ValidationResult{
					Message: "minProperties",
					Name:    fieldName,
					Meta:    tag.meta,
				})
			}
			if tag.maxProperties != nil && l > *tag.maxProperties {
				ret.add(&ValidationResult{
					Message: "maxProperties",
					Name:    fieldName,
					Meta:    tag.meta,
				})
			}
		}
		if tag != nil && len(tag.required) > 0 {
			missing := []string{}
			for _, req := range tag.required {
				key := reflect.ValueOf(req)
				if !value.MapIndex(key).IsValid() {
					missing = append(missing, req)
				}
			}
			if len(missing) > 0 {
				ret.add(&ValidationResult{
					Message: "required",
					Name:    fieldName,
					Meta:    tag.meta,
				})
			}
		}
		for _, key := range value.MapKeys() {
			if tag != nil && tag.patternProperties != nil && !tag.patternProperties.MatchString(toString(key)) {
				ret.add(&ValidationResult{
					Message: "patternProperties",
					Name:    fieldName,
					Meta:    tag.meta,
				})
			}
			if key.Kind() == reflect.Ptr && !key.IsNil() {
				key = key.Elem()
			}
			if ve := v.validate(key, fmt.Sprintf("%s[%v](key)", fieldName, key.Interface()), nil); ve != nil && !ve.isEmpty() {
				ret.add(ve)
			}
			data := value.MapIndex(key)
			if data.Kind() == reflect.Ptr && !data.IsNil() {
				data = data.Elem()
			}
			if ve := v.validate(data, fmt.Sprintf("%s[%v](value)", fieldName, key.Interface()), nil); ve != nil && !ve.isEmpty() {
				ret.add(ve)
			}
		}
		return ret
	case reflect.Slice, reflect.Array:
		ret := newValidationResult()
		if tag != nil && (tag.minItems != nil || tag.maxItems != nil) {
			l := int64(value.Len())
			if tag.minItems != nil && l < *tag.minItems {
				ret.add(&ValidationResult{
					Message: "minItems",
					Name:    fieldName,
					Meta:    tag.meta,
				})
			}
			if tag.maxItems != nil && l > *tag.maxItems {
				ret.add(&ValidationResult{
					Message: "maxItems",
					Name:    fieldName,
					Meta:    tag.meta,
				})
			}
		}
		if tag != nil && tag.uniqueItems != nil && *tag.uniqueItems {
			for i := 1; i < value.Len(); i++ {
				for j := 0; j < i; j++ {
					if value.Index(i).Interface() == value.Index(j).Interface() {
						ret.add(&ValidationResult{
							Message: "uniqueItems",
							Name:    fieldName,
							Meta:    tag.meta,
						})
					}
				}
			}
		}
		for i := 0; i < value.Len(); i++ {
			if ve := v.validate(value.Index(i), fmt.Sprintf("%s[%d]", fieldName, i), tag); ve != nil && !ve.isEmpty() {
				ret.add(ve)
			}
		}
		return ret
	case reflect.String:
		str := value.String()
		if ve := v.validateString(str, fieldName, tag); !ve.isEmpty() {
			return ve
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		str := strconv.FormatInt(value.Int(), 10)
		num, _ := new(big.Float).SetString(str)
		if ve := v.validateNumber(num, fieldName, tag); !ve.isEmpty() {
			return ve
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		str := strconv.FormatUint(value.Uint(), 10)
		num, _ := new(big.Float).SetString(str)
		if ve := v.validateNumber(num, fieldName, tag); !ve.isEmpty() {
			return ve
		}
	case reflect.Float32, reflect.Float64:
		num := big.NewFloat(value.Float())
		if ve := v.validateNumber(num, fieldName, tag); !ve.isEmpty() {
			return ve
		}
	}
	return nil
}

func (v *Validator) validateString(str, fieldName string, tag *tag) *ValidationResult {
	ret := newValidationResult()
	if tag != nil && (tag.minLength != nil || tag.maxLength != nil) {
		l := int64(utf8.RuneCountInString(str))
		if tag.minLength != nil && l < *tag.minLength {
			ret.add(&ValidationResult{
				Message: "minLength",
				Name:    fieldName,
				Meta:    tag.meta,
			})
		}
		if tag.maxLength != nil && l > *tag.maxLength {
			ret.add(&ValidationResult{
				Message: "maxLength",
				Name:    fieldName,
				Meta:    tag.meta,
			})
		}
	}
	if tag != nil && tag.pattern != nil && !tag.pattern.MatchString(str) {
		ret.add(&ValidationResult{
			Message: "pattern",
			Name:    fieldName,
			Meta:    tag.meta,
		})
	}
	if tag != nil && tag.format != nil {
		if e := v.execFormat(*tag.format, str); e != nil {
			ret.add(&ValidationResult{
				Message: "format",
				Name:    fieldName,
				Meta:    tag.meta,
			})
		}
	}
	if tag != nil && len(tag.enum) > 0 && !contains(tag.enum, str) {
		ret.add(&ValidationResult{
			Message: "enum",
			Name:    fieldName,
			Meta:    tag.meta,
		})
	}
	return ret
}

func (v *Validator) validateNumber(num *big.Float, fieldName string, tag *tag) *ValidationResult {
	ret := newValidationResult()
	if tag != nil && tag.minimum != nil {
		if tag.exclusiveMinimumD4 != nil && *tag.exclusiveMinimumD4 && num.Cmp(tag.minimum) <= 0 {
			ret.add(&ValidationResult{
				Message: "minimum exclusiveMinimum",
				Name:    fieldName,
				Meta:    tag.meta,
			})
		}
		if num.Cmp(tag.minimum) < 0 {
			ret.add(&ValidationResult{
				Message: "minimum",
				Name:    fieldName,
				Meta:    tag.meta,
			})
		}
	}
	if tag != nil && tag.exclusiveMinimumD6 != nil && num.Cmp(tag.exclusiveMinimumD6) <= 0 {
		ret.add(&ValidationResult{
			Message: "exclusiveMinimum",
			Name:    fieldName,
			Meta:    tag.meta,
		})
	}
	if tag != nil && tag.maximum != nil {
		if tag.exclusiveMaximumD4 != nil && *tag.exclusiveMaximumD4 && num.Cmp(tag.maximum) >= 0 {
			ret.add(&ValidationResult{
				Message: "maximum exclusiveMaximum",
				Name:    fieldName,
				Meta:    tag.meta,
			})
		}
		if num.Cmp(tag.maximum) > 0 {
			ret.add(&ValidationResult{
				Message: "maximum",
				Name:    fieldName,
				Meta:    tag.meta,
			})
		}
	}
	if tag != nil && tag.exclusiveMaximumD6 != nil && num.Cmp(tag.exclusiveMaximumD6) >= 0 {
		ret.add(&ValidationResult{
			Message: "exclusiveMaximum",
			Name:    fieldName,
			Meta:    tag.meta,
		})
	}
	if tag != nil && tag.multipleOf != nil {
		if m := new(big.Float).Quo(num, tag.multipleOf); !m.IsInt() {
			ret.add(&ValidationResult{
				Message: "multipleOf",
				Name:    fieldName,
				Meta:    tag.meta,
			})
		}
	}
	if tag != nil && len(tag.enum) > 0 {
		if !contains(tag.enum, num.String()) {
			ret.add(&ValidationResult{
				Message: "enum",
				Name:    fieldName,
				Meta:    tag.meta,
			})
		}
	}
	return ret
}
