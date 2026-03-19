package request

type MessageAttachmentRequest struct {
	FileURL  string `json:"fileUrl" binding:"required"`
	FileType string `json:"fileType" binding:"required"`
}

type MetaDataRequest struct {
	ID        uint64   `json:"id"`
	Title     string   `json:"title"`
	Tags      []string `json:"tags"`
	Content   string   `json:"content"`
	MediaURLs []string `json:"mediaUrls"`
}

type SendMessageRequest struct {
	RecipientID uint64                     `json:"recipientId" binding:"required"`
	Content     string                     `json:"content" binding:"required"`
	Attachments []MessageAttachmentRequest `json:"attachments,omitempty"`
	MetaData    *MetaDataRequest           `json:"metadata,omitempty"`
}
