package main

import (
	"fmt"
	"log"
	"net/http"

	"rate-limiter/config"
	redis "rate-limiter/config/redis-config"
	genericratelimiter "rate-limiter/generic-rate-limiter"
)

func main() {
	if err := redis.InitClients(); err != nil {
		log.Fatalf("Failed to initialize Redis clients: %v", err)
	}
	mux := http.NewServeMux()

	mux.HandleFunc("/generic-rate-limiter", genericratelimiter.GenericRateLimitHandler)

	fmt.Printf("Server starting on PORT %s\n", config.GetEnv("PORT"))
	err := http.ListenAndServe(":"+config.GetEnv("PORT"), mux)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
