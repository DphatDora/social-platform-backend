package payload

type PostVoteNotificationPayload struct {
	PostID   uint64 `json:"postId"`
	UserName string `json:"userName"`
	VoteType bool   `json:"voteType"` // true for upvote, false for downvote
}

type PostCommentNotificationPayload struct {
	PostID   uint64 `json:"postId"`
	UserName string `json:"userName"`
}

type CommentVoteNotificationPayload struct {
	CommentID uint64 `json:"commentId"`
	UserName  string `json:"userName"`
	VoteType  bool   `json:"voteType"` // true for upvote, false for downvote
}

type CommentReplyNotificationPayload struct {
	CommentID uint64 `json:"commentId"`
	UserName  string `json:"userName"`
}

type PostStatusNotificationPayload struct {
	PostID uint64 `json:"postId"`
	Status string `json:"status"`
}

type SubscriptionStatusNotificationPayload struct {
	CommunityID   uint64 `json:"communityId"`
	CommunityName string `json:"communityName"`
	Status        string `json:"status"`
}

type UserBanNotificationPayload struct {
	CommunityID     uint64 `json:"communityId"`
	CommunityName   string `json:"communityName"`
	RestrictionType string `json:"restrictionType"`
	Reason          string `json:"reason"`
	ExpiresAt       string `json:"expiresAt,omitempty"`
}

type CommentDeletedNotificationPayload struct {
	CommentID uint64 `json:"commentId"`
	PostID    uint64 `json:"postId"`
}
