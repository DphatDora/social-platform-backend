package service

import (
	"fmt"
	"log"
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"
	"social-platform-backend/internal/interface/dto/response"
	"time"
)

type RecommendationService struct {
	userInterestScoreRepo repository.UserInterestScoreRepository
	userTagPrefRepo       repository.UserTagPreferenceRepository
	postRepo              repository.PostRepository
	communityRepo         repository.CommunityRepository
}

func NewRecommendationService(
	userInterestScoreRepo repository.UserInterestScoreRepository,
	userTagPrefRepo repository.UserTagPreferenceRepository,
	postRepo repository.PostRepository,
	communityRepo repository.CommunityRepository,
) *RecommendationService {
	return &RecommendationService{
		userInterestScoreRepo: userInterestScoreRepo,
		userTagPrefRepo:       userTagPrefRepo,
		postRepo:              postRepo,
		communityRepo:         communityRepo,
	}
}

func (s *RecommendationService) GetRecommendedPosts(userID uint64, page, limit int) ([]*response.PostListResponse, *response.Pagination, error) {
	// Get top communities user is interested in
	topCommunityIDs, err := s.userInterestScoreRepo.GetTopCommunitiesByScore(userID, 10)
	if err != nil {
		log.Printf("[Err] Error getting top communities in RecommendationService.GetRecommendedPosts: %v", err)
		// If no interest data, return empty list
		return []*response.PostListResponse{}, &response.Pagination{
			Total: 0,
			Page:  page,
			Limit: limit,
		}, nil
	}

	if len(topCommunityIDs) == 0 {
		// No interest data yet, return empty
		return []*response.PostListResponse{}, &response.Pagination{
			Total: 0,
			Page:  page,
			Limit: limit,
		}, nil
	}

	// Get user's preferred tags
	userTags, err := s.userTagPrefRepo.GetUserTagPreferences(userID)
	var preferredTags []string
	if err == nil && userTags != nil && userTags.PreferredTags != nil {
		preferredTags = userTags.PreferredTags
	}

	// Get posts from top communities and score them
	allPosts := []*model.Post{}
	for _, communityID := range topCommunityIDs {
		// Get recent posts from this community (last 7 days)
		posts, _, err := s.postRepo.GetPostsByCommunityID(communityID, "new", 1, 20, &userID)
		if err != nil {
			log.Printf("[Warn] Error getting posts for community %d: %v", communityID, err)
			continue
		}
		allPosts = append(allPosts, posts...)
	}

	// Score and sort posts
	scoredPosts := s.scoreAndSortPosts(allPosts, preferredTags, userID)

	start := (page - 1) * limit
	end := start + limit
	if start >= len(scoredPosts) {
		return []*response.PostListResponse{}, &response.Pagination{
			Total: int64(len(scoredPosts)),
			Page:  page,
			Limit: limit,
		}, nil
	}
	if end > len(scoredPosts) {
		end = len(scoredPosts)
	}

	paginatedPosts := scoredPosts[start:end]
	postResponses := make([]*response.PostListResponse, len(paginatedPosts))
	for i, post := range paginatedPosts {
		postResponses[i] = response.NewPostListResponse(post)
	}

	pagination := &response.Pagination{
		Total: int64(len(scoredPosts)),
		Page:  page,
		Limit: limit,
	}
	totalPages := (int64(len(scoredPosts)) + int64(limit) - 1) / int64(limit)
	if int64(page) < totalPages {
		pagination.NextURL = fmt.Sprintf("/api/v1/posts?sortBy=best&page=%d&limit=%d", page+1, limit)
	}

	return postResponses, pagination, nil
}

func (s *RecommendationService) GetRecommendedPostsByCommunity(userID, communityID uint64, page, limit int) ([]*response.PostListResponse, *response.Pagination, error) {
	// Get user's preferred tags
	userTags, err := s.userTagPrefRepo.GetUserTagPreferences(userID)
	var preferredTags []string
	if err == nil && userTags != nil && userTags.PreferredTags != nil {
		preferredTags = userTags.PreferredTags
	}

	// Get all recent posts from the community (last 30 days)
	posts, total, err := s.postRepo.GetPostsByCommunityID(communityID, "new", 1, 100, &userID)
	if err != nil {
		log.Printf("[Err] Error getting posts by community in RecommendationService.GetRecommendedPostsByCommunity: %v", err)
		return nil, nil, fmt.Errorf("failed to get posts")
	}

	// Score and sort posts
	scoredPosts := s.scoreAndSortPosts(posts, preferredTags, userID)

	start := (page - 1) * limit
	end := start + limit
	if start >= len(scoredPosts) {
		return []*response.PostListResponse{}, &response.Pagination{
			Total: total,
			Page:  page,
			Limit: limit,
		}, nil
	}
	if end > len(scoredPosts) {
		end = len(scoredPosts)
	}

	paginatedPosts := scoredPosts[start:end]
	postResponses := make([]*response.PostListResponse, len(paginatedPosts))
	for i, post := range paginatedPosts {
		postResponses[i] = response.NewPostListResponse(post)
	}

	pagination := &response.Pagination{
		Total: int64(len(scoredPosts)),
		Page:  page,
		Limit: limit,
	}
	totalPages := (int64(len(scoredPosts)) + int64(limit) - 1) / int64(limit)
	if int64(page) < totalPages {
		pagination.NextURL = fmt.Sprintf("/api/v1/communities/%d/posts?sortBy=best&page=%d&limit=%d", communityID, page+1, limit)
	}

	return postResponses, pagination, nil
}

func (s *RecommendationService) scoreAndSortPosts(posts []*model.Post, preferredTags []string, userID uint64) []*model.Post {
	type scoredPost struct {
		post  *model.Post
		score float64
	}

	scored := make([]scoredPost, 0, len(posts))

	for _, post := range posts {
		score := s.calculatePostScore(post, preferredTags)
		scored = append(scored, scoredPost{post: post, score: score})
	}

	// Sort by score descending
	for i := 0; i < len(scored)-1; i++ {
		for j := i + 1; j < len(scored); j++ {
			if scored[j].score > scored[i].score {
				scored[i], scored[j] = scored[j], scored[i]
			}
		}
	}

	// Extract sorted posts
	result := make([]*model.Post, len(scored))
	for i, sp := range scored {
		result[i] = sp.post
	}

	return result
}

func (s *RecommendationService) calculatePostScore(post *model.Post, preferredTags []string) float64 {
	score := 0.0

	// Factor 1: Post votes (engagement)
	score += float64(post.Vote) * 0.5

	// Factor 2: Tag matching
	if post.Tags != nil && len(*post.Tags) > 0 {
		tagMatchCount := 0
		for _, postTag := range *post.Tags {
			for _, prefTag := range preferredTags {
				if postTag == prefTag {
					tagMatchCount++
					break
				}
			}
		}
		score += float64(tagMatchCount) * 10.0
	}

	// Factor 3: Freshness (decay over time)
	// Posts from last 24 hours get bonus, then decay
	hoursSincePost := time.Since(post.CreatedAt).Hours()
	if hoursSincePost < 24 {
		score += 20.0 * (1.0 - hoursSincePost/24.0)
	} else if hoursSincePost < 168 { // 7 days
		score += 10.0 * (1.0 - (hoursSincePost-24.0)/144.0)
	}

	// Factor 4: Author karma (quality indicator)
	if post.Author != nil {
		score += float64(post.Author.Karma) / 100.0
	}

	return score
}
