package blimp

import (
	"fmt"
	"net/http"
	"regexp"
)

var StdServer Server = newServer()

type server struct {
	transports map[string]Transport
	handler    map[string]Handler
	conns      map[string]Conn
}

func newServer() Server {
	return &server{
		transports: map[string]Transport{},
		conns:      map[string]Conn{},
	}
}

func (s *server) RegisterTransport(trans Transport) {
	s.transports[trans.Name()] = trans
}

func (s *server) RegisterHandler(handler Handler) {
	s.handler[handler.Trigger()] = handler
}

func (s *server) ProtocolVersion() string {
	return "1"
}

func (s *server) AuthHandler() AuthHandler {
	return nil // TODO: implement
}

func (s *server) GenerateSessionId() string {
	return AlphanumericId(16)
}

func (s *server) Broadcast(m Message) {
	for _, c := range s.conns {
		c.Transport().Send(m)
	}
}

func (s *server) BroadcastExcept(conn Conn, m Message) {
	for _, c := range s.conns {
		if c != conn {
			c.Transport().Send(m)
		}
	}
}

func (s *server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	rxFmt := fmt.Sprintf(`^(.+?)/(%s)(?:/([^/]+)/([^/]+))?/?$`, regexp.QuoteMeta(s.ProtocolVersion()))
	rx := regexp.MustCompile(rxFmt)
	c := rx.FindStringSubmatch(req.URL.Path)
	if len(c) == 0 {
		http.Error(rw, "invalid socket.io url", 404)
		return
	}
	transportId := c[3]
	if len(transportId) == 0 {
		s.handshake(rw, req)
		return
	}
	transport, ok := s.transports[transportId]
	if !ok {
		http.Error(rw, "transport not supported", 400)
		return
	}
	sessionId := c[4]
	conn, ok := s.conns[sessionId]
	if !ok {
		http.Error(rw, "invalid session id", 400)
	}
	if conn.Transport() == nil {
		conn.SetTransport(transport.New(conn))
	}
	conn.Transport().ServeHTTP(rw, req)
}

func (s *server) handshake(rw http.ResponseWriter, req *http.Request) {
	if ah := s.AuthHandler(); ah != nil {
		if !ah.Authorize(req) {
			http.Error(rw, "not authorized", 401)
			return
		}
	}
	h := rw.Header()
	h.Set("Access-Control-Allow-Origin", req.Header.Get("origin"))
	h.Set("Access-Control-Allow-Methods", "GET")
	h.Set("Access-Control-Allow-Credentials", "true")
	sessionId := s.GenerateSessionId()
	if len(sessionId) == 0 {
		http.Error(rw, "could not generate session id", 503)
		return
	}
	s.conns[sessionId] = newConn(s, sessionId)
	// TODO: notify handlers of connect event
}
