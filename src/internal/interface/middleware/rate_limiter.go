package middleware

import (
	"fmt"
	"log"
	"net/http"
	"social-platform-backend/internal/interface/dto/response"
	"social-platform-backend/package/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func LoginRateLimitMiddleware(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		if redisClient == nil {
			c.Next()
			return
		}

		clientIP := c.ClientIP()
		key := fmt.Sprintf("rate_limit:login:%s", clientIP)

		// Check rate limit (5reqs/min)
		allowed, err := util.CheckRateLimit(redisClient, key, 5, time.Minute)
		if err != nil {
			log.Printf("[Err] Error checking login rate limit for IP %s: %v", clientIP, err)
			c.Next()
			return
		}

		if !allowed {
			log.Printf("[Warn] Login rate limit exceeded for IP %s", clientIP)

			ttl, _ := util.GetRateLimitTTL(redisClient, key)

			c.JSON(http.StatusTooManyRequests, response.APIResponse{
				Success: false,
				Message: fmt.Sprintf("Too many login attempts. Please try again in %d seconds", int(ttl.Seconds())),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func APIRateLimitMiddleware(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		if redisClient == nil {
			c.Next()
			return
		}

		userID, err := util.GetUserIDFromContext(c)
		if err != nil {
			// If no user ID, use IP address as fallback
			clientIP := c.ClientIP()
			key := fmt.Sprintf("rate_limit:api:%s", clientIP)

			allowed, err := util.CheckRateLimit(redisClient, key, 30, time.Minute)
			if err != nil {
				log.Printf("[Err] Error checking API rate limit for IP %s: %v", clientIP, err)
				// Gracefully fallback: allow request to proceed when Redis is unavailable
				c.Next()
				return
			}

			if !allowed {
				log.Printf("[Warn] API rate limit exceeded for IP %s", clientIP)

				ttl, _ := util.GetRateLimitTTL(redisClient, key)

				c.JSON(http.StatusTooManyRequests, response.APIResponse{
					Success: false,
					Message: fmt.Sprintf("Too many requests. Please try again in %d seconds", int(ttl.Seconds())),
				})
				c.Abort()
				return
			}

			c.Next()
			return
		}

		// Rate limit by user ID
		key := fmt.Sprintf("rate_limit:api:user:%d", userID)

		allowed, err := util.CheckRateLimit(redisClient, key, 30, time.Minute)
		if err != nil {
			log.Printf("[Err] Error checking API rate limit for user %d: %v", userID, err)
			// Gracefully fallback: allow request to proceed when Redis is unavailable
			c.Next()
			return
		}

		if !allowed {
			log.Printf("[Warn] API rate limit exceeded for user %d", userID)

			ttl, _ := util.GetRateLimitTTL(redisClient, key)

			c.JSON(http.StatusTooManyRequests, response.APIResponse{
				Success: false,
				Message: fmt.Sprintf("Too many requests. Please try again in %d seconds", int(ttl.Seconds())),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
