package service

import (
	"encoding/json"
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/interface/dto/request"

	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) IsEmailExisted(email string) (bool, error) {
	args := m.Called(email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) CreateUser(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByID(id uint64) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(email string) (*model.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByGoogleID(googleID string) (*model.User, error) {
	args := m.Called(googleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) ActivateUser(id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) UpdatePasswordAndSetChangedAt(id uint64, hashedPassword string) error {
	args := m.Called(id, hashedPassword)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateUserProfile(id uint64, updateUser *request.UpdateUserProfileRequest) error {
	args := m.Called(id, updateUser)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateAuthProvider(userID uint64, provider string) error {
	args := m.Called(userID, provider)
	return args.Error(0)
}

func (m *MockUserRepository) LinkGoogleAccount(userID uint64, googleID string, provider string) error {
	args := m.Called(userID, googleID, provider)
	return args.Error(0)
}

func (m *MockUserRepository) GetLatestUserBadge(userID uint64) (*model.UserBadge, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserBadge), args.Error(1)
}

func (m *MockUserRepository) GetUserPostCount(userID uint64) (uint64, error) {
	args := m.Called(userID)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockUserRepository) GetUserCommentCount(userID uint64) (uint64, error) {
	args := m.Called(userID)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockUserRepository) GetUserBadgeHistory(userID uint64) ([]*model.UserBadge, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.UserBadge), args.Error(1)
}

func (m *MockUserRepository) SearchUsers(searchTerm string, page, limit int) ([]*model.User, int64, error) {
	args := m.Called(searchTerm, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.User), args.Get(1).(int64), args.Error(2)
}

type MockUserVerificationRepository struct {
	mock.Mock
}

func (m *MockUserVerificationRepository) CreateVerification(verification *model.UserVerification) error {
	args := m.Called(verification)
	return args.Error(0)
}

func (m *MockUserVerificationRepository) GetVerificationByToken(token string) (*model.UserVerification, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserVerification), args.Error(1)
}

func (m *MockUserVerificationRepository) DeleteVerification(id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserVerificationRepository) DeleteVerificationByUserID(userID uint64) error {
	args := m.Called(userID)
	return args.Error(0)
}

type MockPasswordResetRepository struct {
	mock.Mock
}

func (m *MockPasswordResetRepository) CreatePasswordReset(reset *model.PasswordReset) error {
	args := m.Called(reset)
	return args.Error(0)
}

func (m *MockPasswordResetRepository) GetPasswordResetByToken(token string) (*model.PasswordReset, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.PasswordReset), args.Error(1)
}

func (m *MockPasswordResetRepository) DeletePasswordReset(id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPasswordResetRepository) DeletePasswordResetByUserID(userID uint64) error {
	args := m.Called(userID)
	return args.Error(0)
}

type MockNotificationSettingRepository struct {
	mock.Mock
}

func (m *MockNotificationSettingRepository) CreateNotificationSetting(setting *model.NotificationSetting) error {
	args := m.Called(setting)
	return args.Error(0)
}

func (m *MockNotificationSettingRepository) CreateNotificationSettings(settings []*model.NotificationSetting) error {
	args := m.Called(settings)
	return args.Error(0)
}

func (m *MockNotificationSettingRepository) GetUserNotificationSetting(userID uint64, action string) (*model.NotificationSetting, error) {
	args := m.Called(userID, action)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.NotificationSetting), args.Error(1)
}

func (m *MockNotificationSettingRepository) GetUserNotificationSettings(userID uint64) ([]*model.NotificationSetting, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.NotificationSetting), args.Error(1)
}

func (m *MockNotificationSettingRepository) UpdateNotificationSetting(setting *model.NotificationSetting) error {
	args := m.Called(setting)
	return args.Error(0)
}

func (m *MockNotificationSettingRepository) UpsertNotificationSetting(setting *model.NotificationSetting) error {
	args := m.Called(setting)
	return args.Error(0)
}

type MockBotTaskRepository struct {
	mock.Mock
}

func (m *MockBotTaskRepository) CreateBotTask(botTask *model.BotTask) error {
	args := m.Called(botTask)
	return args.Error(0)
}

type MockCommunityRepository struct {
	mock.Mock
}

func (m *MockCommunityRepository) CreateCommunity(community *model.Community) error {
	args := m.Called(community)
	return args.Error(0)
}

func (m *MockCommunityRepository) GetCommunityByID(id uint64) (*model.Community, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Community), args.Error(1)
}

func (m *MockCommunityRepository) GetCommunityWithMemberCount(id uint64) (*model.Community, int64, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).(*model.Community), args.Get(1).(int64), args.Error(2)
}

func (m *MockCommunityRepository) GetCommunityByIDWithUserSubscription(communityID uint64, userID *uint64) (*model.Community, int64, error) {
	args := m.Called(communityID, userID)
	if args.Get(0) == nil {
		if args.Get(1) == nil {
			return nil, 0, args.Error(2)
		}
		// Handle both int64 and uint64
		var memberCount int64
		switch v := args.Get(1).(type) {
		case int64:
			memberCount = v
		case uint64:
			memberCount = int64(v)
		default:
			memberCount = 0
		}
		return nil, memberCount, args.Error(2)
	}
	var memberCount int64
	switch v := args.Get(1).(type) {
	case int64:
		memberCount = v
	case uint64:
		memberCount = int64(v)
	default:
		memberCount = 0
	}
	return args.Get(0).(*model.Community), memberCount, args.Error(2)
}

func (m *MockCommunityRepository) UpdateCommunity(id uint64, updateCommunity *request.UpdateCommunityRequest) error {
	args := m.Called(id, updateCommunity)
	return args.Error(0)
}

func (m *MockCommunityRepository) DeleteCommunity(id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockCommunityRepository) GetCommunities(page, limit int, userID *uint64) ([]*model.Community, int64, error) {
	args := m.Called(page, limit, userID)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Community), args.Get(1).(int64), args.Error(2)
}

func (m *MockCommunityRepository) SearchCommunitiesByName(name string, page, limit int, userID *uint64) ([]*model.Community, int64, error) {
	args := m.Called(name, page, limit, userID)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Community), args.Get(1).(int64), args.Error(2)
}

func (m *MockCommunityRepository) FilterCommunities(sortBy string, isPrivate *bool, topics []string, page, limit int, userID *uint64) ([]*model.Community, int64, error) {
	args := m.Called(sortBy, isPrivate, topics, page, limit, userID)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Community), args.Get(1).(int64), args.Error(2)
}

func (m *MockCommunityRepository) GetCommunitiesByCreatorID(creatorID uint64) ([]*model.Community, error) {
	args := m.Called(creatorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Community), args.Error(1)
}

func (m *MockCommunityRepository) IsCommunityNameExists(name string) (bool, error) {
	args := m.Called(name)
	return args.Bool(0), args.Error(1)
}

func (m *MockCommunityRepository) UpdateRequiresPostApproval(id uint64, requiresPostApproval bool) error {
	args := m.Called(id, requiresPostApproval)
	return args.Error(0)
}

func (m *MockCommunityRepository) UpdateRequiresMemberApproval(id uint64, requiresMemberApproval bool) error {
	args := m.Called(id, requiresMemberApproval)
	return args.Error(0)
}

type MockCommunityModeratorRepository struct {
	mock.Mock
}

func (m *MockCommunityModeratorRepository) CreateModerator(moderator *model.CommunityModerator) error {
	args := m.Called(moderator)
	return args.Error(0)
}

func (m *MockCommunityModeratorRepository) DeleteModerator(communityID, userID uint64) error {
	args := m.Called(communityID, userID)
	return args.Error(0)
}

func (m *MockCommunityModeratorRepository) GetModeratorRole(communityID, userID uint64) (string, error) {
	args := m.Called(communityID, userID)
	return args.String(0), args.Error(1)
}

func (m *MockCommunityModeratorRepository) GetModeratorCommunitiesByUserID(userID uint64) ([]*model.CommunityModerator, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.CommunityModerator), args.Error(1)
}

func (m *MockCommunityModeratorRepository) GetCommunityModerators(communityID uint64) ([]*model.CommunityModerator, error) {
	args := m.Called(communityID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.CommunityModerator), args.Error(1)
}

func (m *MockCommunityModeratorRepository) UpsertModerator(moderator *model.CommunityModerator) error {
	args := m.Called(moderator)
	return args.Error(0)
}

func (m *MockCommunityModeratorRepository) GetSuperAdminCommunitiesByUserID(userID uint64) ([]*model.CommunityModerator, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.CommunityModerator), args.Error(1)
}

type MockUserSavedPostRepository struct {
	mock.Mock
}

func (m *MockUserSavedPostRepository) GetUserSavedPostCount(userID uint64) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

type MockPostRepository struct {
	mock.Mock
}

func (m *MockPostRepository) CreatePost(post *model.Post) error {
	args := m.Called(post)
	return args.Error(0)
}

func (m *MockPostRepository) GetPostByID(id uint64) (*model.Post, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Post), args.Error(1)
}

func (m *MockPostRepository) GetPostDetailByID(id uint64, userID *uint64) (*model.Post, error) {
	args := m.Called(id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Post), args.Error(1)
}

func (m *MockPostRepository) UpdatePostText(id uint64, updatePost *request.UpdatePostTextRequest) error {
	args := m.Called(id, updatePost)
	return args.Error(0)
}

func (m *MockPostRepository) UpdatePostLink(id uint64, updatePost *request.UpdatePostLinkRequest) error {
	args := m.Called(id, updatePost)
	return args.Error(0)
}

func (m *MockPostRepository) UpdatePostMedia(id uint64, updatePost *request.UpdatePostMediaRequest) error {
	args := m.Called(id, updatePost)
	return args.Error(0)
}

func (m *MockPostRepository) UpdatePostPoll(id uint64, updatePost *request.UpdatePostPollRequest) error {
	args := m.Called(id, updatePost)
	return args.Error(0)
}

func (m *MockPostRepository) UpdatePostStatus(id uint64, status string) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *MockPostRepository) UpdatePollData(postID uint64, pollData *json.RawMessage) error {
	args := m.Called(postID, pollData)
	return args.Error(0)
}

func (m *MockPostRepository) DeletePost(id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPostRepository) GetAllPosts(sortBy string, page, limit int, tags []string, userID *uint64) ([]*model.Post, int64, error) {
	args := m.Called(sortBy, page, limit, tags, userID)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Post), args.Get(1).(int64), args.Error(2)
}

func (m *MockPostRepository) GetPostsByCommunityID(communityID uint64, sortBy string, page, limit int, tags []string, userID *uint64) ([]*model.Post, int64, error) {
	args := m.Called(communityID, sortBy, page, limit, tags, userID)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Post), args.Get(1).(int64), args.Error(2)
}

func (m *MockPostRepository) GetCommunityPostsForModerator(communityID uint64, status, searchTitle string, page, limit int) ([]*model.Post, int64, error) {
	args := m.Called(communityID, status, searchTitle, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Post), args.Get(1).(int64), args.Error(2)
}

func (m *MockPostRepository) SearchPostsByTitle(title, sortBy string, page, limit int, tags []string, userID *uint64) ([]*model.Post, int64, error) {
	args := m.Called(title, sortBy, page, limit, tags, userID)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Post), args.Get(1).(int64), args.Error(2)
}

func (m *MockPostRepository) GetPostsByUserID(userID uint64, sortBy string, page, limit int) ([]*model.Post, int64, error) {
	args := m.Called(userID, sortBy, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Post), args.Get(1).(int64), args.Error(2)
}

type MockCommentRepository struct {
	mock.Mock
}

func (m *MockCommentRepository) CreateComment(comment *model.Comment) error {
	args := m.Called(comment)
	return args.Error(0)
}

func (m *MockCommentRepository) GetCommentByID(id uint64) (*model.Comment, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Comment), args.Error(1)
}

func (m *MockCommentRepository) GetCommentsByPostID(postID uint64, limit, offset int) ([]*model.Comment, int64, error) {
	args := m.Called(postID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Comment), args.Get(1).(int64), args.Error(2)
}

func (m *MockCommentRepository) UpdateComment(id uint64, content string, mediaURL *string) error {
	args := m.Called(id, content, mediaURL)
	return args.Error(0)
}

func (m *MockCommentRepository) DeleteComment(commentID uint64, parentCommentID *uint64) error {
	args := m.Called(commentID, parentCommentID)
	return args.Error(0)
}

func (m *MockCommentRepository) GetRepliesByParentID(parentID uint64) ([]*model.Comment, error) {
	args := m.Called(parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Comment), args.Error(1)
}

func (m *MockCommentRepository) GetCommentsByUserID(userID uint64, sortBy string, page, limit int) ([]*model.Comment, int64, error) {
	args := m.Called(userID, sortBy, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Comment), args.Get(1).(int64), args.Error(2)
}

type MockSubscriptionRepository struct {
	mock.Mock
}

func (m *MockSubscriptionRepository) CreateSubscription(subscription *model.Subscription) error {
	args := m.Called(subscription)
	return args.Error(0)
}

func (m *MockSubscriptionRepository) DeleteSubscription(userID, communityID uint64) error {
	args := m.Called(userID, communityID)
	return args.Error(0)
}

func (m *MockSubscriptionRepository) IsUserSubscribed(userID, communityID uint64) (bool, error) {
	args := m.Called(userID, communityID)
	return args.Bool(0), args.Error(1)
}

func (m *MockSubscriptionRepository) GetCommunityMembers(communityID uint64, sortBy, searchName, status string, page, limit int) ([]*model.Subscription, int64, error) {
	args := m.Called(communityID, sortBy, searchName, status, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Subscription), args.Get(1).(int64), args.Error(2)
}

func (m *MockSubscriptionRepository) UpdateSubscriptionStatus(userID, communityID uint64, status string) error {
	args := m.Called(userID, communityID, status)
	return args.Error(0)
}

type MockTopicRepository struct {
	mock.Mock
}

func (m *MockTopicRepository) GetAllTopics() ([]*model.Topic, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Topic), args.Error(1)
}

type MockPostReportRepository struct {
	mock.Mock
}

func (m *MockPostReportRepository) CreateReport(report *model.PostReport) error {
	args := m.Called(report)
	return args.Error(0)
}

type MockBotTaskService struct {
	mock.Mock
}

type MockNotificationService struct {
	mock.Mock
}

type MockRecommendationService struct {
	mock.Mock
}

type MockCommentVoteRepository struct {
	mock.Mock
}

type MockPostVoteRepository struct {
	mock.Mock
}

type MockTagRepository struct {
	mock.Mock
}

// MockNotificationRepository
type MockNotificationRepository struct {
	mock.Mock
}

func (m *MockNotificationRepository) CreateNotification(notification *model.Notification) error {
	args := m.Called(notification)
	return args.Error(0)
}

func (m *MockNotificationRepository) GetNotificationByID(id uint64) (*model.Notification, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Notification), args.Error(1)
}

func (m *MockNotificationRepository) GetUserNotifications(userID uint64, limit, offset int) ([]*model.Notification, int64, error) {
	args := m.Called(userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Notification), args.Get(1).(int64), args.Error(2)
}

func (m *MockNotificationRepository) MarkAsRead(id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockNotificationRepository) MarkAllAsRead(userID uint64) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockNotificationRepository) DeleteNotification(id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockNotificationRepository) GetUnreadCount(userID uint64) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

// MockMessageRepository
type MockMessageRepository struct {
	mock.Mock
}

func (m *MockMessageRepository) CreateMessage(message *model.Message) error {
	args := m.Called(message)
	return args.Error(0)
}

func (m *MockMessageRepository) GetMessageByID(id uint64) (*model.Message, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Message), args.Error(1)
}

func (m *MockMessageRepository) GetConversationMessages(conversationID uint64, page, limit int) ([]*model.Message, int64, error) {
	args := m.Called(conversationID, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Message), args.Get(1).(int64), args.Error(2)
}

func (m *MockMessageRepository) MarkMessageAsRead(messageID, userID uint64) error {
	args := m.Called(messageID, userID)
	return args.Error(0)
}

func (m *MockMessageRepository) MarkConversationMessagesAsRead(conversationID, userID uint64) error {
	args := m.Called(conversationID, userID)
	return args.Error(0)
}

func (m *MockMessageRepository) GetUnreadCount(conversationID, userID uint64) (int64, error) {
	args := m.Called(conversationID, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockMessageRepository) DeleteMessage(messageID uint64) error {
	args := m.Called(messageID)
	return args.Error(0)
}

// MockConversationRepository
type MockConversationRepository struct {
	mock.Mock
}

func (m *MockConversationRepository) CreateOrGetConversation(user1ID, user2ID uint64) (*model.Conversation, error) {
	args := m.Called(user1ID, user2ID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Conversation), args.Error(1)
}

func (m *MockConversationRepository) GetConversationByID(id uint64) (*model.Conversation, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Conversation), args.Error(1)
}

func (m *MockConversationRepository) GetConversationByUsers(user1ID, user2ID uint64) (*model.Conversation, error) {
	args := m.Called(user1ID, user2ID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Conversation), args.Error(1)
}

func (m *MockConversationRepository) GetUserConversations(userID uint64, page, limit int) ([]*model.Conversation, int64, error) {
	args := m.Called(userID, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Conversation), args.Get(1).(int64), args.Error(2)
}

func (m *MockConversationRepository) UpdateLastMessage(conversationID, messageID uint64) error {
	args := m.Called(conversationID, messageID)
	return args.Error(0)
}

func (m *MockConversationRepository) CheckUserInConversation(conversationID, userID uint64) (bool, error) {
	args := m.Called(conversationID, userID)
	return args.Bool(0), args.Error(1)
}

// MockMessageAttachmentRepository
type MockMessageAttachmentRepository struct {
	mock.Mock
}

func (m *MockMessageAttachmentRepository) CreateMessageAttachments(attachments []model.MessageAttachment) error {
	args := m.Called(attachments)
	return args.Error(0)
}

// MockSSEService
type MockSSEService struct {
	mock.Mock
}

func (m *MockSSEService) BroadcastToUser(userID uint64, event interface{}) {
	m.Called(userID, event)
}
