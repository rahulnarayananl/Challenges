package main

import (
	"fmt"
	"sync"
	"time"
)

type TokenBucket struct {
	capacity       int
	tokens         int
	lastRefillTime time.Time
	refillRate     int
	mu             sync.Mutex
}

func (tb *TokenBucket) AllowRequest() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefillTime).Seconds()
	refillTokens := int(elapsed) * tb.refillRate

	if refillTokens > 0 {
		tb.tokens = min(tb.capacity, tb.tokens+refillTokens)
		tb.lastRefillTime = now
	}

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	return false
}

type TokenRateLimiter struct {
	buckets         map[string]*TokenBucket // ip Vs TokenBuckets
	mu              sync.Mutex
	cleanupInterval time.Duration
}

func (rl *TokenRateLimiter) GetBucket(ip string) *TokenBucket {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if bucket, exists := rl.buckets[ip]; exists {
		return bucket
	}

	bucket := &TokenBucket{
		capacity:       10,
		tokens:         10,
		lastRefillTime: time.Now(),
		refillRate:     1,
	}
	rl.buckets[ip] = bucket
	return bucket
}

func (rl *TokenRateLimiter) CleanupExpiredBuckets() {
	for {
		time.Sleep(rl.cleanupInterval)

		rl.mu.Lock()
		for ip, bucket := range rl.buckets {
			// Cleanup logic: Remove buckets that have not been accessed for a while (e.g., 1 minute)
			if time.Since(bucket.lastRefillTime) > 1*time.Minute {
				fmt.Printf("Cleaning up expired bucket for IP: %s\n", ip)
				delete(rl.buckets, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func NewTokenBucketRateLimiter() *TokenRateLimiter {
	return &TokenRateLimiter{
		buckets:         make(map[string]*TokenBucket),
		cleanupInterval: 10 * time.Second,
	}
}
