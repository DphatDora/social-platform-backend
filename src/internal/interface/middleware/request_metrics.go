package middleware

import (
	"social-platform-backend/package/logger"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestMetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logger.ObserveRequest(time.Since(start), c.Writer.Status())
	}
}
