package response

import (
	"social-platform-backend/internal/domain/model"
	"time"
)

type LastMessagePreview struct {
	Content string     `json:"content"`
	IsRead  bool       `json:"isRead"`
	SentAt  time.Time  `json:"sentAt"`
	ReadAt  *time.Time `json:"readAt,omitempty"`
}

type OtherUserInfo struct {
	ID       uint64  `json:"id"`
	Username string  `json:"username"`
	Avatar   *string `json:"avatar,omitempty"`
}

type ConversationListResponse struct {
	ID          uint64              `json:"id"`
	OtherUser   OtherUserInfo       `json:"otherUser"`
	LastMessage *LastMessagePreview `json:"lastMessage,omitempty"`
	UnreadCount int64               `json:"unreadCount"`
	CreatedAt   time.Time           `json:"createdAt"`
	UpdatedAt   *time.Time          `json:"updatedAt,omitempty"`
}

func NewConversationListResponse(conversation *model.Conversation, currentUserID uint64, unreadCount int64) *ConversationListResponse {
	resp := &ConversationListResponse{
		ID:          conversation.ID,
		UnreadCount: unreadCount,
		CreatedAt:   conversation.CreatedAt,
		UpdatedAt:   conversation.UpdatedAt,
	}

	var otherUser *model.User
	if conversation.User1ID == currentUserID {
		otherUser = conversation.User2
	} else {
		otherUser = conversation.User1
	}

	if otherUser != nil {
		resp.OtherUser = OtherUserInfo{
			ID:       otherUser.ID,
			Username: otherUser.Username,
			Avatar:   otherUser.Avatar,
		}
	}

	if conversation.LastMessage != nil {
		resp.LastMessage = &LastMessagePreview{
			Content: conversation.LastMessage.Content,
			IsRead:  conversation.LastMessage.IsRead,
			SentAt:  conversation.LastMessage.CreatedAt,
			ReadAt:  conversation.LastMessage.ReadAt,
		}
	}

	return resp
}

type MessageAttachmentResponse struct {
	ID       uint64 `json:"id"`
	FileURL  string `json:"fileUrl"`
	FileType string `json:"fileType"`
}

type SenderInfo struct {
	ID       uint64  `json:"id"`
	Username string  `json:"username"`
	Avatar   *string `json:"avatar,omitempty"`
}

type MessageResponse struct {
	ID             uint64                      `json:"id"`
	ConversationID uint64                      `json:"conversationId"`
	Sender         SenderInfo                  `json:"sender"`
	Type           string                      `json:"type"`
	Content        string                      `json:"content"`
	IsRead         bool                        `json:"isRead"`
	ReadAt         *time.Time                  `json:"readAt,omitempty"`
	Attachments    []MessageAttachmentResponse `json:"attachments,omitempty"`
	CreatedAt      time.Time                   `json:"createdAt"`
}

func NewMessageResponse(message *model.Message) *MessageResponse {
	resp := &MessageResponse{
		ID:             message.ID,
		ConversationID: message.ConversationID,
		Type:           message.Type,
		Content:        message.Content,
		IsRead:         message.IsRead,
		ReadAt:         message.ReadAt,
		CreatedAt:      message.CreatedAt,
	}

	if message.Sender != nil {
		resp.Sender = SenderInfo{
			ID:       message.Sender.ID,
			Username: message.Sender.Username,
			Avatar:   message.Sender.Avatar,
		}
	}

	if len(message.Attachments) > 0 {
		resp.Attachments = make([]MessageAttachmentResponse, len(message.Attachments))
		for i, att := range message.Attachments {
			resp.Attachments[i] = MessageAttachmentResponse{
				ID:       att.ID,
				FileURL:  att.FileURL,
				FileType: att.FileType,
			}
		}
	}

	return resp
}

// SSEEvent represents a Server-Sent Event
type SSEEvent struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

// NewMessageEvent represents a new message event for SSE
type NewMessageEvent struct {
	ConversationID uint64          `json:"conversationId"`
	Message        MessageResponse `json:"message"`
}

// ConversationUpdatedEvent represents a conversation update event for SSE
type ConversationUpdatedEvent struct {
	Conversation ConversationListResponse `json:"conversation"`
}
