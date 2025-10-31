package repository

import "social-platform-backend/internal/domain/model"

type ConversationRepository interface {
	CreateOrGetConversation(user1ID, user2ID uint64) (*model.Conversation, error) // create new conversation or get existing one
	GetConversationByID(id uint64) (*model.Conversation, error)
	GetConversationByUsers(user1ID, user2ID uint64) (*model.Conversation, error)
	GetUserConversations(userID uint64, page, limit int) ([]*model.Conversation, int64, error)
	UpdateLastMessage(conversationID, messageID uint64) error
	CheckUserInConversation(conversationID, userID uint64) (bool, error) // checks if user is part of the conversation
}
