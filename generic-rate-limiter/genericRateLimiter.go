package genericratelimiter

import (
	"encoding/json"
	"net/http"
	"rate-limiter/helper"
	rateconfig "rate-limiter/rate-limiters"
	slidingwindow "rate-limiter/rate-limiters/sliding-window"
)

func GenericRateLimitHandler(w http.ResponseWriter, r *http.Request) {
	ip, err := helper.GetIpAddressFromHeaders(r.Header)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	key := "test" + "-" + ip

	capacity := rateconfig.ConCurrentRateLimitConfig.Capacity
	ttl := rateconfig.ConCurrentRateLimitConfig.TTL

	allowed, count, err := slidingwindow.SlidingWindowRateLimiter(key, capacity, ttl)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !allowed {
		response := rateconfig.RateLimitResponse{
			Success:      false,
			Message:      "Rate limit exceeded",
			CurrentCount: count,
			Remaining:    0,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Success response
	response := rateconfig.RateLimitResponse{
		Success:      true,
		Message:      "Request allowed",
		CurrentCount: count + 1,
		Remaining:    capacity - count - 1,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
