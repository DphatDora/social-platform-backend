package cache

import (
	"context"
	"fmt"
	"log"
	"social-platform-backend/config"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func NewRedisClient(conf *config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", conf.Redis.Host, conf.Redis.Port),
		Password:     conf.Redis.Password,
		DB:           conf.Redis.DB,
		PoolSize:     conf.Redis.PoolSize,
		MinIdleConns: 5,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Printf("[Info] Redis connected successfully to %s:%s", conf.Redis.Host, conf.Redis.Port)
	redisClient = client
	return client, nil
}

func GetRedisClient() *redis.Client {
	return redisClient
}

func CloseRedis() error {
	if redisClient != nil {
		log.Println("[Info] Closing Redis connection...")
		return redisClient.Close()
	}
	return nil
}
