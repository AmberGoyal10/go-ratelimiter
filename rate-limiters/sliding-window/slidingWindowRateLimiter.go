package slidingwindow

import (
	"context"
	"errors"
	"log"
	"os"
	redisconfig "rate-limiter/config/redis-config"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func SlidingWindowRateLimiter(key string, capacity int, ttlSeconds int) (bool, int, error) {
	redisClient := redisconfig.GetClient(redisconfig.RATE_LIMITER_DB)

	timestamp := time.Now().Unix()
	script, err := os.ReadFile("rate-limiters/sliding-window/sliding_window_rate_limiter.lua")
	id := uuid.New().String()

	if err != nil {
		log.Fatalf("Failed to read concurrent_rate_limiter.lua: %v", err)
	}

	redisClient.ZRemRangeByScore(context.Background(), key, "0", strconv.FormatInt(timestamp-int64(ttlSeconds), 10))

	keys := []string{key}
	args := []string{strconv.FormatInt(int64(capacity), 10), strconv.FormatInt(timestamp, 10), id, strconv.FormatInt(int64(ttlSeconds), 10)}
	result := redisClient.Eval(context.Background(), string(script), keys, args)

	response := result.Val()
	responseErr := result.Err()

	// Handle errors
	if responseErr != nil {
		err := errors.New("error executing Redis Eval")
		log.Println(err)
		return false, 0, err
	}

	// Check if response is valid
	if response == nil {
		err := errors.New("no response from Redis")
		log.Println("No response from Redis")
		return false, 0, err
	}

	// Type assert the response
	responseArray, ok := response.([]interface{})
	if !ok {
		err := errors.New("invalid response format")
		log.Printf("Invalid response format: %v", response)
		return false, 0, err
	}

	if len(responseArray) < 2 {
		log.Printf("Insufficient response data: %v", responseArray)
		err := errors.New("insufficient response data")
		return false, 0, err
	}

	// Extract allowed and count values
	allowed := responseArray[0]
	count := responseArray[1]

	// Check if request is allowed
	if allowed == nil || allowed == false {
		return false, 0, nil
	}

	// Convert count to int
	countInt, ok := count.(int64)
	if !ok {
		log.Printf("Invalid count format: %v", count)
		err := errors.New("invalid count format")
		return false, 0, err
	}

	return true, int(countInt), nil
}
