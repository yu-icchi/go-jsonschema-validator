package jsonschema

import (
	"errors"
	"net"
	"net/mail"
	"net/url"
	"reflect"
	"strings"
	"time"
)

func dateTime(value *reflect.Value, field *reflect.StructField) error {
	if value.Kind() != reflect.String {
		return errors.New("format/dateTime: invalid value kind")
	}
	data := value.String()
	if _, err := time.Parse(time.RFC3339, data); err != nil {
		return err
	}
	if _, err := time.Parse(time.RFC3339Nano, data); err != nil {
		return err
	}
	return nil
}

func email(value *reflect.Value, field *reflect.StructField) error {
	if value.Kind() != reflect.String {
		return errors.New("format/email: invalid value kind")
	}
	data := value.String()
	if len(data) > 254 {
		return errors.New("")
	}
	at := strings.LastIndexByte(data, '@')
	if at == -1 {
		return errors.New("")
	}
	local := data[0:at]
	if len(local) > 64 {
		return errors.New("")
	}
	domain := data[at+1:]
	domain = strings.TrimSuffix(domain, ".")
	if len(domain) > 253 {
		return errors.New("")
	}
	for _, label := range strings.Split(domain, ".") {
		if l := len(label); l < 1 || l > 63 {
			return errors.New("")
		}
		if f := domain[0]; f >= '0' && f <= '9' || f == '-' {
			return errors.New("")
		}
		if label[len(label)-1] == '-' {
			return errors.New("")
		}
		for _, c := range label {
			if valid := c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '-'; !valid {
				return errors.New("")
			}
		}
	}
	_, err := mail.ParseAddress(data)
	return err
}

func hostname(value *reflect.Value, field *reflect.StructField) error {
	if value.Kind() != reflect.String {
		return errors.New("format/hostname: invalid value kind")
	}
	data := value.String()
	data = strings.TrimSuffix(data, ".")
	if len(data) > 253 {
		return errors.New("")
	}
	for _, label := range strings.Split(data, ".") {
		if l := len(label); l < 1 || l > 63 {
			return errors.New("")
		}
		if f := data[0]; f >= '0' && f <= '9' || f == '-' {
			return errors.New("")
		}
		if label[len(label)-1] == '-' {
			return errors.New("")
		}
		for _, c := range label {
			if valid := c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '-'; !valid {
				return errors.New("")
			}
		}
	}
	return nil
}

func ipv4(value *reflect.Value, field *reflect.StructField) error {
	if value.Kind() != reflect.String {
		return errors.New("format/ipv4: invalid value kind")
	}
	data := value.String()
	if g := strings.Split(data, "."); len(g) != 4 {
		return errors.New("")
	}
	if net.ParseIP(data) == nil {
		return errors.New("")
	}
	return nil
}

func ipv6(value *reflect.Value, field *reflect.StructField) error {
	if value.Kind() != reflect.String {
		return errors.New("format/ipv6: invalid value kind")
	}
	data := value.String()
	if !strings.Contains(data, ":") {
		return errors.New("")
	}
	if net.ParseIP(data) == nil {
		return errors.New("")
	}
	return nil
}

func uri(value *reflect.Value, field *reflect.StructField) error {
	if value.Kind() != reflect.String {
		return errors.New("format/url: invalid value kind")
	}
	data := value.String()
	u, err := url.Parse(data)
	if err != nil {
		return err
	}
	if !u.IsAbs() {
		return errors.New("")
	}
	return nil
}

func uriReference(value *reflect.Value, field *reflect.StructField) error {
	if value.Kind() != reflect.String {
		return errors.New("format/uriReference: invalid value kind")
	}
	data := value.String()
	_, err := url.Parse(data)
	return err
}

func jsonPointer(value *reflect.Value, field *reflect.StructField) error {
	if value.Kind() != reflect.String {
		return errors.New("format/jsonPointer: invalid value kind")
	}
	data := value.String()
	for _, item := range strings.Split(data, "/") {
		for i := 0; i < len(item); i++ {
			if item[i] == '~' {
				if i == len(item)-1 {
					return errors.New("")
				}
				switch item[i+1] {
				case '~', '0', '1':
					// valid
				default:
					return errors.New("")
				}
			}
		}
	}
	return nil
}
