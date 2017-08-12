package jsonschema

import (
	"bytes"
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

type ValidationError struct {
	Name    string
	Message string
	Causes  []*ValidationError
}

func (v *ValidationError) Error() string {
	buf := bytes.NewBufferString(v.Message)
	write(buf, v)
	return buf.String()
}

func write(buf *bytes.Buffer, e *ValidationError) {
	for _, err := range e.Causes {
		if err.Message != "" {
			buf.WriteString(err.Message)
		}
		if len(err.Causes) > 0 {
			write(buf, err)
		}
	}
}

func (v *ValidationError) isEmpty() bool {
	return v.Name == "" && v.Message == "" && len(v.Causes) == 0
}

func (v *ValidationError) add(err *ValidationError) {
	if err != nil {
		v.Causes = append(v.Causes, err)
	}
}

func newValidationError() *ValidationError {
	return &ValidationError{
		Causes: []*ValidationError{},
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
func (v *Validator) Validate(data interface{}) error {
	rt, rv := reflect.TypeOf(data), reflect.ValueOf(data)

	if rv.Kind() == reflect.Ptr && !rv.IsNil() {
		return v.Validate(rv.Elem().Interface())
	}
	if rv.Kind() != reflect.Struct && rv.Kind() != reflect.Interface {
		return errors.New("")
	}

	result := newValidationError()
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

		err := v.validate(value, name, tag)
		ret, ok := err.(*ValidationError)
		if ok && ret != nil && !ret.isEmpty() {
			result.add(ret)
		}
	}

	return result
}

func (v *Validator) parseTag(meta string) *tag {
	tag := newTag(meta)
	r := newReader([]byte(meta))
	tag.read(r)
	return tag
}

func (v *Validator) validate(value reflect.Value, fieldName string, tag *tag) error {
	switch value.Kind() {
	case reflect.Struct:
		return v.Validate(value.Interface())
	case reflect.Map:
		result := newValidationError()
		if tag != nil && (tag.minProperties != nil || tag.maxProperties != nil) {
			l := int64(value.Len())
			if tag.minProperties != nil && l < *tag.minProperties {
				result.add(&ValidationError{
					Message: fmt.Sprintf("minProperties:%d", *tag.minProperties),
					Name:    fieldName,
				})
			}
			if tag.maxProperties != nil && l > *tag.maxProperties {
				result.add(&ValidationError{
					Message: fmt.Sprintf("maxProperties:%d", *tag.maxProperties),
					Name:    fieldName,
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
				result.add(&ValidationError{
					Message: "required",
					Name:    fieldName,
				})
			}
		}
		for _, key := range value.MapKeys() {
			if tag != nil && tag.patternProperties != nil && !tag.patternProperties.MatchString(toString(key)) {
				result.add(&ValidationError{
					Message: fmt.Sprintf("patternProperties:%s", tag.patternProperties.String()),
					Name:    fieldName,
				})
			}
			if key.Kind() == reflect.Ptr && !key.IsNil() {
				key = key.Elem()
			}
			err := v.validate(key, fmt.Sprintf("%s[%v](key)", fieldName, key.Interface()), nil)
			ret, ok := err.(*ValidationError)
			if ok && ret != nil && !ret.isEmpty() {
				result.add(ret)
			}

			data := value.MapIndex(key)
			if data.Kind() == reflect.Ptr && !data.IsNil() {
				data = data.Elem()
			}
			err = v.validate(data, fmt.Sprintf("%s[%v](value)", fieldName, key.Interface()), nil)
			ret, ok = err.(*ValidationError)
			if ok && ret != nil && !ret.isEmpty() {
				ret.add(ret)
			}
		}
		return result
	case reflect.Slice, reflect.Array:
		result := newValidationError()
		if tag != nil && (tag.minItems != nil || tag.maxItems != nil) {
			l := int64(value.Len())
			if tag.minItems != nil && l < *tag.minItems {
				result.add(&ValidationError{
					Message: fmt.Sprintf("minItems:%d", *tag.minItems),
					Name:    fieldName,
				})
			}
			if tag.maxItems != nil && l > *tag.maxItems {
				result.add(&ValidationError{
					Message: fmt.Sprintf("maxItems:%d", *tag.maxItems),
					Name:    fieldName,
				})
			}
		}
		if tag != nil && tag.uniqueItems != nil && *tag.uniqueItems {
			for i := 1; i < value.Len(); i++ {
				for j := 0; j < i; j++ {
					if value.Index(i).Interface() == value.Index(j).Interface() {
						result.add(&ValidationError{
							Message: "uniqueItems:true",
							Name:    fieldName,
						})
					}
				}
			}
		}
		for i := 0; i < value.Len(); i++ {
			err := v.validate(value.Index(i), fmt.Sprintf("%s[%d]", fieldName, i), tag)
			ret, ok := err.(*ValidationError)
			if ok && ret != nil && !ret.isEmpty() {
				result.add(ret)
			}
		}
		return result
	case reflect.String:
		str := value.String()
		if ret := v.validateString(str, fieldName, tag); !ret.isEmpty() {
			return ret
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		str := strconv.FormatInt(value.Int(), 10)
		num, flag := new(big.Float).SetString(str)
		if !flag {
			return errors.New("invalid integer")
		}
		if ret := v.validateNumber(num, fieldName, tag); !ret.isEmpty() {
			return ret
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		str := strconv.FormatUint(value.Uint(), 10)
		num, flag := new(big.Float).SetString(str)
		if !flag {
			return errors.New("invalid unsigned integer")
		}
		if ret := v.validateNumber(num, fieldName, tag); !ret.isEmpty() {
			return ret
		}
	case reflect.Float32, reflect.Float64:
		num := big.NewFloat(value.Float())
		if ret := v.validateNumber(num, fieldName, tag); !ret.isEmpty() {
			return ret
		}
	}
	return nil
}

func (v *Validator) validateString(str, fieldName string, tag *tag) *ValidationError {
	ret := newValidationError()
	if tag != nil && (tag.minLength != nil || tag.maxLength != nil) {
		l := int64(utf8.RuneCountInString(str))
		if tag.minLength != nil && l < *tag.minLength {
			ret.add(&ValidationError{
				Message: fmt.Sprintf("minLength:%d(%d)", *tag.minLength, l),
				Name:    fieldName,
			})
		}
		if tag.maxLength != nil && l > *tag.maxLength {
			ret.add(&ValidationError{
				Message: fmt.Sprintf("maxLength:%d(%d)", *tag.maxLength, l),
				Name:    fieldName,
			})
		}
	}
	if tag != nil && tag.pattern != nil && !tag.pattern.MatchString(str) {
		ret.add(&ValidationError{
			Message: fmt.Sprintf("pattern:%s(%s)", tag.pattern.String(), str),
			Name:    fieldName,
		})
	}
	if tag != nil && tag.format != nil {
		if e := v.execFormat(*tag.format, str); e != nil {
			ret.add(&ValidationError{
				Message: e.Error(),
				Name:    fieldName,
			})
		}
	}
	if tag != nil && len(tag.enum) > 0 && !contains(tag.enum, str) {
		ret.add(&ValidationError{
			Message: "enum",
			Name:    fieldName,
		})
	}
	return ret
}

func (v *Validator) validateNumber(num *big.Float, fieldName string, tag *tag) *ValidationError {
	ret := newValidationError()
	if tag != nil && tag.minimum != nil {
		if tag.exclusiveMinimumD4 != nil && *tag.exclusiveMinimumD4 && num.Cmp(tag.minimum) <= 0 {
			ret.add(&ValidationError{
				Message: fmt.Sprintf("exclusiveMinimum:%s", tag.minimum.String()),
				Name:    fieldName,
			})
		}
		if num.Cmp(tag.minimum) < 0 {
			ret.add(&ValidationError{
				Message: fmt.Sprintf("minimum:%s", tag.minimum.String()),
				Name:    fieldName,
			})
		}
	}
	if tag != nil && tag.exclusiveMinimumD6 != nil && num.Cmp(tag.exclusiveMinimumD6) <= 0 {
		ret.add(&ValidationError{
			Message: fmt.Sprintf("exclusiveMinimum:%s", tag.exclusiveMinimumD6.String()),
			Name:    fieldName,
		})
	}
	if tag != nil && tag.maximum != nil {
		if tag.exclusiveMaximumD4 != nil && *tag.exclusiveMaximumD4 && num.Cmp(tag.maximum) >= 0 {
			ret.add(&ValidationError{
				Message: fmt.Sprintf("exclusiveMaximum:%s", tag.maximum.String()),
				Name:    fieldName,
			})
		}
		if num.Cmp(tag.maximum) > 0 {
			ret.add(&ValidationError{
				Message: fmt.Sprintf("maximum:%s", tag.maximum.String()),
				Name:    fieldName,
			})
		}
	}
	if tag != nil && tag.exclusiveMaximumD6 != nil && num.Cmp(tag.exclusiveMaximumD6) >= 0 {
		ret.add(&ValidationError{
			Message: fmt.Sprintf("exclusiveMaximum:%s", tag.exclusiveMaximumD6.String()),
			Name:    fieldName,
		})
	}
	if tag != nil && tag.multipleOf != nil {
		if m := new(big.Float).Quo(num, tag.multipleOf); !m.IsInt() {
			ret.add(&ValidationError{
				Message: fmt.Sprintf("multipleOf:%s", tag.multipleOf.String()),
				Name:    fieldName,
			})
		}
	}
	if tag != nil && len(tag.enum) > 0 {
		if !contains(tag.enum, num.String()) {
			ret.add(&ValidationError{
				Message: "enum",
				Name:    fieldName,
			})
		}
	}
	return ret
}
