package constant

const (
	NOTIFICATION_ACTION_GET_POST_VOTE               = "get_post_vote"
	NOTIFICATION_ACTION_GET_POST_NEW_COMMENT        = "get_post_new_comment"
	NOTIFICATION_ACTION_GET_COMMENT_VOTE            = "get_comment_vote"
	NOTIFICATION_ACTION_GET_COMMENT_REPLY           = "get_comment_reply"
	NOTIFICATION_ACTION_POST_STATUS_UPDATED         = "post_status_updated"
	NOTIFICATION_ACTION_POST_DELETED                = "post_deleted"
	NOTIFICATION_ACTION_COMMENT_DELETED             = "comment_deleted"
	NOTIFICATION_ACTION_SUBSCRIPTION_STATUS_UPDATED = "subscription_status_updated"
	NOTIFICATION_ACTION_USER_BANNED                 = "user_banned"
	NOTIFICATION_ACTION_CONTENT_VIOLATION_POST      = "content_violation_post"
	NOTIFICATION_ACTION_CONTENT_VIOLATION_COMMENT   = "content_violation_comment"
)

var EmailSubjectMap = map[string]string{
	NOTIFICATION_ACTION_GET_POST_VOTE:               "Post Vote Received",
	NOTIFICATION_ACTION_GET_POST_NEW_COMMENT:        "New Comment on Your Post",
	NOTIFICATION_ACTION_GET_COMMENT_VOTE:            "Comment Vote Received",
	NOTIFICATION_ACTION_GET_COMMENT_REPLY:           "New Reply to Your Comment",
	NOTIFICATION_ACTION_POST_STATUS_UPDATED:         "Post Status Updated",
	NOTIFICATION_ACTION_POST_DELETED:                "Post Deleted",
	NOTIFICATION_ACTION_COMMENT_DELETED:             "Comment Deleted",
	NOTIFICATION_ACTION_SUBSCRIPTION_STATUS_UPDATED: "Subscription Status Updated",
	NOTIFICATION_ACTION_USER_BANNED:                 "User Restriction Notification",
	NOTIFICATION_ACTION_CONTENT_VIOLATION_POST:      "Content Violation - Post",
	NOTIFICATION_ACTION_CONTENT_VIOLATION_COMMENT:   "Content Violation - Comment",
}
