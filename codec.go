package blimp

import (
	"bytes"
	"io"
	"strconv"
)

var StdCodec Codec = &codec{}

type codec struct{}

func (c *codec) Encode(m Message, w io.Writer) (int, error) {
	// Write head
	var b bytes.Buffer
	b.WriteString(strconv.Itoa(m.Type()))
	b.WriteByte(':')
	if m.Id() != 0 {
		b.WriteString(strconv.Itoa(m.Id()))
		if m.Ack() {
			b.WriteByte('+')
		}
	}
	b.WriteByte(':')
	b.WriteString(m.Endpoint())
	b.WriteByte(':')
	n, err := w.Write(b.Bytes())
	if err != nil {
		return n, err
	}

	// Write data
	n2, err := w.Write(m.Bytes())
	n += n2

	return n, err
}

func (c *codec) Decode(b []byte) (Message, error) {
	panic("TODO: implement")
}
