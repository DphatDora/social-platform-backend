package service

import (
	"fmt"
	"log"
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"
	"social-platform-backend/internal/interface/dto/request"
	"social-platform-backend/internal/interface/dto/response"
	"social-platform-backend/package/constant"
)

type CommentService struct {
	commentRepo repository.CommentRepository
	postRepo    repository.PostRepository
}

func NewCommentService(
	commentRepo repository.CommentRepository,
	postRepo repository.PostRepository,
) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		postRepo:    postRepo,
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
