package router

import (
	"social-platform-backend/config"
	domainRepository "social-platform-backend/internal/domain/repository"
	"social-platform-backend/internal/infrastructure/db/repository"
	"social-platform-backend/internal/interface/handler"
	"social-platform-backend/internal/interface/middleware"
	"social-platform-backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AppHandler struct {
	userHandler         *handler.UserHandler
	communityHandler    *handler.CommunityHandler
	postHandler         *handler.PostHandler
	commentHandler      *handler.CommentHandler
	authHandler         *handler.AuthHandler
	messageHandler      *handler.MessageHandler
	notificationHandler *handler.NotificationHandler
	sseHandler          *handler.SSEHandler
}

func SetupRoutes(db *gorm.DB, redisClient *redis.Client, conf *config.Config) *gin.Engine {
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
	postReportRepo := repository.NewPostReportRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	commentVoteRepo := repository.NewCommentVoteRepository(db)
	conversationRepo := repository.NewConversationRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	messageAttachmentRepo := repository.NewMessageAttachmentRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)
	notificationSettingRepo := repository.NewNotificationSettingRepository(db)
	userSavedPostRepo := repository.NewUserSavedPostRepository(db)
	userInterestScoreRepo := repository.NewUserInterestScoreRepository(db)
	userTagPrefRepo := repository.NewUserTagPreferenceRepository(db)
	tagRepo := repository.NewTagRepository(db)
	topicRepo := repository.NewTopicRepository(db)

	// Initialize services
	sseService := service.NewSSEService()
	botTaskService := service.NewBotTaskService(botTaskRepo)
	recommendService := service.NewRecommendationService(userInterestScoreRepo, userTagPrefRepo, postRepo, communityRepo)
	notificationService := service.NewNotificationService(notificationRepo, notificationSettingRepo, botTaskRepo, userRepo, sseService, botTaskService)
	authService := service.NewAuthService(userRepo, verificationRepo, passwordResetRepo, botTaskRepo, communityModeratorRepo, notificationSettingRepo, botTaskService, redisClient)
	userService := service.NewUserService(userRepo, communityRepo, communityModeratorRepo, userSavedPostRepo, postRepo, botTaskService, redisClient)
	messageService := service.NewMessageService(conversationRepo, messageRepo, messageAttachmentRepo, userRepo, sseService)
	postService := service.NewPostService(postRepo, communityRepo, subscriptionRepo, postVoteRepo, postReportRepo, botTaskRepo, userRepo, tagRepo, notificationService, botTaskService, recommendService)
	commentService := service.NewCommentService(commentRepo, postRepo, commentVoteRepo, botTaskRepo, userRepo, userSavedPostRepo, notificationService, botTaskService)
	communityService := service.NewCommunityService(communityRepo, subscriptionRepo, communityModeratorRepo, postRepo, postReportRepo, topicRepo, notificationService, botTaskService)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	communityHandler := handler.NewCommunityHandler(communityService)
	postHandler := handler.NewPostHandler(postService)
	commentHandler := handler.NewCommentHandler(commentService)
	messageHandler := handler.NewMessageHandler(messageService)
	notificationHandler := handler.NewNotificationHandler(notificationService)
	sseHandler := handler.NewSSEHandler(sseService)

	appHandler := &AppHandler{
		userHandler:         userHandler,
		communityHandler:    communityHandler,
		postHandler:         postHandler,
		commentHandler:      commentHandler,
		authHandler:         authHandler,
		messageHandler:      messageHandler,
		notificationHandler: notificationHandler,
		sseHandler:          sseHandler,
	}

	// Setup route groups
	api := router.Group("/api/v1")
	{
		setupPublicRoutes(api, appHandler, conf, redisClient)
		setupProtectedRoutes(api, appHandler, conf, redisClient, userRepo)
	}

	return router
}

func setupPublicRoutes(rg *gin.RouterGroup, appHandler *AppHandler, conf *config.Config, redisClient *redis.Client) {
	// Health check
	rg.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	auth := rg.Group("/auth")
	{
		auth.POST("/register", appHandler.authHandler.Register)
		auth.GET("/verify", appHandler.authHandler.VerifyEmail)
		auth.POST("/login", middleware.LoginRateLimitMiddleware(redisClient), appHandler.authHandler.Login)
		auth.POST("/google-login", middleware.LoginRateLimitMiddleware(redisClient), appHandler.authHandler.GoogleLogin)
		auth.POST("/forgot-password", appHandler.authHandler.ForgotPassword)
		auth.GET("/verify-reset", appHandler.authHandler.VerifyResetToken)
		auth.POST("/reset-password", appHandler.authHandler.ResetPassword)
		auth.POST("/resend-verification", appHandler.authHandler.ResendVerificationEmail)
		auth.POST("/resend-reset-password", appHandler.authHandler.ResendResetPasswordEmail)
	}

	communities := rg.Group("/communities")
	communities.Use(middleware.OptionalAuthMiddleware(conf))
	{
		communities.GET("", appHandler.communityHandler.GetCommunities)
		communities.GET("/search", appHandler.communityHandler.SearchCommunities)
		communities.GET("/filter", appHandler.communityHandler.FilterCommunities)
		communities.GET("/:id", appHandler.communityHandler.GetCommunityByID)
		communities.GET("/:id/posts", appHandler.postHandler.GetPostsByCommunity)
		communities.POST("/verify-name", appHandler.communityHandler.VerifyCommunityName)
		communities.GET("/topics", appHandler.communityHandler.GetAllTopics)
	}

	posts := rg.Group("/posts")
	posts.Use(middleware.OptionalAuthMiddleware(conf))
	{
		posts.GET("", appHandler.postHandler.GetAllPosts)
		posts.GET("/search", appHandler.postHandler.SearchPosts)
		posts.GET("/:id", appHandler.postHandler.GetPostDetail)
		posts.GET("/:id/comments", appHandler.commentHandler.GetCommentsOnPost)
		posts.GET("/tags", appHandler.postHandler.GetAllTags)
	}

	users := rg.Group("/users")
	{
		users.GET("/search", appHandler.userHandler.SearchUsers)
		users.GET("/:id", appHandler.userHandler.GetUserByID)
		users.GET("/:id/posts", appHandler.postHandler.GetPostsByUser)
		users.GET("/:id/comments", appHandler.commentHandler.GetCommentsByUser)
		users.GET("/:id/badge-history", appHandler.userHandler.GetUserBadgeHistory)
		users.GET("/:id/communities/super-admin", appHandler.userHandler.GetUserSuperAdminCommunities)
		users.GET("/:id/communities/admin", appHandler.userHandler.GetUserAdminCommunities)
	}
}

func setupProtectedRoutes(rg *gin.RouterGroup, appHandler *AppHandler, conf *config.Config, redisClient *redis.Client, userRepo domainRepository.UserRepository) {
	protected := rg.Group("")
	protected.Use(middleware.AuthMiddleware(conf, redisClient, userRepo))
	{
		users := protected.Group("/users")
		{
			users.GET("/me", appHandler.userHandler.GetCurrentUser)
			users.PUT("/me", appHandler.userHandler.UpdateUserProfile)
			users.PUT("/change-password", appHandler.userHandler.ChangePassword)
			users.GET("/config", appHandler.userHandler.GetUserConfig)
			users.GET("/notification-settings", appHandler.notificationHandler.GetNotificationSettings)
			users.PATCH("/notification-settings", appHandler.notificationHandler.UpdateNotificationSetting)
			users.GET("/saved-posts", appHandler.userHandler.GetUserSavedPosts)
			users.POST("/saved-posts", appHandler.userHandler.CreateUserSavedPost)
			users.PATCH("/saved-posts/:postId", appHandler.userHandler.UpdateUserSavedPostFollowStatus)
			users.DELETE("/saved-posts/:postId", appHandler.userHandler.DeleteUserSavedPost)
		}

		communities := protected.Group("/communities")
		communities.Use(middleware.APIRateLimitMiddleware(redisClient))
		{
			communities.POST("", appHandler.communityHandler.CreateCommunity)
			communities.POST("/:id/join", appHandler.communityHandler.JoinCommunity)
			communities.DELETE("/:id/join", appHandler.communityHandler.UnjoinCommunity)
			communities.PUT("/:id", appHandler.communityHandler.UpdateCommunity)
			communities.DELETE("/:id", appHandler.communityHandler.DeleteCommunity)
			communities.GET("/:id/members", appHandler.communityHandler.GetCommunityMembers)
			communities.PUT("/:id/moderators/:userId", appHandler.communityHandler.UpdateMemberRole)
			communities.DELETE("/:id/members/:memberId", appHandler.communityHandler.RemoveMember)
			communities.GET("/:id/role", appHandler.communityHandler.GetUserRoleInCommunity)
			communities.PATCH("/:id/requires-post-approval", appHandler.communityHandler.UpdateRequiresPostApproval)
			communities.PATCH("/:id/requires-member-approval", appHandler.communityHandler.UpdateRequiresMemberApproval)
			communities.GET("/:id/manage/posts", appHandler.communityHandler.GetCommunityPostsForModerator)
			communities.PATCH("/:id/manage/posts/:postId/status", appHandler.communityHandler.UpdatePostStatusByModerator)
			communities.DELETE("/:id/manage/posts/:postId", appHandler.communityHandler.DeletePostByModerator)
			communities.GET("/:id/manage/reports", appHandler.communityHandler.GetCommunityPostReports)
			communities.DELETE("/:id/manage/reports/:reportId", appHandler.communityHandler.DeletePostReport)
			communities.PATCH("/:id/manage/subscriptions/:userId/status", appHandler.communityHandler.UpdateSubscriptionStatus)
		}

		posts := protected.Group("/posts")
		posts.Use(middleware.APIRateLimitMiddleware(redisClient))
		{
			posts.POST("", appHandler.postHandler.CreatePost)
			posts.PUT("/:id", appHandler.postHandler.UpdatePost)
			posts.DELETE("/:id", appHandler.postHandler.DeletePost)
			posts.POST("/:id/vote", appHandler.postHandler.VotePost)
			posts.DELETE("/:id/vote", appHandler.postHandler.UnvotePost)
			posts.POST("/:id/poll/vote", appHandler.postHandler.VotePoll)
			posts.DELETE("/:id/poll/vote", appHandler.postHandler.UnvotePoll)
			posts.POST("/:id/report", appHandler.postHandler.ReportPost)
		}

		comments := protected.Group("/comments")
		comments.Use(middleware.APIRateLimitMiddleware(redisClient))
		{
			comments.POST("", appHandler.commentHandler.CreateComment)
			comments.PUT("/:id", appHandler.commentHandler.UpdateComment)
			comments.DELETE("/:id", appHandler.commentHandler.DeleteComment)
			comments.POST("/:id/vote", appHandler.commentHandler.VoteComment)
			comments.DELETE("/:id/vote", appHandler.commentHandler.UnvoteComment)
		}

		messages := protected.Group("/messages")
		{
			messages.POST("", appHandler.messageHandler.SendMessage)
			messages.GET("/conversations", appHandler.messageHandler.GetConversations)
			messages.GET("/conversations/:conversationId/messages", appHandler.messageHandler.GetMessages)
			messages.PATCH("/conversations/:conversationId/read", appHandler.messageHandler.MarkConversationAsRead)
			messages.PATCH("/:messageId/read", appHandler.messageHandler.MarkAsRead)
			messages.DELETE("/:messageId", appHandler.messageHandler.DeleteMessage)
		}

		notifications := protected.Group("/notifications")
		{
			notifications.GET("", appHandler.notificationHandler.GetNotifications)
			notifications.GET("/unread-count", appHandler.notificationHandler.GetUnreadCount)
			notifications.PATCH("/:id/read", appHandler.notificationHandler.MarkAsRead)
			notifications.PATCH("/read-all", appHandler.notificationHandler.MarkAllAsRead)
			notifications.DELETE("/:id", appHandler.notificationHandler.DeleteNotification)
		}

		sse := protected.Group("")
		{
			sse.GET("/stream", appHandler.sseHandler.Stream)
			sse.GET("/conversations/:conversationId", appHandler.sseHandler.StreamConversationMessages)
		}
	}
}
