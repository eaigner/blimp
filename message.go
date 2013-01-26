package blimp

import (
	"io"
)

type message struct {
	typ      int
	id       int
	ack      bool
	endpoint string
	reader   io.Reader
}

func newMessage(typ, id int, ack bool, endpoint string, r io.Reader) Message {
	return &message{
		typ:      typ,
		id:       id,
		ack:      ack,
		endpoint: endpoint,
		reader:   r,
	}
}

func (m *message) Type() int {
	return m.typ
}

func (m *message) Id() int {
	return m.id
}

func (m *message) Ack() bool {
	return m.ack
}

func (m *message) Endpoint() string {
	return m.endpoint
}

func (m *message) Read(p []byte) (n int, err error) {
	if m.reader != nil {
		return m.reader.Read(p)
	}
	return 0, nil
}
