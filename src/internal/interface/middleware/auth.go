package middleware

import (
	"log"
	"net/http"
	"social-platform-backend/config"
	"social-platform-backend/internal/domain/repository"
	"social-platform-backend/internal/interface/dto/response"
	"social-platform-backend/package/util"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func AuthMiddleware(conf *config.Config, redisClient *redis.Client, userRepo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Printf("[Err] Missing Authorization header")
			c.JSON(http.StatusUnauthorized, response.APIResponse{
				Success: false,
				Message: "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Check if it's a Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Printf("[Err] Invalid Authorization header format")
			c.JSON(http.StatusUnauthorized, response.APIResponse{
				Success: false,
				Message: "Invalid authorization header format. Expected 'Bearer <token>'",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Verify JWT token
		claims, err := util.VerifyJWT(tokenString, conf.Auth.JWTSecret)
		if err != nil {
			log.Printf("[Err] Invalid JWT token: %v", err)
			c.JSON(http.StatusUnauthorized, response.APIResponse{
				Success: false,
				Message: "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Validate password_changed_at from token against current value
		if redisClient != nil && userRepo != nil {
			isValid, err := util.ValidatePasswordChangedAt(redisClient, userRepo, claims.UserID, claims.PasswordChangedAt)
			if err != nil {
				log.Printf("[Err] Error validating password_changed_at for user %d: %v", claims.UserID, err)
				c.JSON(http.StatusInternalServerError, response.APIResponse{
					Success: false,
					Message: "Failed to validate token",
				})
				c.Abort()
				return
			}

			if !isValid {
				log.Printf("[Warn] Token invalidated due to password change for user %d", claims.UserID)
				c.JSON(http.StatusUnauthorized, response.APIResponse{
					Success: false,
					Message: "Token is no longer valid. Please login again",
				})
				c.Abort()
				return
			}
		}

		// Set user information in context
		c.Set("userID", claims.UserID)

		c.Next()
	}
}

func OptionalAuthMiddleware(conf *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		tokenString := parts[1]

		// Verify JWT token
		claims, err := util.VerifyJWT(tokenString, conf.Auth.JWTSecret)
		if err != nil {
			log.Printf("[Warn] Invalid JWT token in optional auth: %v", err)
			c.Next()
			return
		}

		c.Set("userID", claims.UserID)

		c.Next()
	}
}
