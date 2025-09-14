package router

import (
	"social-platform-backend/config"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB, conf *config.Config) *gin.Engine {
	router := gin.Default()
	// --- IGNORE ---

	return router
}
