package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

type Client struct {
	Name string
	Room string
	conn net.Conn
	mu sync.Mutex
	limiter *RateLimiter
}

func NewClient(name string, conn net.Conn) *Client {
	return &Client{Name: name, Room: "general", conn: conn, limiter: NewRateLimiter(6, 10*time.Second)}
}

func (c *Client) Send(message string) {
	c.mu.Lock(); defer c.mu.Unlock()
	if c.conn != nil { _, _ = fmt.Fprintln(c.conn, message) }
}

func (c *Client) Serve(h *Hub) {
	defer c.conn.Close()
	if err := h.Register(c); err != nil { c.Send("error: " + err.Error()); return }
	defer h.Remove(c)
	c.Send("Welcome, " + c.Name + ". Type /help for commands.")
	h.Broadcast(c.Room, "* "+c.Name+" joined #"+c.Room, c)
	defer h.Broadcast(c.Room, "* "+c.Name+" left #"+c.Room, c)

	s := bufio.NewScanner(c.conn)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" { continue }
		if !c.limiter.Allow() { c.Send("error: rate limit exceeded"); continue }
		if strings.HasPrefix(line, "/") {
			if !handleCommand(h, c, line) { return }
			continue
		}
		h.Broadcast(c.Room, fmt.Sprintf("[%s] %s: %s", c.Room, c.Name, line), nil)
	}
}
