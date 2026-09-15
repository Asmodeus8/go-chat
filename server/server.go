package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type Server struct { Hub *Hub }

func New() *Server { return &Server{Hub: NewHub()} }

func (s *Server) Serve(listener net.Listener) error {
	for {
		conn, err := listener.Accept()
		if err != nil { return err }
		go s.accept(conn)
	}
}

func (s *Server) accept(conn net.Conn) {
	_, _ = fmt.Fprintln(conn, "Adewale Go Chat")
	_, _ = fmt.Fprint(conn, "username: ")
	r := bufio.NewReader(conn)
	name, err := r.ReadString('\n')
	if err != nil { _ = conn.Close(); return }
	name = strings.TrimSpace(name)
	if !validName(name) { _, _ = fmt.Fprintln(conn, "invalid username"); _ = conn.Close(); return }
	NewClient(name, &bufferedConn{Conn: conn, r: r}).Serve(s.Hub)
}

func validName(s string) bool {
	if len(s) < 2 || len(s) > 24 { return false }
	for _, r := range s {
		if !(r == '_' || r == '-' || r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9') { return false }
	}
	return true
}

type bufferedConn struct { net.Conn; r *bufio.Reader }
func (c *bufferedConn) Read(p []byte) (int, error) { return c.r.Read(p) }
