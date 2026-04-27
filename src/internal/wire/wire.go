//go:build wireinject
// +build wireinject

package wire

import (
	"social-platform-backend/config"
	domainrepo "social-platform-backend/internal/domain/repository"
	dbrepository "social-platform-backend/internal/infrastructure/db/repository"
	"social-platform-backend/internal/interface/handler"
	"social-platform-backend/internal/service"

	"github.com/google/wire"
	"gorm.io/gorm"
)

type AppHandler struct {
	AuthHandler         *handler.AuthHandler
	UserHandler         *handler.UserHandler
	CommunityHandler    *handler.CommunityHandler
	PostHandler         *handler.PostHandler
	CommentHandler      *handler.CommentHandler
	MessageHandler      *handler.MessageHandler
	NotificationHandler *handler.NotificationHandler
	SSEHandler          *handler.SSEHandler
	ChatbotHandler      *handler.ChatbotHandler

	// Repos dùng trực tiếp trong route middleware
	UserRestrictionRepo domainrepo.UserRestrictionRepository
	PostRepo            domainrepo.PostRepository
}

var RepositorySet = wire.NewSet(
	dbrepository.NewUserRepository,
	dbrepository.NewUserVerificationRepository,
	dbrepository.NewPasswordResetRepository,
	dbrepository.NewRefreshTokenRepository,
	dbrepository.NewBotTaskRepository,
	dbrepository.NewCommunityRepository,
	dbrepository.NewSubscriptionRepository,
	dbrepository.NewCommunityModeratorRepository,
	dbrepository.NewPostRepository,
	dbrepository.NewPostVoteRepository,
	dbrepository.NewPostReportRepository,
	dbrepository.NewCommentRepository,
	dbrepository.NewCommentVoteRepository,
	dbrepository.NewCommentReportRepository,
	dbrepository.NewUserRestrictionRepository,
	wire.Bind(new(domainrepo.UserRestrictionRepository), new(*dbrepository.UserRestrictionRepositoryImpl)),
	dbrepository.NewConversationRepository,
	dbrepository.NewMessageRepository,
	dbrepository.NewMessageAttachmentRepository,
	dbrepository.NewNotificationRepository,
	dbrepository.NewNotificationSettingRepository,
	dbrepository.NewUserSavedPostRepository,
	dbrepository.NewUserInterestScoreRepository,
	dbrepository.NewUserTagPreferenceRepository,
	dbrepository.NewTagRepository,
	dbrepository.NewTopicRepository,
)

var ServiceSet = wire.NewSet(
	service.NewSSEService,
	service.NewBotTaskService,
	service.NewRecommendationService,
	service.NewNotificationService,
	service.NewAuthService,
	service.NewUserService,
	service.NewMessageService,
	service.NewPostService,
	service.NewCommentService,
	service.NewCommunityService,
	service.NewChatbotService,
)

var HandlerSet = wire.NewSet(
	handler.NewAuthHandler,
	handler.NewUserHandler,
	handler.NewCommunityHandler,
	handler.NewPostHandler,
	handler.NewCommentHandler,
	handler.NewMessageHandler,
	handler.NewNotificationHandler,
	handler.NewSSEHandler,
	handler.NewChatbotHandler,
)

var ProviderSet = wire.NewSet(
	RepositorySet,
	ServiceSet,
	HandlerSet,
	wire.Struct(new(AppHandler), "*"),
)

func InitAppContainer(db *gorm.DB, conf *config.Config) *AppHandler {
	wire.Build(ProviderSet)
	return nil
}
