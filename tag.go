package jsonschema

import (
	"bytes"
	"math/big"
	"regexp"
	"strconv"
)

var (
	preMinimum        = []byte("mini")
	preMaximum        = []byte("maxi")
	preExclusive      = []byte("excl")
	exclusiveMaximum  = []byte("usiveMaximum")
	exclusiveMinimum  = []byte("usiveMinimum")
	preMultipleOf     = []byte("mult")
	preMinLength      = []byte("minL")
	preMaxLength      = []byte("maxL")
	prePattern        = []byte("patt")
	patternProperties = []byte("roperties:")
	preFormat         = []byte("form")
	preMinItems       = []byte("minI")
	preMaxItems       = []byte("maxI")
	preUniqueItems    = []byte("uniq")
	preMinProperties  = []byte("minP")
	preMaxProperties  = []byte("maxP")
	preRequired       = []byte("requ")
	preEnum           = []byte("enum")
)

func newTag(meta string) *tag {
	return &tag{
		meta: meta,
		enum: []string{},
	}
}

type tag struct {
	meta string
	// number validators
	minimum            *big.Float
	maximum            *big.Float
	exclusiveMinimumD4 *bool      // draft 4
	exclusiveMaximumD4 *bool      // draft 4
	exclusiveMinimumD6 *big.Float // draft 6
	exclusiveMaximumD6 *big.Float // draft 6
	multipleOf         *big.Float
	// string validators
	minLength *int64
	maxLength *int64
	pattern   *regexp.Regexp
	format    *string
	// array validations
	minItems    *int64
	maxItems    *int64
	uniqueItems *bool
	// object validations
	minProperties     *int64
	maxProperties     *int64
	patternProperties *regexp.Regexp
	required          []string
	// all validations
	enum []string
}

func (t *tag) read(r *reader) error {
	if r.IsEOF() {
		return nil
	}
	prefix := r.ReadBytes(4)
	switch {
	case bytes.Equal(prefix[:], preMinimum):
		r.SkipDelimiter()
		value := r.ReadSeparator()
		min, _ := strconv.ParseFloat(value, 64)
		t.minimum = big.NewFloat(min)
		if !r.IsEOF() {
			t.read(r)
		}
	case bytes.Equal(prefix[:], preMaximum):
		r.SkipDelimiter()
		value := r.ReadSeparator()
		max, _ := strconv.ParseFloat(value, 64)
		t.maximum = big.NewFloat(max)
		if !r.IsEOF() {
			t.read(r)
		}
	case bytes.Equal(prefix[:], preExclusive):
		prefix = r.ReadBytes(12)
		r.SkipDelimiter()
		value := r.ReadSeparator()
		if bytes.Equal(prefix[:], exclusiveMinimum) {
			if value == "true" || value == "false" {
				exclusive, _ := strconv.ParseBool(value)
				t.exclusiveMinimumD4 = &exclusive
			} else {
				exclusive, _ := strconv.ParseFloat(value, 64)
				t.exclusiveMinimumD6 = big.NewFloat(exclusive)
			}
		}
		if bytes.Equal(prefix[:], exclusiveMaximum) {
			if value == "true" || value == "false" {
				exclusive, _ := strconv.ParseBool(value)
				t.exclusiveMaximumD4 = &exclusive
			} else {
				exclusive, _ := strconv.ParseFloat(value, 64)
				t.exclusiveMaximumD6 = big.NewFloat(exclusive)
			}
		}
		if !r.IsEOF() {
			t.read(r)
		}
	case bytes.Equal(prefix[:], preMultipleOf):
		r.SkipDelimiter()
		value := r.ReadSeparator()
		num, _ := strconv.ParseFloat(value, 64)
		t.multipleOf = big.NewFloat(num)
		if !r.IsEOF() {
			t.read(r)
		}
	case bytes.Equal(prefix[:], preMinLength):
		r.SkipDelimiter()
		value := r.ReadSeparator()
		num, _ := strconv.ParseInt(value, 10, 64)
		t.minLength = &num
		if !r.IsEOF() {
			t.read(r)
		}
	case bytes.Equal(prefix[:], preMaxLength):
		r.SkipDelimiter()
		value := r.ReadSeparator()
		num, _ := strconv.ParseInt(value, 10, 64)
		t.maxLength = &num
		if !r.IsEOF() {
			t.read(r)
		}
	case bytes.Equal(prefix[:], prePattern):
		prefix = r.ReadBytes(3)
		// todo...ern check
		hasSep, _ := r.ReadByte()
		if hasSep == byte(':') || hasSep == byte('=') {
			t.pattern = regexp.MustCompile(r.ReadSeparator())
		} else {
			prefix = r.ReadBytes(10)
			if bytes.Equal(prefix[:], patternProperties) {
				t.patternProperties = regexp.MustCompile(r.ReadSeparator())
			}
		}
		if !r.IsEOF() {
			t.read(r)
		}
	case bytes.Equal(prefix[:], preFormat):
		r.SkipDelimiter()
		value := r.ReadSeparator()
		t.format = &value
		if !r.IsEOF() {
			t.read(r)
		}
	case bytes.Equal(prefix[:], preMinItems):
		r.SkipDelimiter()
		value := r.ReadSeparator()
		num, _ := strconv.ParseInt(value, 10, 64)
		t.minItems = &num
		if !r.IsEOF() {
			t.read(r)
		}
	case bytes.Equal(prefix[:], preMaxItems):
		r.SkipDelimiter()
		value := r.ReadSeparator()
		num, _ := strconv.ParseInt(value, 10, 64)
		t.maxItems = &num
		if !r.IsEOF() {
			t.read(r)
		}
	case bytes.Equal(prefix[:], preUniqueItems):
		r.SkipDelimiter()
		value := r.ReadSeparator()
		uniq, _ := strconv.ParseBool(value)
		t.uniqueItems = &uniq
		if !r.IsEOF() {
			t.read(r)
		}
	case bytes.Equal(prefix[:], preMinProperties):
		r.SkipDelimiter()
		value := r.ReadSeparator()
		num, _ := strconv.ParseInt(value, 10, 64)
		t.minProperties = &num
		if !r.IsEOF() {
			t.read(r)
		}
	case bytes.Equal(prefix[:], preMaxProperties):
		r.SkipDelimiter()
		value := r.ReadSeparator()
		num, _ := strconv.ParseInt(value, 10, 64)
		t.maxProperties = &num
		if !r.IsEOF() {
			t.read(r)
		}
	case bytes.Equal(prefix[:], preRequired):
		r.SkipDelimiter()
		buf := []byte{}
		for {
			b, _ := r.ReadByte()
			if b == byte('[') {
				continue
			}
			if b == byte(']') {
				t.required = append(t.required, string(buf))
				break
			}
			if b == byte(',') {
				t.required = append(t.required, string(buf))
				buf = []byte{}
				continue
			}

			buf = append(buf, b)
		}
		if !r.IsEOF() {
			t.read(r)
		}
	case bytes.Equal(prefix[:], preEnum):
		r.SkipDelimiter()
		buf := []byte{}
		for {
			b, _ := r.ReadByte()
			if b == byte('[') {
				continue
			}
			if b == byte(']') {
				t.enum = append(t.enum, string(buf))
				break
			}
			if b == byte(',') {
				t.enum = append(t.enum, string(buf))
				buf = []byte{}
				continue
			}

			buf = append(buf, b)
		}
		if !r.IsEOF() {
			t.read(r)
		}
	}
	if !r.IsEOF() {
		r.ReadByte()
		t.read(r)
	}
	return nil
}
