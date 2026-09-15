package server

import (
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	r := NewRateLimiter(2, time.Hour)
	if !r.Allow() || !r.Allow() { t.Fatal("first two events should be allowed") }
	if r.Allow() { t.Fatal("third event should be rejected") }
}

func TestHubRooms(t *testing.T) {
	h := NewHub()
	a := NewClient("alice", nil)
	b := NewClient("bob", nil)
	if err := h.Register(a); err != nil { t.Fatal(err) }
	if err := h.Register(b); err != nil { t.Fatal(err) }
	if got := len(h.Users("general")); got != 2 { t.Fatalf("users=%d want 2", got) }
	h.Move(b, "backend")
	if got := len(h.Users("backend")); got != 1 { t.Fatalf("backend users=%d want 1", got) }
}
