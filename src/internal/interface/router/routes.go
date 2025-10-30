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

type AppHandler struct {
	userHandler      *handler.UserHandler
	communityHandler *handler.CommunityHandler
	postHandler      *handler.PostHandler
	commentHandler   *handler.CommentHandler
	authHandler      *handler.AuthHandler
}

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
	postVoteRepo := repository.NewPostVoteRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	commentVoteRepo := repository.NewCommentVoteRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, verificationRepo, passwordResetRepo, botTaskRepo, communityModeratorRepo)
	userService := service.NewUserService(userRepo, communityModeratorRepo)
	communityService := service.NewCommunityService(communityRepo, subscriptionRepo, communityModeratorRepo)
	postService := service.NewPostService(postRepo, communityRepo, postVoteRepo)
	commentService := service.NewCommentService(commentRepo, postRepo, commentVoteRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	communityHandler := handler.NewCommunityHandler(communityService)
	postHandler := handler.NewPostHandler(postService)
	commentHandler := handler.NewCommentHandler(commentService)

	appHandler := &AppHandler{
		userHandler:      userHandler,
		communityHandler: communityHandler,
		postHandler:      postHandler,
		commentHandler:   commentHandler,
		authHandler:      authHandler,
	}

	// Setup route groups
	api := router.Group("/api/v1")
	{
		setupPublicRoutes(api, appHandler)
		setupProtectedRoutes(api, appHandler, conf)
	}

	return router
}

func setupPublicRoutes(rg *gin.RouterGroup, appHandler *AppHandler) {
	// Health check
	rg.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	// Auth routes (public)
	auth := rg.Group("/auth")
	{
		auth.POST("/register", appHandler.authHandler.Register)
		auth.GET("/verify", appHandler.authHandler.VerifyEmail)
		auth.POST("/login", appHandler.authHandler.Login)
		auth.POST("/forgot-password", appHandler.authHandler.ForgotPassword)
		auth.GET("/verify-reset", appHandler.authHandler.VerifyResetToken)
		auth.POST("/reset-password", appHandler.authHandler.ResetPassword)
		auth.POST("/resend-verification", appHandler.authHandler.ResendVerificationEmail)
		auth.POST("/resend-reset-password", appHandler.authHandler.ResendResetPasswordEmail)
	}
	// Community routes (public)
	communities := rg.Group("/communities")
	{
		communities.GET("", appHandler.communityHandler.GetCommunities)
		communities.GET("/search", appHandler.communityHandler.SearchCommunities)
		communities.GET("/filter", appHandler.communityHandler.FilterCommunities)
		communities.GET("/:id", appHandler.communityHandler.GetCommunityByID)
		communities.GET("/:id/posts", appHandler.postHandler.GetPostsByCommunity)
		communities.POST("/verify-name", appHandler.communityHandler.VerifyCommunityName)
	}

	// Post routes (public)
	posts := rg.Group("/posts")
	{
		posts.GET("", appHandler.postHandler.GetAllPosts)
		posts.GET("/search", appHandler.postHandler.SearchPosts)
		posts.GET("/:postId/comments", appHandler.commentHandler.GetCommentsByPostID)
	}
}

func setupProtectedRoutes(rg *gin.RouterGroup, appHandler *AppHandler, conf *config.Config) {
	protected := rg.Group("")
	protected.Use(middleware.AuthMiddleware(conf))
	{
		users := protected.Group("/users")
		{
			users.GET("/me", appHandler.userHandler.GetCurrentUser)
			users.PUT("/me", appHandler.userHandler.UpdateUserProfile)
			users.PUT("/change-password", appHandler.userHandler.ChangePassword)
			users.GET("/config", appHandler.userHandler.GetUserConfig)
		}

		communities := protected.Group("/communities")
		{
			communities.POST("", appHandler.communityHandler.CreateCommunity)
			communities.POST("/:id/join", appHandler.communityHandler.JoinCommunity)
			communities.PUT("/:id", appHandler.communityHandler.UpdateCommunity)
			communities.DELETE("/:id", appHandler.communityHandler.DeleteCommunity)
			communities.GET("/:id/members", appHandler.communityHandler.GetCommunityMembers)
			communities.DELETE("/:id/members/:memberId", appHandler.communityHandler.RemoveMember)
			communities.GET("/:id/role", appHandler.communityHandler.GetUserRoleInCommunity)
		}

		posts := protected.Group("/posts")
		{
			posts.POST("", appHandler.postHandler.CreatePost)
			posts.PUT("/:id", appHandler.postHandler.UpdatePost)
			posts.DELETE("/:id", appHandler.postHandler.DeletePost)
			posts.POST("/:id/vote", appHandler.postHandler.VotePost)
			posts.DELETE("/:id/vote", appHandler.postHandler.UnvotePost)
		}

		comments := protected.Group("/comments")
		{
			comments.POST("", appHandler.commentHandler.CreateComment)
			comments.PUT("/:id", appHandler.commentHandler.UpdateComment)
			comments.DELETE("/:id", appHandler.commentHandler.DeleteComment)
			comments.POST("/:id/vote", appHandler.commentHandler.VoteComment)
			comments.DELETE("/:id/vote", appHandler.commentHandler.UnvoteComment)
		}
	}
}
