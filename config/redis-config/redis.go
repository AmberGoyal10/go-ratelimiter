package redisconfig

import (
	"context"
	"fmt"
	"rate-limiter/config"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClients map[int]*redis.Client = make(map[int]*redis.Client)

func GetClient(db int) *redis.Client {
	if RedisClients[db] != nil {
		return RedisClients[db]
	}

	client := redis.NewClient(&redis.Options{
		Addr: config.GetEnv("REDIS_URL"),
		DB:   db,
	})

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		fmt.Printf("Failed to connect to Redis DB %d: %v\n", db, err)
		return nil
	}

	fmt.Println("Redis client created and connected successfully for db", db)
	RedisClients[db] = client
	return client
}

func InitClients() error {
	configClient := GetClient(CONFIG_REDIS_DB)
	if configClient == nil {
		return fmt.Errorf("failed to initialize Redis client for config DB %d", CONFIG_REDIS_DB)
	}

	rateLimiterClient := GetClient(RATE_LIMITER_DB)
	if rateLimiterClient == nil {
		return fmt.Errorf("failed to initialize Redis client for rate limiter DB %d", RATE_LIMITER_DB)
	}

	fmt.Println("All Redis clients initialized successfully")
	return nil
}
