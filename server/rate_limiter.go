package server

import (
	"sync"
	"time"
)

type RateLimiter struct {
	mu sync.Mutex
	limit int
	window time.Duration
	events []time.Time
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{limit: limit, window: window}
}

func (r *RateLimiter) Allow() bool {
	r.mu.Lock(); defer r.mu.Unlock()
	now := time.Now()
	cutoff := now.Add(-r.window)
	keep := r.events[:0]
	for _, t := range r.events { if t.After(cutoff) { keep = append(keep, t) } }
	r.events = keep
	if len(r.events) >= r.limit { return false }
	r.events = append(r.events, now)
	return true
}
