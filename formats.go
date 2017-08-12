package jsonschema

import (
	"errors"
	"net"
	"net/mail"
	"net/url"
	"strings"
	"time"
)

func dataTime(data string) error {
	if _, err := time.Parse(time.RFC3339, data); err != nil {
		return err
	}
	if _, err := time.Parse(time.RFC3339Nano, data); err != nil {
		return err
	}
	return nil
}

func email(data string) error {
	if len(data) > 254 {
		return errors.New("")
	}
	at := strings.LastIndexByte(data, '@')
	if at == -1 {
		return errors.New("")
	}
	local := data[0:at]
	domain := data[at+1:]
	if len(local) > 64 {
		return errors.New("")
	}
	if err := hostname(domain); err != nil {
		return err
	}
	_, err := mail.ParseAddress(data)
	return err
}

func hostname(data string) error {
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

func ipv4(data string) error {
	if g := strings.Split(data, "."); len(g) != 4 {
		return errors.New("")
	}
	if net.ParseIP(data) == nil {
		return errors.New("")
	}
	return nil
}

func ipv6(data string) error {
	if !strings.Contains(data, ":") {
		return errors.New("")
	}
	if net.ParseIP(data) == nil {
		return errors.New("")
	}
	return nil
}

func uri(data string) error {
	u, err := url.Parse(data)
	if err != nil {
		return err
	}
	if !u.IsAbs() {
		return errors.New("")
	}
	return nil
}

func uriReference(data string) error {
	_, err := url.Parse(data)
	return err
}

func jsonPointer(data string) error {
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
