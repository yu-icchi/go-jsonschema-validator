package jsonschema

import (
	"bytes"
	"encoding/binary"
)

type reader struct {
	buf *bytes.Reader
	pos int
	max int
}

func newReader(b []byte) *reader {
	return &reader{buf: bytes.NewReader(b), pos: 0, max: len(b)}
}

func (r *reader) IsEOF() bool {
	return r.pos >= r.max
}

func (r *reader) ReadByte() (byte, error) {
	var value byte
	if err := binary.Read(r.buf, binary.BigEndian, &value); err != nil {
		return 0, err
	}
	r.pos++
	return value, nil
}

func (r *reader) ReadBytes(num int) []byte {
	value := make([]byte, num)
	if err := binary.Read(r.buf, binary.BigEndian, value); err != nil {
		panic(err)
	}
	r.pos += num
	return value
}

func (r *reader) ReadSeparator() string {
	buf := make([]byte, 0, 16)
	for {
		b, _ := r.ReadByte()
		if b == byte(',') {
			break
		}
		buf = append(buf, b)
		if r.IsEOF() {
			break
		}
	}
	return string(buf)
}

func (r *reader) SkipDelimiter() {
	for {
		b, _ := r.ReadByte()
		if !valueIs(b, byte(':'), byte('=')) {
			continue
		}
		break
	}
}
