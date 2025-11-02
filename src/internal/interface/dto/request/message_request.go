package request

type SendMessageRequest struct {
	RecipientID uint64   `json:"recipientId" binding:"required"`
	Content     string   `json:"content" binding:"required"`
	Type        string   `json:"type"` // "text", "image", "video", "file"
	Attachments []string `json:"attachments,omitempty"`
}
