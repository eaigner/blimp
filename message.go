package blimp

type message struct {
	data []byte
}

func newMessage(data []byte) Message {
	return &message{data: data}
}

func (m *message) Bytes() []byte {
	return m.data
}
