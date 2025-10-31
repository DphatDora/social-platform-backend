package repository

import "social-platform-backend/internal/domain/model"

type MessageRepository interface {
	CreateMessage(message *model.Message) error
	GetMessageByID(id uint64) (*model.Message, error)
	GetConversationMessages(conversationID uint64, page, limit int) ([]*model.Message, int64, error)
	MarkMessageAsRead(messageID, userID uint64) error
	MarkConversationMessagesAsRead(conversationID, userID uint64) error
	GetUnreadCount(conversationID, userID uint64) (int64, error)
}
