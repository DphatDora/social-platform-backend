package router

import (
	"social-platform-backend/config"
	"social-platform-backend/internal/interface/handler"
	"social-platform-backend/internal/interface/middleware"
	"social-platform-backend/internal/wire"
	"strings"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(appHandler *wire.AppHandler, conf *config.Config) *gin.Engine {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware(conf.App.Whitelist))

	// Setup route groups
	api := router.Group("/api/v1")
	{
		api.Use(middleware.RequestMetricsMiddleware())
		setupPublicRoutes(api, appHandler, conf)
		setupProtectedRoutes(api, appHandler, conf)
	}

	if strings.TrimSpace(conf.Log.DashboardToken) != "" {
		dash := router.Group("")
		dash.Use(middleware.LogDashboardAuth(conf))
		dash.GET("/admin/logs", handler.ServeAdminLogsDashboard)

		apiAdmin := router.Group("/api/v1/admin")
		apiAdmin.Use(middleware.LogDashboardAuth(conf))
		apiAdmin.GET("/logs/files", handler.GetAdminLogFiles)
		apiAdmin.GET("/logs", handler.GetAdminLogs)
		apiAdmin.GET("/metrics", handler.GetAdminMetrics)
	}

	return router
}

func setupPublicRoutes(rg *gin.RouterGroup, appHandler *wire.AppHandler, conf *config.Config) {
	// Health check
	rg.Use(middleware.OptionalAuthMiddleware(conf))
	rg.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	auth := rg.Group("/auth")
	{
		auth.POST("/register", appHandler.AuthHandler.Register)
		auth.GET("/verify", appHandler.AuthHandler.VerifyEmail)
		auth.POST("/login", appHandler.AuthHandler.Login)
		auth.POST("/google-login", appHandler.AuthHandler.GoogleLogin)
		auth.POST("/refresh", appHandler.AuthHandler.RefreshToken)
		auth.POST("/logout", appHandler.AuthHandler.Logout)
		auth.POST("/forgot-password", appHandler.AuthHandler.ForgotPassword)
		auth.GET("/verify-reset", appHandler.AuthHandler.VerifyResetToken)
		auth.POST("/reset-password", appHandler.AuthHandler.ResetPassword)
		auth.POST("/resend-verification", appHandler.AuthHandler.ResendVerificationEmail)
		auth.POST("/resend-reset-password", appHandler.AuthHandler.ResendResetPasswordEmail)
	}

	communities := rg.Group("/communities")
	{
		communities.GET("", appHandler.CommunityHandler.GetCommunities)
		communities.GET("/search", appHandler.CommunityHandler.SearchCommunities)
		communities.GET("/filter", appHandler.CommunityHandler.FilterCommunities)
		communities.GET("/:id", appHandler.CommunityHandler.GetCommunityByID)
		communities.GET("/:id/posts", appHandler.PostHandler.GetPostsByCommunity)
		communities.POST("/verify-name", appHandler.CommunityHandler.VerifyCommunityName)
		communities.GET("/topics", appHandler.CommunityHandler.GetAllTopics)
	}

	posts := rg.Group("/posts")
	{
		posts.GET("", appHandler.PostHandler.GetAllPosts)
		posts.GET("/search", appHandler.PostHandler.SearchPosts)
		posts.GET("/:id", appHandler.PostHandler.GetPostDetail)
		posts.GET("/:id/comments", appHandler.CommentHandler.GetCommentsOnPost)
		posts.GET("/tags", appHandler.PostHandler.GetAllTags)
	}

	users := rg.Group("/users")
	{
		users.GET("/search", appHandler.UserHandler.SearchUsers)
		users.GET("/:id", appHandler.UserHandler.GetUserByID)
		users.GET("/:id/posts", appHandler.PostHandler.GetPostsByUser)
		users.GET("/:id/comments", appHandler.CommentHandler.GetCommentsByUser)
		users.GET("/:id/badge-history", appHandler.UserHandler.GetUserBadgeHistory)
		// users.GET("/:id/communities/super-admin", appHandler.UserHandler.GetUserSuperAdminCommunities)
		// users.GET("/:id/communities/admin", appHandler.UserHandler.GetUserAdminCommunities)
		users.GET("/:id/communities/joined", appHandler.UserHandler.GetUserJoinedCommunities)
	}
}

func setupProtectedRoutes(rg *gin.RouterGroup, appHandler *wire.AppHandler, conf *config.Config) {
	protected := rg.Group("")
	protected.Use(middleware.AuthMiddleware(conf))
	{
		users := protected.Group("/users")
		{
			users.GET("/me", appHandler.UserHandler.GetCurrentUser)
			users.PUT("/me", appHandler.UserHandler.UpdateUserProfile)
			users.PUT("/change-password", appHandler.UserHandler.ChangePassword)
			users.GET("/config", appHandler.UserHandler.GetUserConfig)
			users.GET("/notification-settings", appHandler.NotificationHandler.GetNotificationSettings)
			users.PATCH("/notification-settings", appHandler.NotificationHandler.UpdateNotificationSetting)
			users.GET("/saved-posts", appHandler.UserHandler.GetUserSavedPosts)
			users.POST("/saved-posts", appHandler.UserHandler.CreateUserSavedPost)
			users.PATCH("/saved-posts/:postId", appHandler.UserHandler.UpdateUserSavedPostFollowStatus)
			users.DELETE("/saved-posts/:postId", appHandler.UserHandler.DeleteUserSavedPost)
		}

		communities := protected.Group("/communities")
		{
			communities.POST("", appHandler.CommunityHandler.CreateCommunity)
			communities.POST("/:id/join", appHandler.CommunityHandler.JoinCommunity)
			communities.DELETE("/:id/join", appHandler.CommunityHandler.UnjoinCommunity)
			communities.PUT("/:id", appHandler.CommunityHandler.UpdateCommunity)
			communities.DELETE("/:id", appHandler.CommunityHandler.DeleteCommunity)
			communities.GET("/:id/members", appHandler.CommunityHandler.GetCommunityMembers)
			communities.PUT("/:id/moderators/:userId", appHandler.CommunityHandler.UpdateMemberRole)
			communities.DELETE("/:id/members/:memberId", appHandler.CommunityHandler.RemoveMember)
			communities.GET("/:id/role", appHandler.CommunityHandler.GetUserRoleInCommunity)
			communities.PATCH("/:id/requires-post-approval", appHandler.CommunityHandler.UpdateRequiresPostApproval)
			communities.PATCH("/:id/requires-member-approval", appHandler.CommunityHandler.UpdateRequiresMemberApproval)
			communities.GET("/:id/manage/posts", appHandler.CommunityHandler.GetCommunityPostsForModerator)
			communities.PATCH("/:id/manage/posts/:postId/status", appHandler.CommunityHandler.UpdatePostStatusByModerator)
			communities.DELETE("/:id/manage/posts/:postId", appHandler.CommunityHandler.DeletePostByModerator)
			communities.DELETE("/:id/manage/comments/:commentId", appHandler.CommunityHandler.DeleteCommentByModerator)
			communities.GET("/:id/manage/reports", appHandler.CommunityHandler.GetCommunityPostReports)
			communities.DELETE("/:id/manage/reports/:reportId", appHandler.CommunityHandler.DeletePostReport)
			communities.GET("/:id/manage/comment-reports", appHandler.CommunityHandler.GetCommunityCommentReports)
			communities.DELETE("/:id/manage/comment-reports/:reportId", appHandler.CommunityHandler.DeleteCommentReport)
			communities.POST("/:id/manage/ban-user", appHandler.CommunityHandler.BanUser)
			communities.GET("/:id/manage/restrictions/user/:userId", appHandler.CommunityHandler.GetUserRestrictionHistory)
			communities.DELETE("/:id/manage/restrictions/:restrictionId", appHandler.CommunityHandler.RemoveUserRestriction)
			communities.PATCH("/:id/manage/subscriptions/:userId/status", appHandler.CommunityHandler.UpdateSubscriptionStatus)
		}

		posts := protected.Group("/posts")
		{
			posts.POST("", middleware.CheckUserRestrictionForPostMiddleware(appHandler.UserRestrictionRepo), appHandler.PostHandler.CreatePost)
			posts.PUT("/:id", appHandler.PostHandler.UpdatePost)
			posts.DELETE("/:id", appHandler.PostHandler.DeletePost)
			posts.POST("/:id/vote", appHandler.PostHandler.VotePost)
			posts.DELETE("/:id/vote", appHandler.PostHandler.UnvotePost)
			posts.POST("/:id/poll/vote", appHandler.PostHandler.VotePoll)
			posts.DELETE("/:id/poll/vote", appHandler.PostHandler.UnvotePoll)
			posts.POST("/:id/report", appHandler.PostHandler.ReportPost)
		}

		comments := protected.Group("/comments")
		{
			comments.POST("", middleware.CheckUserRestrictionForCommentMiddleware(appHandler.UserRestrictionRepo, appHandler.PostRepo), appHandler.CommentHandler.CreateComment)
			comments.PUT("/:id", appHandler.CommentHandler.UpdateComment)
			comments.DELETE("/:id", appHandler.CommentHandler.DeleteComment)
			comments.POST("/:id/vote", appHandler.CommentHandler.VoteComment)
			comments.DELETE("/:id/vote", appHandler.CommentHandler.UnvoteComment)
			comments.POST("/:id/report", appHandler.CommentHandler.ReportComment)
		}

		messages := protected.Group("/messages")
		{
			messages.POST("", appHandler.MessageHandler.SendMessage)
			messages.GET("/conversations", appHandler.MessageHandler.GetConversations)
			messages.GET("/conversations/:conversationId/messages", appHandler.MessageHandler.GetMessages)
			messages.PATCH("/conversations/:conversationId/read", appHandler.MessageHandler.MarkConversationAsRead)
			messages.PATCH("/:messageId/read", appHandler.MessageHandler.MarkAsRead)
			messages.DELETE("/:messageId", appHandler.MessageHandler.DeleteMessage)
		}

		notifications := protected.Group("/notifications")
		{
			notifications.GET("", appHandler.NotificationHandler.GetNotifications)
			notifications.GET("/unread-count", appHandler.NotificationHandler.GetUnreadCount)
			notifications.PATCH("/:id/read", appHandler.NotificationHandler.MarkAsRead)
			notifications.PATCH("/read-all", appHandler.NotificationHandler.MarkAllAsRead)
			notifications.DELETE("/:id", appHandler.NotificationHandler.DeleteNotification)
		}

		sse := protected.Group("")
		{
			sse.GET("/stream", appHandler.SSEHandler.Stream)
			sse.GET("/conversations/:conversationId", appHandler.SSEHandler.StreamConversationMessages)
		}

		chatbot := protected.Group("/chatbot")
		{
			chatbot.POST("/stream", appHandler.ChatbotHandler.StreamChat)
		}
	}
}
