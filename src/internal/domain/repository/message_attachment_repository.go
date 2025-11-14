package repository

import "social-platform-backend/internal/domain/model"

type MessageAttachmentRepository interface {
	CreateMessageAttachments(attachments []model.MessageAttachment) error
	GetAttachmentsByMessageID(messageID uint64) ([]model.MessageAttachment, error)
}
