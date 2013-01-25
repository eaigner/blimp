package blimp

import (
	"code.google.com/p/go.net/websocket"
	"net/http"
)

func init() {
	StdServer.RegisterTransport(&websocketTransport{})
}

type websocketTransport struct {
	conn   Conn
	wsConn *websocket.Conn
	wsOpen bool
}

func (t *websocketTransport) New(conn Conn) Transport {
	return &websocketTransport{conn: conn}
}

func (t *websocketTransport) Conn() Conn {
	return t.conn
}

func (t *websocketTransport) Name() string {
	return "websocket"
}

func (t *websocketTransport) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	websocket.Handler(func(conn *websocket.Conn) {
		t.webSocketHandler(conn)
	}).ServeHTTP(rw, req)
}

func (t *websocketTransport) webSocketHandler(ws *websocket.Conn) {
	if t.wsOpen {
		t.Close()
	}
	t.wsConn = ws
	t.wsOpen = true
	t.conn.Opened()
	for {
		var b []byte
		err := websocket.Message.Receive(ws, &b)
		if err != nil {
			t.Close()
			return
		}
		// TODO: decode frame and packet type, then pass to receive
		// t.conn.Receive(newMessage(b))
	}
}

func (t *websocketTransport) Send(m Message) {
	websocket.Message.Send(t.wsConn, m)
}

func (t *websocketTransport) Close() error {
	t.wsOpen = false
	return t.wsConn.Close()
}

func (t *websocketTransport) UseHeartbeat() bool {
	return true
}
