package router

import (
	"social-platform-backend/config"
	"social-platform-backend/internal/infrastructure/db/repository"
	"social-platform-backend/internal/interface/handler"
	"social-platform-backend/internal/interface/middleware"
	"social-platform-backend/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB, conf *config.Config) *gin.Engine {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware(conf.App.Whitelist))

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	verificationRepo := repository.NewUserVerificationRepository(db)
	passwordResetRepo := repository.NewPasswordResetRepository(db)
	botTaskRepo := repository.NewBotTaskRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, verificationRepo, passwordResetRepo, botTaskRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)

	// Setup route groups
	api := router.Group("/api/v1")
	{
		setupPublicRoutes(api, authHandler)
		setupProtectedRoutes(api, authHandler, conf)
	}

	return router
}

func setupPublicRoutes(rg *gin.RouterGroup, authHandler *handler.AuthHandler) {
	// Health check
	rg.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	// Auth routes (public)
	auth := rg.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.GET("/verify", authHandler.VerifyEmail)
		auth.POST("/login", authHandler.Login)
		auth.POST("/forgot-password", authHandler.ForgotPassword)
		auth.GET("/verify-reset", authHandler.VerifyResetToken)
		auth.POST("/reset-password", authHandler.ResetPassword)
	}
}

func setupProtectedRoutes(rg *gin.RouterGroup, authHandler *handler.AuthHandler, conf *config.Config) {
	// protected := rg.Group("")
	// protected.Use(middleware.AuthMiddleware(conf))
	// {
	// 	users := protected.Group("/users")
	// 	{
	// 		// users.GET("/me", userHandler.GetCurrentUser)
	// 	}
	// }
}
