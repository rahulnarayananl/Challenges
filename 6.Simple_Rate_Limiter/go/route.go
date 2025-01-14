package main

import (
	"net/http"
	"time"
)

func Route(option string) *http.ServeMux {
	router := http.NewServeMux()
	switch option {
	case "tokenbucket":
		limiter := &TokenRateLimiter{
			buckets:         make(map[string]*TokenBucket),
			cleanupInterval: 1 * time.Minute,
		}
		router.Handle("/limited", limiter.TokenBucketMiddleware(http.HandlerFunc(Limited)))
	case "leakybucket":
		limiter := &LeakyBucketLimiter{
			buckets: make(map[string]*LeakyBucket),
		}
		router.Handle("/limited", limiter.LeakyBucketMiddleware(Limited))
	case "fixedwindow":
		limiter := &FixedWindowCounter{
			limit:       10,
			windowSize:  time.Minute,
			requests:    make(map[string]int),
			windowStart: time.Now(),
		}
		router.Handle("/limited", limiter.FixedWindowCounterMiddleware(Limited))
	// case "4":
	// 	router.HandleFunc("/limited", limiter.SlidingWindowLogMiddlewareRL(Limited))
	// case "5":
	// 	router.HandleFunc("/limited", limiter.SlidingWindowCounterMiddlewareRL(Limited))
	default:
		router.HandleFunc("/limited", Limited)
	}
	router.HandleFunc("/unlimited", Unlimited)
	return router
}
