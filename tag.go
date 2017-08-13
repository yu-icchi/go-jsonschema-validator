package jsonschema

import (
	"bytes"
	"errors"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

var (
	// ErrTagSyntax -
	ErrTagSyntax = errors.New("tag syntax error")
)

var (
	preMinimum        = []byte("mini")
	preMaximum        = []byte("maxi")
	mum               = []byte("mum")
	preExclusive      = []byte("excl")
	exclusiveMaximum  = []byte("usiveMaximum")
	exclusiveMinimum  = []byte("usiveMinimum")
	preMultipleOf     = []byte("mult")
	multipleOf        = []byte("ipleOf")
	preMinLength      = []byte("minL")
	preMaxLength      = []byte("maxL")
	length            = []byte("ength")
	prePattern        = []byte("patt")
	pattern           = []byte("ern")
	patternProperties = []byte("roperties:")
	preFormat         = []byte("form")
	format            = []byte("at")
	preMinItems       = []byte("minI")
	preMaxItems       = []byte("maxI")
	items             = []byte("tems")
	preUniqueItems    = []byte("uniq")
	uniqueItems       = []byte("ueItems")
	preMinProperties  = []byte("minP")
	preMaxProperties  = []byte("maxP")
	properties        = []byte("roperties")
	preRequired       = []byte("requ")
	required          = []byte("ired")
	preEnum           = []byte("enum")
)

func newTag() *tag {
	return &tag{
		required: []string{},
		enum:     []string{},
	}
}

type tag struct {
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
		prefix = r.ReadBytes(3)
		if !bytes.Equal(prefix[:], mum) {
			return ErrTagSyntax
		}
		r.SkipDelimiter()
		value := r.ReadSeparator()
		min, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return ErrTagSyntax
		}
		t.minimum = big.NewFloat(min)
		if !r.IsEOF() {
			t.read(r)
		}
	case bytes.Equal(prefix[:], preMaximum):
		prefix = r.ReadBytes(3)
		if !bytes.Equal(prefix[:], mum) {
			return ErrTagSyntax
		}
		r.SkipDelimiter()
		value := r.ReadSeparator()
		max, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return ErrTagSyntax
		}
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
				exclusive, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return ErrTagSyntax
				}
				t.exclusiveMinimumD6 = big.NewFloat(exclusive)
			}
		}
		if bytes.Equal(prefix[:], exclusiveMaximum) {
			if value == "true" || value == "false" {
				exclusive, _ := strconv.ParseBool(value)
				t.exclusiveMaximumD4 = &exclusive
			} else {
				exclusive, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return ErrTagSyntax
				}
				t.exclusiveMaximumD6 = big.NewFloat(exclusive)
			}
		}
		if !r.IsEOF() {
			t.read(r)
		}
	case bytes.Equal(prefix[:], preMultipleOf):
		prefix = r.ReadBytes(6)
		if !bytes.Equal(prefix[:], multipleOf) {
			return ErrTagSyntax
		}
		r.SkipDelimiter()
		value := r.ReadSeparator()
		num, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return ErrTagSyntax
		}
		t.multipleOf = big.NewFloat(num)
		if !r.IsEOF() {
			t.read(r)
		}
	case bytes.Equal(prefix[:], preMinLength):
		prefix = r.ReadBytes(5)
		if !bytes.Equal(prefix[:], length) {
			return ErrTagSyntax
		}
		r.SkipDelimiter()
		value := r.ReadSeparator()
		num, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return ErrTagSyntax
		}
		t.minLength = &num
		if !r.IsEOF() {
			t.read(r)
		}
	case bytes.Equal(prefix[:], preMaxLength):
		prefix = r.ReadBytes(5)
		if !bytes.Equal(prefix[:], length) {
			return ErrTagSyntax
		}
		r.SkipDelimiter()
		value := r.ReadSeparator()
		num, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return ErrTagSyntax
		}
		t.maxLength = &num
		if !r.IsEOF() {
			t.read(r)
		}
	case bytes.Equal(prefix[:], prePattern):
		prefix = r.ReadBytes(3)
		if !bytes.Equal(prefix[:], pattern) {
			return ErrTagSyntax
		}
		hasSep, err := r.ReadByte()
		if err != nil {
			return ErrTagSyntax
		}
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
		prefix = r.ReadBytes(2)
		if !bytes.Equal(prefix[:], format) {
			return ErrTagSyntax
		}
		r.SkipDelimiter()
		value := r.ReadSeparator()
		t.format = &value
		if !r.IsEOF() {
			t.read(r)
		}
	case bytes.Equal(prefix[:], preMinItems):
		prefix = r.ReadBytes(4)
		if !bytes.Equal(prefix[:], items) {
			return ErrTagSyntax
		}
		r.SkipDelimiter()
		value := r.ReadSeparator()
		num, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return ErrTagSyntax
		}
		t.minItems = &num
		if !r.IsEOF() {
			t.read(r)
		}
	case bytes.Equal(prefix[:], preMaxItems):
		prefix = r.ReadBytes(4)
		if !bytes.Equal(prefix[:], items) {
			return ErrTagSyntax
		}
		r.SkipDelimiter()
		value := r.ReadSeparator()
		num, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return ErrTagSyntax
		}
		t.maxItems = &num
		if !r.IsEOF() {
			t.read(r)
		}
	case bytes.Equal(prefix[:], preUniqueItems):
		prefix = r.ReadBytes(7)
		if !bytes.Equal(prefix[:], uniqueItems) {
			return ErrTagSyntax
		}
		r.SkipDelimiter()
		value := r.ReadSeparator()
		uniq, err := strconv.ParseBool(value)
		if err != nil {
			return ErrTagSyntax
		}
		t.uniqueItems = &uniq
		if !r.IsEOF() {
			t.read(r)
		}
	case bytes.Equal(prefix[:], preMinProperties):
		prefix = r.ReadBytes(9)
		if !bytes.Equal(prefix[:], properties) {
			return ErrTagSyntax
		}
		r.SkipDelimiter()
		value := r.ReadSeparator()
		num, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return ErrTagSyntax
		}
		t.minProperties = &num
		if !r.IsEOF() {
			t.read(r)
		}
	case bytes.Equal(prefix[:], preMaxProperties):
		prefix = r.ReadBytes(9)
		if !bytes.Equal(prefix[:], properties) {
			return ErrTagSyntax
		}
		r.SkipDelimiter()
		value := r.ReadSeparator()
		num, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return ErrTagSyntax
		}
		t.maxProperties = &num
		if !r.IsEOF() {
			t.read(r)
		}
	case bytes.Equal(prefix[:], preRequired):
		prefix = r.ReadBytes(4)
		if !bytes.Equal(prefix[:], required) {
			return ErrTagSyntax
		}
		r.SkipDelimiter()
		buf := []byte{}
		for {
			b, err := r.ReadByte()
			if err != nil {
				return ErrTagSyntax
			}
			if b == byte('[') {
				continue
			}
			if b == byte(']') {
				v := strings.TrimSpace(string(buf))
				t.required = append(t.required, v)
				break
			}
			if b == byte(',') {
				v := strings.TrimSpace(string(buf))
				t.required = append(t.required, v)
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
			b, err := r.ReadByte()
			if err != nil {
				return ErrTagSyntax
			}
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
