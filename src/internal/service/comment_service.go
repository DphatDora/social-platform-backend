package service

import (
	"encoding/json"
	"fmt"
	"log"
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"
	"social-platform-backend/internal/interface/dto/request"
	"social-platform-backend/internal/interface/dto/response"
	"social-platform-backend/package/constant"
	"social-platform-backend/package/template/payload"
	"time"
)

type CommentService struct {
	commentRepo     repository.CommentRepository
	postRepo        repository.PostRepository
	commentVoteRepo repository.CommentVoteRepository
	botTaskRepo     repository.BotTaskRepository
}

func NewCommentService(
	commentRepo repository.CommentRepository,
	postRepo repository.PostRepository,
	commentVoteRepo repository.CommentVoteRepository,
	botTaskRepo repository.BotTaskRepository,
) *CommentService {
	return &CommentService{
		commentRepo:     commentRepo,
		postRepo:        postRepo,
		commentVoteRepo: commentVoteRepo,
		botTaskRepo:     botTaskRepo,
	}
}

func (s *CommentService) CreateComment(userID uint64, req *request.CreateCommentRequest) error {
	// Check if post exists
	_, err := s.postRepo.GetPostByID(req.PostID)
	if err != nil {
		log.Printf("[Err] Post not found in CommentService.CreateComment: %v", err)
		return fmt.Errorf("post not found")
	}

	// If it is a reply, check if parent comment exists and belongs to the same post
	if req.ParentCommentID != nil {
		parentComment, err := s.commentRepo.GetCommentByID(*req.ParentCommentID)
		if err != nil {
			log.Printf("[Err] Parent comment not found in CommentService.CreateComment: %v", err)
			return fmt.Errorf("parent comment not found")
		}
		if parentComment.PostID != req.PostID {
			log.Printf("[Err] Parent comment does not belong to the same post in CommentService.CreateComment")
			return fmt.Errorf("parent comment does not belong to this post")
		}
	}

	comment := &model.Comment{
		PostID:          req.PostID,
		AuthorID:        userID,
		ParentCommentID: req.ParentCommentID,
		Content:         req.Content,
		MediaURL:        req.MediaURL,
	}

	if err := s.commentRepo.CreateComment(comment); err != nil {
		log.Printf("[Err] Error creating comment in CommentService.CreateComment: %v", err)
		return fmt.Errorf("failed to create comment")
	}

	karmaPayload := payload.UpdateUserKarmaPayload{
		UserId:    userID,
		TargetId:  nil,
		Action:    constant.KARMA_ACTION_CREATE_COMMENT,
		UpdatedAt: time.Now(),
	}

	payloadBytes, err := json.Marshal(karmaPayload)
	if err != nil {
		log.Printf("[Err] Error marshaling karma payload in CommentService.CreateComment: %v", err)
		return nil
	}

	rawPayload := json.RawMessage(payloadBytes)
	now := time.Now()
	botTask := &model.BotTask{
		Action:     constant.BOT_TASK_ACTION_UPDATE_KARMA,
		Payload:    &rawPayload,
		CreatedAt:  now,
		ExecutedAt: &now,
	}

	if err := s.botTaskRepo.CreateBotTask(botTask); err != nil {
		log.Printf("[Err] Error creating bot task in CommentService.CreateComment: %v", err)
	}

	return nil
}

func (s *CommentService) GetCommentsByPostID(postID uint64, page, limit int) ([]*response.CommentResponse, *response.Pagination, error) {
	// Validate pagination
	if page <= 0 {
		page = constant.DEFAULT_PAGE
	}
	if limit <= 0 || limit > 100 {
		limit = constant.DEFAULT_LIMIT
	}

	offset := (page - 1) * limit

	// Get top-level comments
	comments, total, err := s.commentRepo.GetCommentsByPostID(postID, limit, offset)
	if err != nil {
		log.Printf("[Err] Error getting comments in CommentService.GetCommentsByPostID: %v", err)
		return nil, nil, fmt.Errorf("failed to get comments")
	}

	// Build response with nested replies
	var commentResponses []*response.CommentResponse
	for _, comment := range comments {
		commentResp := response.NewCommentResponse(comment)

		// Get all replies for this comment
		replies, err := s.loadRepliesRecursively(comment.ID)
		if err != nil {
			log.Printf("[Err] Error loading replies in CommentService.GetCommentsByPostID: %v", err)
			// Continue without replies rather than failing entirely
			replies = []*response.CommentResponse{}
		}
		commentResp.Replies = replies

		commentResponses = append(commentResponses, commentResp)
	}

	pagination := &response.Pagination{
		Total: total,
		Page:  page,
		Limit: limit,
	}
	totalPages := (total + int64(limit) - 1) / int64(limit)
	if int64(page) < totalPages {
		pagination.NextURL = fmt.Sprintf("/api/v1/posts/%d/comments?page=%d&limit=%d", postID, page+1, limit)
	}

	return commentResponses, pagination, nil
}

func (s *CommentService) loadRepliesRecursively(parentID uint64) ([]*response.CommentResponse, error) {
	replies, err := s.commentRepo.GetRepliesByParentID(parentID)
	if err != nil {
		return nil, err
	}

	var replyResponses []*response.CommentResponse
	for _, reply := range replies {
		replyResp := response.NewCommentResponse(reply)

		// Load child replies
		nestedReplies, err := s.loadRepliesRecursively(reply.ID)
		if err != nil {
			log.Printf("[Err] Error loading nested replies: %v", err)
			nestedReplies = []*response.CommentResponse{}
		}
		replyResp.Replies = nestedReplies

		replyResponses = append(replyResponses, replyResp)
	}

	return replyResponses, nil
}

func (s *CommentService) UpdateComment(userID, commentID uint64, req *request.UpdateCommentRequest) error {
	comment, err := s.commentRepo.GetCommentByID(commentID)
	if err != nil {
		log.Printf("[Err] Comment not found in CommentService.UpdateComment: %v", err)
		return fmt.Errorf("comment not found")
	}

	// Check if user is the author
	if comment.AuthorID != userID {
		log.Printf("[Err] User does not have permission to update comment in CommentService.UpdateComment: userID=%d, commentID=%d", userID, commentID)
		return fmt.Errorf("permission denied")
	}

	if err := s.commentRepo.UpdateComment(commentID, req.Content, req.MediaURL); err != nil {
		log.Printf("[Err] Error updating comment in CommentService.UpdateComment: %v", err)
		return fmt.Errorf("failed to update comment")
	}

	return nil
}

func (s *CommentService) DeleteComment(userID, commentID uint64) error {
	comment, err := s.commentRepo.GetCommentByID(commentID)
	if err != nil {
		log.Printf("[Err] Comment not found in CommentService.DeleteComment: %v", err)
		return fmt.Errorf("comment not found")
	}

	// Check if user is the author
	if comment.AuthorID != userID {
		log.Printf("[Err] User does not have permission to delete comment in CommentService.DeleteComment: userID=%d, commentID=%d", userID, commentID)
		return fmt.Errorf("permission denied")
	}

	// Delete comment with transaction (updates replies' parent_comment_id, then deletes the comment)
	if err := s.commentRepo.DeleteComment(commentID, comment.ParentCommentID); err != nil {
		log.Printf("[Err] Error deleting comment in CommentService.DeleteComment: %v", err)
		return fmt.Errorf("failed to delete comment")
	}

	return nil
}

func (s *CommentService) VoteComment(userID, commentID uint64, vote bool) error {
	// Check if comment exists
	comment, err := s.commentRepo.GetCommentByID(commentID)
	if err != nil {
		log.Printf("[Err] Comment not found in CommentService.VoteComment: %v", err)
		return fmt.Errorf("comment not found")
	}

	commentVote := &model.CommentVote{
		UserID:    userID,
		CommentID: commentID,
		Vote:      vote,
	}

	if err := s.commentVoteRepo.UpsertCommentVote(commentVote); err != nil {
		log.Printf("[Err] Error voting comment in CommentService.VoteComment: %v", err)
		return fmt.Errorf("failed to vote comment")
	}

	// Create bot task for updating karma
	var action string
	if vote {
		action = constant.KARMA_ACTION_UPVOTE_COMMENT
	} else {
		action = constant.KARMA_ACTION_DOWNVOTE_COMMENT
	}

	karmaPayload := payload.UpdateUserKarmaPayload{
		UserId:    userID,
		TargetId:  &comment.AuthorID,
		Action:    action,
		UpdatedAt: time.Now(),
	}

	payloadBytes, err := json.Marshal(karmaPayload)
	if err != nil {
		log.Printf("[Err] Error marshaling karma payload in CommentService.VoteComment: %v", err)
		return nil
	}

	rawPayload := json.RawMessage(payloadBytes)
	now := time.Now()
	botTask := &model.BotTask{
		Action:     constant.BOT_TASK_ACTION_UPDATE_KARMA,
		Payload:    &rawPayload,
		CreatedAt:  now,
		ExecutedAt: &now,
	}

	if err := s.botTaskRepo.CreateBotTask(botTask); err != nil {
		log.Printf("[Err] Error creating bot task in CommentService.VoteComment: %v", err)
	}

	return nil
}

func (s *CommentService) UnvoteComment(userID, commentID uint64) error {
	// Check if comment exists
	_, err := s.commentRepo.GetCommentByID(commentID)
	if err != nil {
		log.Printf("[Err] Comment not found in CommentService.UnvoteComment: %v", err)
		return fmt.Errorf("comment not found")
	}

	// Delete vote
	if err := s.commentVoteRepo.DeleteCommentVote(userID, commentID); err != nil {
		log.Printf("[Err] Error unvoting comment in CommentService.UnvoteComment: %v", err)
		return fmt.Errorf("failed to unvote comment")
	}

	return nil
}
