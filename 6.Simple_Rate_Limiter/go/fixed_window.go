package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"6.Simple_Rate_Limiter/util"
)

type FixedWindowCounter struct {
	limit       int
	windowSize  time.Duration
	requests    map[string]int
	windowStart time.Time
	mu          sync.Mutex
}

func (f *FixedWindowCounter) Allow(clientIP string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()

	now := time.Now()
	if now.Sub(f.windowStart) >= f.windowSize {
		f.requests = make(map[string]int)
		f.windowStart = now
	}

	f.requests[clientIP]++
	if f.requests[clientIP] > f.limit {
		return false
	}

	return true
}

func (f *FixedWindowCounter) FixedWindowCounterMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientIP := util.ReadUserIP(r)
		if !f.Allow(clientIP) {
			w.WriteHeader(http.StatusTooManyRequests)
			fmt.Fprintf(w, "Too many requests\n")
			return
		}

		next(w, r)
	}
}
