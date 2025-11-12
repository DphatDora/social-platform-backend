package util

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func CheckRateLimit(redisClient *redis.Client, key string, maxRequests int, window time.Duration) (bool, error) {
	ctx := context.Background()

	// Increment counter
	count, err := redisClient.Incr(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to increment rate limit counter: %w", err)
	}

	// Set expiration only on first request (when count == 1)
	if count == 1 {
		if err := redisClient.Expire(ctx, key, window).Err(); err != nil {
			return false, fmt.Errorf("failed to set rate limit expiration: %w", err)
		}
	}

	if count > int64(maxRequests) {
		return false, nil
	}

	return true, nil
}

func GetRateLimitTTL(redisClient *redis.Client, key string) (time.Duration, error) {
	ctx := context.Background()

	ttl, err := redisClient.TTL(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get rate limit TTL: %w", err)
	}

	return ttl, nil
}

// GetRemainingAttempts returns the number of remaining attempts for a key
func GetRemainingAttempts(redisClient *redis.Client, key string, maxRequests int) (int, error) {
	ctx := context.Background()

	count, err := redisClient.Get(ctx, key).Int()
	if err != nil {
		if err == redis.Nil {
			return maxRequests, nil
		}
		return 0, fmt.Errorf("failed to get rate limit count: %w", err)
	}

	remaining := maxRequests - count
	if remaining < 0 {
		remaining = 0
	}

	return remaining, nil
}
