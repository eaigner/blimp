package blimp

type conn struct {
	server Server
	sid    string
	trans  Transport
}

func newConn(server Server, sessionId string) Conn {
	return &conn{server: server, sid: sessionId}
}

func (c *conn) Server() Server {
	return c.server
}

func (c *conn) SessionId() string {
	return c.sid
}

func (c *conn) Transport() Transport {
	return c.trans
}

func (c *conn) SetTransport(trans Transport) {
	c.trans = trans
}

func (c *conn) Opened() {
	// TODO: implement
}

func (c *conn) Receive(m Message) {
	// TODO: implement
}
