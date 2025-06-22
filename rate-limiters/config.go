package concurrentratelimiter

type RateLimitResponse struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	CurrentCount int    `json:"current_count"`
	Remaining    int    `json:"remaining"`
}

type RateLimitConfig struct {
	Capacity int `json:"capacity"`
	TTL      int `json:"ttl"`
}

var ConCurrentRateLimitConfig = RateLimitConfig{
	Capacity: 10,
	TTL:      100,
}
