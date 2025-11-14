package repository

import (
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"

	"gorm.io/gorm"
)

type MessageAttachmentRepositoryImpl struct {
	db *gorm.DB
}

func NewMessageAttachmentRepository(db *gorm.DB) repository.MessageAttachmentRepository {
	return &MessageAttachmentRepositoryImpl{db: db}
}

func (r *MessageAttachmentRepositoryImpl) CreateMessageAttachments(attachments []model.MessageAttachment) error {
	if len(attachments) == 0 {
		return nil
	}
	return r.db.Create(&attachments).Error
}

func (r *MessageAttachmentRepositoryImpl) GetAttachmentsByMessageID(messageID uint64) ([]model.MessageAttachment, error) {
	var attachments []model.MessageAttachment
	err := r.db.Where("message_id = ?", messageID).Find(&attachments).Error
	return attachments, err
}
