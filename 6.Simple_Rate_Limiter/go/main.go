package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"6.Simple_Rate_Limiter/util"
)

type RateLimiter interface {
	GetBucket(string) *TokenBucket
	CleanupExpiredBuckets()
}

var (
	algorithm = flag.String("b", "tokenbucket", "Rate limiting algorithm")
)

var rateLimiter RateLimiter

func main() {
	if *algorithm == "tokenbucket" {
		rateLimiter = NewTokenBucketRateLimiter()
		go rateLimiter.CleanupExpiredBuckets()

		fmt.Println("Server started")

		http.HandleFunc("/unlimited", Unlimited)
		http.HandleFunc("/limited", Limited)
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("Nope bye")
	}
}

func Unlimited(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Unlimited! Let's Go!")
}

func Limited(w http.ResponseWriter, r *http.Request) {
	ip := util.ReadUserIP(r)
	tb := rateLimiter.GetBucket(ip)
	if tb.AllowRequest() {
		fmt.Fprintf(w, "Cosmos Allows you ! \n")
	} else {
		// w.WriteHeader(http.StatusTooManyRequests)
		http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
	}
}
