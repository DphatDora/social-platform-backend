package repository

import (
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"
	"time"

	"gorm.io/gorm"
)

type ConversationRepositoryImpl struct {
	db *gorm.DB
}

func NewConversationRepository(db *gorm.DB) repository.ConversationRepository {
	return &ConversationRepositoryImpl{db: db}
}

func (r *ConversationRepositoryImpl) CreateOrGetConversation(user1ID, user2ID uint64) (*model.Conversation, error) {
	// Ensure user1ID < user2ID for consistency
	if user1ID > user2ID {
		user1ID, user2ID = user2ID, user1ID
	}

	var conversation model.Conversation
	err := r.db.Where("(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)",
		user1ID, user2ID, user2ID, user1ID).
		First(&conversation).Error

	if err == gorm.ErrRecordNotFound {
		// Create new conversation
		conversation = model.Conversation{
			User1ID:   user1ID,
			User2ID:   user2ID,
			CreatedAt: time.Now(),
		}
		if err := r.db.Create(&conversation).Error; err != nil {
			return nil, err
		}
		return &conversation, nil
	}

	if err != nil {
		return nil, err
	}

	return &conversation, nil
}

func (r *ConversationRepositoryImpl) GetConversationByID(id uint64) (*model.Conversation, error) {
	var conversation model.Conversation
	err := r.db.Preload("User1").
		Preload("User2").
		Preload("LastMessage").
		Where("id = ?", id).
		First(&conversation).Error
	if err != nil {
		return nil, err
	}
	return &conversation, nil
}

func (r *ConversationRepositoryImpl) GetConversationByUsers(user1ID, user2ID uint64) (*model.Conversation, error) {
	var conversation model.Conversation
	err := r.db.Where("(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)",
		user1ID, user2ID, user2ID, user1ID).
		First(&conversation).Error
	if err != nil {
		return nil, err
	}
	return &conversation, nil
}

func (r *ConversationRepositoryImpl) GetUserConversations(userID uint64, page, limit int) ([]*model.Conversation, int64, error) {
	var conversations []*model.Conversation
	var total int64

	offset := (page - 1) * limit

	// Count total conversations
	query := r.db.Model(&model.Conversation{}).
		Where("user1_id = ? OR user2_id = ?", userID, userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get conversations with preloads, sorted by last_message_at DESC
	err := r.db.Preload("User1").
		Preload("User2").
		Preload("LastMessage").
		Where("user1_id = ? OR user2_id = ?", userID, userID).
		Order("COALESCE(last_message_at, created_at) DESC").
		Limit(limit).
		Offset(offset).
		Find(&conversations).Error

	if err != nil {
		return nil, 0, err
	}

	return conversations, total, nil
}

func (r *ConversationRepositoryImpl) UpdateLastMessage(conversationID, messageID uint64) error {
	return r.db.Model(&model.Conversation{}).
		Where("id = ?", conversationID).
		Updates(map[string]interface{}{
			"last_message_id": messageID,
			"last_message_at": time.Now(),
			"updated_at":      time.Now(),
		}).Error
}

func (r *ConversationRepositoryImpl) CheckUserInConversation(conversationID, userID uint64) (bool, error) {
	var count int64
	err := r.db.Model(&model.Conversation{}).
		Where("id = ? AND (user1_id = ? OR user2_id = ?)", conversationID, userID, userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
