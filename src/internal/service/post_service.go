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

type PostService struct {
	postRepo      repository.PostRepository
	communityRepo repository.CommunityRepository
	postVoteRepo  repository.PostVoteRepository
}

func NewPostService(
	postRepo repository.PostRepository,
	communityRepo repository.CommunityRepository,
	postVoteRepo repository.PostVoteRepository,
) *PostService {
	return &PostService{
		postRepo:      postRepo,
		communityRepo: communityRepo,
		postVoteRepo:  postVoteRepo,
	}
}

func (s *PostService) CreatePost(userID uint64, req *request.CreatePostRequest) error {
	// Check if community exists
	_, err := s.communityRepo.GetCommunityByID(req.CommunityID)
	if err != nil {
		log.Printf("[Err] Community not found in PostService.CreatePost: %v", err)
		return fmt.Errorf("community not found")
	}

	// Validate post type and required fields
	switch req.Type {
	case constant.PostTypeText:
		// text post: only title and content required
	case constant.PostTypeLink:
		if req.URL == nil {
			log.Printf("[Err] URL is required for link post in PostService.CreatePost")
			return fmt.Errorf("url is required for link post")
		}
	case constant.PostTypeMedia:
		if req.MediaURLs == nil || len(*req.MediaURLs) == 0 {
			log.Printf("[Err] Media URLs are required for media post in PostService.CreatePost")
			return fmt.Errorf("media_urls are required for media post")
		}
	case constant.PostTypePoll:
		if req.PollData == nil {
			log.Printf("[Err] Poll data is required for poll post in PostService.CreatePost")
			return fmt.Errorf("poll_data is required for poll post")
		}
	default:
		log.Printf("[Err] Invalid post type in PostService.CreatePost: %s", req.Type)
		return fmt.Errorf("invalid post type")
	}

	post := &model.Post{
		CommunityID: req.CommunityID,
		AuthorID:    userID,
		Title:       req.Title,
		Type:        req.Type,
		Content:     req.Content,
		URL:         req.URL,
		MediaURLs:   req.MediaURLs,
		PollData:    req.PollData,
		Tags:        req.Tags,
	}

	if err := s.postRepo.CreatePost(post); err != nil {
		log.Printf("[Err] Error creating post in PostService.CreatePost: %v", err)
		return fmt.Errorf("failed to create post")
	}

	return nil
}

func (s *PostService) UpdatePost(userID, postID uint64, postType string, reqBody interface{}) error {
	post, err := s.postRepo.GetPostByID(postID)
	if err != nil {
		log.Printf("[Err] Post not found in PostService.UpdatePost: %v", err)
		return fmt.Errorf("post not found")
	}

	// Check if user is the author
	if post.AuthorID != userID {
		log.Printf("[Err] User does not have permission to update post in PostService.UpdatePost: userID=%d, postID=%d", userID, postID)
		return fmt.Errorf("permission denied")
	}

	// Check if post type matches
	if post.Type != postType {
		log.Printf("[Err] Post type mismatch in PostService.UpdatePost: expected=%s, actual=%s", postType, post.Type)
		return fmt.Errorf("post type mismatch")
	}

	// Update based on post type
	switch postType {
	case constant.PostTypeText:
		req, ok := reqBody.(*request.UpdatePostTextRequest)
		if !ok {
			return fmt.Errorf("invalid request body for text post")
		}
		if err := s.postRepo.UpdatePostText(postID, req); err != nil {
			log.Printf("[Err] Error updating text post in PostService.UpdatePost: %v", err)
			return fmt.Errorf("failed to update post")
		}
	case constant.PostTypeLink:
		req, ok := reqBody.(*request.UpdatePostLinkRequest)
		if !ok {
			return fmt.Errorf("invalid request body for link post")
		}
		if err := s.postRepo.UpdatePostLink(postID, req); err != nil {
			log.Printf("[Err] Error updating link post in PostService.UpdatePost: %v", err)
			return fmt.Errorf("failed to update post")
		}
	case constant.PostTypeMedia:
		req, ok := reqBody.(*request.UpdatePostMediaRequest)
		if !ok {
			return fmt.Errorf("invalid request body for media post")
		}
		if err := s.postRepo.UpdatePostMedia(postID, req); err != nil {
			log.Printf("[Err] Error updating media post in PostService.UpdatePost: %v", err)
			return fmt.Errorf("failed to update post")
		}
	case constant.PostTypePoll:
		req, ok := reqBody.(*request.UpdatePostPollRequest)
		if !ok {
			return fmt.Errorf("invalid request body for poll post")
		}
		if err := s.postRepo.UpdatePostPoll(postID, req); err != nil {
			log.Printf("[Err] Error updating poll post in PostService.UpdatePost: %v", err)
			return fmt.Errorf("failed to update post")
		}
	default:
		return fmt.Errorf("invalid post type")
	}

	return nil
}

func (s *PostService) DeletePost(userID, postID uint64) error {
	post, err := s.postRepo.GetPostByID(postID)
	if err != nil {
		log.Printf("[Err] Post not found in PostService.DeletePost: %v", err)
		return fmt.Errorf("post not found")
	}

	// Check if user is the author
	if post.AuthorID != userID {
		log.Printf("[Err] User does not have permission to delete post in PostService.DeletePost: userID=%d, postID=%d", userID, postID)
		return fmt.Errorf("permission denied")
	}

	if err := s.postRepo.DeletePost(postID); err != nil {
		log.Printf("[Err] Error deleting post in PostService.DeletePost: %v", err)
		return fmt.Errorf("failed to delete post")
	}

	return nil
}

func (s *PostService) GetAllPosts(sortBy string, page, limit int) ([]*response.PostListResponse, *response.Pagination, error) {
	posts, total, err := s.postRepo.GetAllPosts(sortBy, page, limit)
	if err != nil {
		log.Printf("[Err] Error getting all posts in PostService.GetAllPosts: %v", err)
		return nil, nil, fmt.Errorf("failed to get posts")
	}

	postResponses := make([]*response.PostListResponse, len(posts))
	for i, post := range posts {
		postResponses[i] = response.NewPostListResponse(post)
	}

	// Set pagination
	pagination := &response.Pagination{
		Total: total,
		Page:  page,
		Limit: limit,
	}
	totalPages := (total + int64(limit) - 1) / int64(limit)
	if int64(page) < totalPages {
		pagination.NextURL = fmt.Sprintf("/api/v1/posts?sortBy=%s&page=%d&limit=%d", sortBy, page+1, limit)
	}

	return postResponses, pagination, nil
}

func (s *PostService) GetPostsByCommunityID(communityID uint64, sortBy string, page, limit int) ([]*response.PostListResponse, *response.Pagination, error) {
	// Check if community exists
	_, err := s.communityRepo.GetCommunityByID(communityID)
	if err != nil {
		log.Printf("[Err] Community not found in PostService.GetPostsByCommunityID: %v", err)
		return nil, nil, fmt.Errorf("community not found")
	}

	posts, total, err := s.postRepo.GetPostsByCommunityID(communityID, sortBy, page, limit)
	if err != nil {
		log.Printf("[Err] Error getting posts by community ID in PostService.GetPostsByCommunityID: %v", err)
		return nil, nil, fmt.Errorf("failed to get posts")
	}

	postResponses := make([]*response.PostListResponse, len(posts))
	for i, post := range posts {
		postResponses[i] = response.NewPostListResponse(post)
	}

	pagination := &response.Pagination{
		Total: total,
		Page:  page,
		Limit: limit,
	}
	totalPages := (total + int64(limit) - 1) / int64(limit)
	if int64(page) < totalPages {
		pagination.NextURL = fmt.Sprintf("/api/v1/communities/%d/posts?sortBy=%s&page=%d&limit=%d", communityID, sortBy, page+1, limit)
	}

	return postResponses, pagination, nil
}

func (s *PostService) SearchPostsByTitle(title, sortBy string, page, limit int) ([]*response.PostListResponse, *response.Pagination, error) {
	posts, total, err := s.postRepo.SearchPostsByTitle(title, sortBy, page, limit)
	if err != nil {
		log.Printf("[Err] Error searching posts by title in PostService.SearchPostsByTitle: %v", err)
		return nil, nil, fmt.Errorf("failed to search posts")
	}

	postResponses := make([]*response.PostListResponse, len(posts))
	for i, post := range posts {
		postResponses[i] = response.NewPostListResponse(post)
	}

	// Set pagination
	pagination := &response.Pagination{
		Total: total,
		Page:  page,
		Limit: limit,
	}
	totalPages := (total + int64(limit) - 1) / int64(limit)
	if int64(page) < totalPages {
		pagination.NextURL = fmt.Sprintf("/api/v1/posts/search?search=%s&sortBy=%s&page=%d&limit=%d", title, sortBy, page+1, limit)
	}

	return postResponses, pagination, nil
}

func (s *PostService) VotePost(userID, postID uint64, vote bool) error {
	// Check if post exists
	_, err := s.postRepo.GetPostByID(postID)
	if err != nil {
		log.Printf("[Err] Post not found in PostService.VotePost: %v", err)
		return fmt.Errorf("post not found")
	}

	postVote := &model.PostVote{
		UserID: userID,
		PostID: postID,
		Vote:   vote,
	}

	if err := s.postVoteRepo.UpsertPostVote(postVote); err != nil {
		log.Printf("[Err] Error upserting post vote in PostService.VotePost: %v", err)
		return fmt.Errorf("failed to vote post")
	}

	return nil
}

func (s *PostService) UnvotePost(userID, postID uint64) error {
	// Check if post exists
	_, err := s.postRepo.GetPostByID(postID)
	if err != nil {
		log.Printf("[Err] Post not found in PostService.UnvotePost: %v", err)
		return fmt.Errorf("post not found")
	}

	// Delete vote
	if err := s.postVoteRepo.DeletePostVote(userID, postID); err != nil {
		log.Printf("[Err] Error deleting post vote in PostService.UnvotePost: %v", err)
		return fmt.Errorf("failed to unvote post")
	}

	return nil
}
