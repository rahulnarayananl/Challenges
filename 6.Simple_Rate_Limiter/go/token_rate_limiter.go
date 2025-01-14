package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"6.Simple_Rate_Limiter/util"
)

type TokenBucket struct {
	capacity       int
	tokens         int
	lastRefillTime time.Time
	refillRate     int
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
	log.Println("Get Bucket")
	return bucket
}

func (rl *TokenRateLimiter) CleanupExpiredBuckets() {
	for {
		time.Sleep(rl.cleanupInterval)

		rl.mu.Lock()
		for ip, bucket := range rl.buckets {
			if time.Since(bucket.lastRefillTime) > 1*time.Minute {
				fmt.Printf("Cleaning up expired bucket for IP: %s\n", ip)
				delete(rl.buckets, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *TokenRateLimiter) TokenBucketMiddleware(next http.Handler) http.HandlerFunc {
	go rl.CleanupExpiredBuckets()
	return func(w http.ResponseWriter, r *http.Request) {
		ip := util.ReadUserIP(r)
		bucket := rl.GetBucket(ip)

		rl.mu.Lock()
		now := time.Now()
		elapsed := now.Sub(bucket.lastRefillTime).Seconds()
		bucket.tokens += int(elapsed) * bucket.refillRate
		if bucket.tokens > bucket.capacity {
			bucket.tokens = bucket.capacity
		}
		bucket.lastRefillTime = now

		if bucket.tokens > 0 {
			bucket.tokens--
			rl.mu.Unlock()
			next.ServeHTTP(w, r)
		} else {
			rl.mu.Unlock()
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
		}
	}
}
