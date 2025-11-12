package util

import (
	"context"
	"fmt"
	"log"
	"social-platform-backend/internal/domain/repository"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func GetPasswordChangedAt(redisClient *redis.Client, userRepo repository.UserRepository, userID uint64) (*time.Time, error) {
	ctx := context.Background()
	key := fmt.Sprintf("pwd_changed_at:%d", userID)

	// Try to get from cache first
	cachedValue, err := redisClient.Get(ctx, key).Result()
	if err == nil {
		// Cache hit
		if cachedValue == "nil" {
			return nil, nil
		}

		timestamp, err := strconv.ParseInt(cachedValue, 10, 64)
		if err != nil {
			log.Printf("[Err] Error parsing cached password_changed_at for user %d: %v", userID, err)
		} else {
			t := time.Unix(timestamp, 0)
			return &t, nil
		}
	} else if err != redis.Nil {
		log.Printf("[Warn] Error getting password_changed_at from cache for user %d: %v", userID, err)
	}

	// Cache miss or error, query database
	user, err := userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user from database: %w", err)
	}

	// Cache the result (TTL: 1 hour)
	if user.PasswordChangedAt != nil {
		cacheValue := strconv.FormatInt(user.PasswordChangedAt.Unix(), 10)
		if err := redisClient.Set(ctx, key, cacheValue, time.Hour).Err(); err != nil {
			log.Printf("[Warn] Error caching password_changed_at for user %d: %v", userID, err)
		}
	} else {
		// Cache "nil" value to prevent repeated DB queries
		if err := redisClient.Set(ctx, key, "nil", time.Hour).Err(); err != nil {
			log.Printf("[Warn] Error caching nil password_changed_at for user %d: %v", userID, err)
		}
	}

	return user.PasswordChangedAt, nil
}

// InvalidatePasswordCache removes the cached password_changed_at for a user (after password change)
func InvalidatePasswordCache(redisClient *redis.Client, userID uint64) error {
	ctx := context.Background()
	key := fmt.Sprintf("pwd_changed_at:%d", userID)

	if err := redisClient.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to invalidate password cache: %w", err)
	}

	log.Printf("[Info] Password cache invalidated for user %d", userID)
	return nil
}

func ValidatePasswordChangedAt(redisClient *redis.Client, userRepo repository.UserRepository, userID uint64, tokenPasswordChangedAt *int64) (bool, error) {
	// Get current password_changed_at from cache/DB
	currentPasswordChangedAt, err := GetPasswordChangedAt(redisClient, userRepo, userID)
	if err != nil {
		return false, fmt.Errorf("failed to get current password_changed_at: %w", err)
	}

	// If user has never changed password
	if currentPasswordChangedAt == nil {
		return tokenPasswordChangedAt == nil, nil
	}

	// If token doesn't have password_changed_at but user has changed password
	if tokenPasswordChangedAt == nil {
		return false, nil
	}

	currentUnix := currentPasswordChangedAt.Unix()
	return currentUnix == *tokenPasswordChangedAt, nil
}
