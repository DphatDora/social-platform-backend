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
