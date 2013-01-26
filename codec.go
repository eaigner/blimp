package blimp

import (
	"bytes"
	"errors"
	"io"
	"strconv"
)

var StdCodec Codec = &codec{}

type codec struct{}

func (c *codec) Encode(m Message, w io.Writer) (int64, error) {
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
		return int64(n), err
	}

	// Write data
	n2, err := io.Copy(w, m)
	n2 += int64(n)

	return n2, err
}

func (c *codec) Decode(b []byte) (Message, error) {
	parts := bytes.SplitN(b, []byte(":"), 4)
	if len(parts) < 3 {
		return nil, errors.New("invalid packet")
	}
	mtype, err := strconv.Atoi(string(parts[0]))
	if err != nil {
		return nil, err
	}
	idPart := parts[1]
	ack := bytes.HasSuffix(idPart, []byte("+"))
	if ack {
		idPart = idPart[:len(idPart)-1]
	}
	pid, err := strconv.Atoi(string(idPart))
	if err != nil {
		return nil, err
	}
	endpoint := string(parts[2])
	var r io.Reader
	if len(parts) == 4 {
		r = bytes.NewBuffer(parts[3])
	}
	return newMessage(mtype, pid, ack, endpoint, r), nil
}
