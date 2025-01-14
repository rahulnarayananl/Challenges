package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"6.Simple_Rate_Limiter/util"
)

type LeakyBucket struct {
	capacity int
	queue    []time.Time
	leakRate time.Duration
	mu       sync.Mutex
}

type LeakyBucketLimiter struct {
	buckets map[string]*LeakyBucket
	mu      sync.Mutex
}

func NewLeakyBucket(capacity int, leakRate time.Duration) *LeakyBucket {
	return &LeakyBucket{
		capacity: capacity,
		queue:    make([]time.Time, 0, capacity),
		leakRate: leakRate,
	}
}

func (b *LeakyBucket) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()

	for len(b.queue) > 0 && now.Sub(b.queue[0]) >= b.leakRate {
		b.queue = b.queue[1:] // Remove the oldest request
	}

	if len(b.queue) < b.capacity {
		b.queue = append(b.queue, now)
		return true
	}

	return false
}

func (l *LeakyBucketLimiter) GetBucket(clientIP string) *LeakyBucket {
	l.mu.Lock()
	defer l.mu.Unlock()

	bucket, exists := l.buckets[clientIP]
	if !exists {
		bucket = NewLeakyBucket(10, time.Second)
		l.buckets[clientIP] = bucket
	}

	return bucket
}

func (l *LeakyBucketLimiter) LeakyBucketMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientIP := util.ReadUserIP(r)
		bucket := l.GetBucket(clientIP)
		if !bucket.Allow() {
			w.WriteHeader(http.StatusTooManyRequests)
			fmt.Fprintf(w, "Too many requests\n")
			return
		}
		next(w, r)
	}
}
