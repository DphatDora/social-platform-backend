package repository

import (
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"
	"time"

	"gorm.io/gorm"
)

type MessageRepositoryImpl struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) repository.MessageRepository {
	return &MessageRepositoryImpl{db: db}
}

func (r *MessageRepositoryImpl) CreateMessage(message *model.Message) error {
	return r.db.Create(message).Error
}

func (r *MessageRepositoryImpl) GetMessageByID(id uint64) (*model.Message, error) {
	var message model.Message
	err := r.db.Unscoped().
		Preload("Sender").
		Preload("Attachments").
		Where("id = ?", id).
		First(&message).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *MessageRepositoryImpl) GetConversationMessages(conversationID uint64, page, limit int) ([]*model.Message, int64, error) {
	var messages []*model.Message
	var total int64

	offset := (page - 1) * limit

	// Count total messages (including soft deleted)
	query := r.db.Unscoped().Model(&model.Message{}).
		Where("conversation_id = ?", conversationID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get messages with preloads, sorted by created_at DESC (newest first)
	// Use Unscoped() to include soft deleted messages
	err := r.db.Unscoped().
		Preload("Sender").
		Preload("Attachments").
		Where("conversation_id = ?", conversationID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error

	if err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

func (r *MessageRepositoryImpl) MarkMessageAsRead(messageID, userID uint64) error {
	// Only mark as read if the user is not the sender
	return r.db.Model(&model.Message{}).
		Where("id = ? AND sender_id != ? AND is_read = false", messageID, userID).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": time.Now(),
		}).Error
}

func (r *MessageRepositoryImpl) MarkConversationMessagesAsRead(conversationID, userID uint64) error {
	// Mark all unread messages in conversation as read (except messages sent by the user)
	return r.db.Model(&model.Message{}).
		Where("conversation_id = ? AND sender_id != ? AND is_read = false", conversationID, userID).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": time.Now(),
		}).Error
}

func (r *MessageRepositoryImpl) GetUnreadCount(conversationID, userID uint64) (int64, error) {
	var count int64
	err := r.db.Model(&model.Message{}).
		Where("conversation_id = ? AND sender_id != ? AND is_read = false", conversationID, userID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *MessageRepositoryImpl) DeleteMessage(messageID uint64) error {
	return r.db.Delete(&model.Message{}, messageID).Error
}
