package request

type MessageAttachmentRequest struct {
	FileURL  string `json:"fileUrl" binding:"required"`
	FileType string `json:"fileType" binding:"required"`
}

type SendMessageRequest struct {
	RecipientID uint64                     `json:"recipientId" binding:"required"`
	Content     string                     `json:"content" binding:"required"`
	Attachments []MessageAttachmentRequest `json:"attachments,omitempty"`
}
