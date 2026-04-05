package middleware

import (
	"net/http"
	"social-platform-backend/config"
	"strings"

	"github.com/gin-gonic/gin"
)

// LogDashboardAuth protects log dashboard routes when log.dashboardToken is set.
// Token may be sent as query ?token= or header X-Log-Dashboard-Token.
func LogDashboardAuth(conf *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		expected := strings.TrimSpace(conf.Log.DashboardToken)
		if expected == "" {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		t := strings.TrimSpace(c.Query("token"))
		if t == "" {
			t = strings.TrimSpace(c.GetHeader("X-Log-Dashboard-Token"))
		}
		if t != expected {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}
