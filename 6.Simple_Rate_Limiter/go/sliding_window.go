package main

import (
	"sync"
	"time"
)

type SlidingWindow struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

func NewSlidingWindowRateLimiter() *SlidingWindow {
	return &SlidingWindow{
		requests: make(map[string][]time.Time),
		limit:    10,
		window:   time.Minute,
	}
}

func (s *SlidingWindow) GetBucket(key string) *TokenBucket {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-s.window)
	requests := s.requests[key]

	var validRequests []time.Time
	for _, t := range requests {
		if t.After(windowStart) {
			validRequests = append(validRequests, t)
		}
	}
	s.requests[key] = validRequests

	if len(validRequests) < s.limit {
		s.requests[key] = append(s.requests[key], now)
		return &TokenBucket{tokens: 1} // Allow request
	}
	return &TokenBucket{tokens: 0} // Deny request
}

func (s *SlidingWindow) CleanupExpiredBuckets() {
	for {
		time.Sleep(time.Minute)
		s.mu.Lock()
		for key, requests := range s.requests {
			windowStart := time.Now().Add(-s.window)
			var validRequests []time.Time
			for _, t := range requests {
				if t.After(windowStart) {
					validRequests = append(validRequests, t)
				}
			}
			s.requests[key] = validRequests
		}
		s.mu.Unlock()
	}
}
