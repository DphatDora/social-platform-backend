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
	communityRepo := repository.NewCommunityRepository(db)
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	communityModeratorRepo := repository.NewCommunityModeratorRepository(db)
	postRepo := repository.NewPostRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, verificationRepo, passwordResetRepo, botTaskRepo, communityModeratorRepo)
	userService := service.NewUserService(userRepo)
	communityService := service.NewCommunityService(communityRepo, subscriptionRepo, communityModeratorRepo)
	postService := service.NewPostService(postRepo, communityRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	communityHandler := handler.NewCommunityHandler(communityService)
	postHandler := handler.NewPostHandler(postService)

	// Setup route groups
	api := router.Group("/api/v1")
	{
		setupPublicRoutes(api, authHandler, communityHandler)
		setupProtectedRoutes(api, userHandler, communityHandler, postHandler, conf)
	}

	return router
}

func setupPublicRoutes(rg *gin.RouterGroup, authHandler *handler.AuthHandler, communityHandler *handler.CommunityHandler) {
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
		auth.POST("/resend-verification", authHandler.ResendVerificationEmail)
		auth.POST("/resend-reset-password", authHandler.ResendResetPasswordEmail)
	}
	// Community routes (public)
	communities := rg.Group("/communities")
	{
		communities.GET("", communityHandler.GetCommunities)
		communities.GET("/search", communityHandler.SearchCommunities)
		communities.GET("/filter", communityHandler.FilterCommunities)
		communities.GET("/:id", communityHandler.GetCommunityByID)
	}
}

func setupProtectedRoutes(rg *gin.RouterGroup, userHandler *handler.UserHandler, communityHandler *handler.CommunityHandler, postHandler *handler.PostHandler, conf *config.Config) {
	protected := rg.Group("")
	protected.Use(middleware.AuthMiddleware(conf))
	{
		users := protected.Group("/users")
		{
			users.GET("/me", userHandler.GetCurrentUser)
			users.PUT("/me", userHandler.UpdateUserProfile)
			users.PUT("/change-password", userHandler.ChangePassword)
		}

		communities := protected.Group("/communities")
		{
			communities.POST("", communityHandler.CreateCommunity)
			communities.POST("/:id/join", communityHandler.JoinCommunity)
			communities.PUT("/:id", communityHandler.UpdateCommunity)
			communities.DELETE("/:id", communityHandler.DeleteCommunity)
			communities.GET("/:id/members", communityHandler.GetCommunityMembers)
			communities.DELETE("/:id/members/:memberId", communityHandler.RemoveMember)
			communities.GET("/:id/role", communityHandler.GetUserRoleInCommunity)
		}

		posts := protected.Group("/posts")
		{
			posts.POST("", postHandler.CreatePost)
			posts.PUT("/:id", postHandler.UpdatePost)
			posts.DELETE("/:id", postHandler.DeletePost)
		}
	}
}
