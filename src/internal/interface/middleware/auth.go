package middleware

import (
	"log"
	"net/http"
	"social-platform-backend/config"
	"social-platform-backend/internal/interface/dto/response"
	"social-platform-backend/package/util"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT token and sets user information in context
func AuthMiddleware(conf *config.Config) gin.HandlerFunc {
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

		// Set user information in context
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)

		c.Next()
	}
}

// OptionalAuthMiddleware validates JWT token if present, but allows request to continue if not
func OptionalAuthMiddleware(conf *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// No token provided, continue without setting user context
			c.Next()
			return
		}

		// Check if it's a Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			// Invalid format, but allow request to continue
			c.Next()
			return
		}

		tokenString := parts[1]

		// Verify JWT token
		claims, err := util.VerifyJWT(tokenString, conf.Auth.JWTSecret)
		if err != nil {
			// Invalid token, but allow request to continue
			log.Printf("[Warn] Invalid JWT token in optional auth: %v", err)
			c.Next()
			return
		}

		// Set user information in context if token is valid
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)

		c.Next()
	}
}
