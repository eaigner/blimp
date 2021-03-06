package blimp

import (
	"io"
	"net/http"
)

const (
	TypeDisconnect = int(iota)
	TypeConnect
	TypeHeartbeat
	TypeMessage
	TypeJsonMessage
	TypeEvent
	TypeAck
	TypeError
	TypeNoop
)

// Conn represents a session based client connection
type Conn interface {
	Server() Server
	SessionId() string
	Transport() Transport
	SetTransport(trans Transport)
	Opened()
	Receive(m Message)
	// TODO: RemoteAddr() net.TCPAddr
}

// Handler represents a session handler
type Handler interface {
	Trigger() string
	Connected(conn Conn)
	Disconnected(conn Conn)
	Received(conn Conn, m Message)
}

type Message interface {
	io.Reader
	Type() int
	Id() int
	Ack() bool
	Endpoint() string
}

type Codec interface {
	Encode(m Message, w io.Writer) (n int64, err error)
	Decode(b []byte) (m Message, err error)
}

type Transport interface {
	New(conn Conn) Transport
	Conn() Conn
	Name() string
	ServeHTTP(rw http.ResponseWriter, req *http.Request)
	Send(m Message)
	Close() error
	UseHeartbeat() bool
}

type Server interface {
	http.Handler
	RegisterTransport(trans Transport)
	RegisterHandler(handler Handler)
	ProtocolVersion() string
	AuthHandler() AuthHandler
	GenerateSessionId() string
	Broadcast(m Message)
	BroadcastExcept(conn Conn, m Message)
}

type AuthHandler interface {
	Authorize(req *http.Request) bool
}
