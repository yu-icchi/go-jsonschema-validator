package jsonschema

import (
	"math/big"
	"reflect"
	"strconv"
)

func contains(strs []string, str string) bool {
	for _, s := range strs {
		if s == str {
			return true
		}
	}
	return false
}

func valueIs(i byte, numbers ...byte) bool {
	for n := range numbers {
		if i == numbers[n] {
			return true
		}
	}
	return false
}

func toString(value reflect.Value) string {
	switch value.Kind() {
	case reflect.String:
		return value.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(value.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(value.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return big.NewFloat(value.Float()).String()
	case reflect.Bool:
	}
	return ""
}
