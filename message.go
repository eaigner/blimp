package blimp

type message struct {
	typ      int
	id       int
	ack      bool
	endpoint string
	data     []byte
}

func newMessage(typ, id int, ack bool, endpoint string, data []byte) Message {
	return &message{
		typ:      typ,
		id:       id,
		ack:      ack,
		endpoint: endpoint,
		data:     data,
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

func (m *message) Bytes() []byte {
	return m.data
}
