package server

import (
	"fmt"
	"sort"
	"sync"
)

// Hub owns the shared chat state. All mutations are protected by mu.
type Hub struct {
	mu      sync.RWMutex
	clients map[string]*Client
	rooms   map[string]map[string]*Client
}

func NewHub() *Hub {
	h := &Hub{clients: make(map[string]*Client), rooms: make(map[string]map[string]*Client)}
	h.rooms["general"] = make(map[string]*Client)
	return h
}

func (h *Hub) Register(c *Client) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, exists := h.clients[c.Name]; exists {
		return fmt.Errorf("username %q is already in use", c.Name)
	}
	h.clients[c.Name] = c
	if h.rooms[c.Room] == nil { h.rooms[c.Room] = make(map[string]*Client) }
	h.rooms[c.Room][c.Name] = c
	return nil
}

func (h *Hub) Remove(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, c.Name)
	if room := h.rooms[c.Room]; room != nil { delete(room, c.Name) }
}

func (h *Hub) Move(c *Client, room string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if old := h.rooms[c.Room]; old != nil { delete(old, c.Name) }
	if h.rooms[room] == nil { h.rooms[room] = make(map[string]*Client) }
	c.Room = room
	h.rooms[room][c.Name] = c
}

func (h *Hub) Broadcast(room, message string, except *Client) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, c := range h.rooms[room] {
		if c != except { c.Send(message) }
	}
}

func (h *Hub) Direct(name, message string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	c, ok := h.clients[name]
	if ok { c.Send(message) }
	return ok
}

func (h *Hub) Users(room string) []string {
	h.mu.RLock(); defer h.mu.RUnlock()
	out := make([]string, 0, len(h.rooms[room]))
	for name := range h.rooms[room] { out = append(out, name) }
	sort.Strings(out)
	return out
}

func (h *Hub) Rooms() []string {
	h.mu.RLock(); defer h.mu.RUnlock()
	out := make([]string, 0, len(h.rooms))
	for name := range h.rooms { out = append(out, name) }
	sort.Strings(out)
	return out
}
