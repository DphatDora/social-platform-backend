package service

import (
	"fmt"
	"log"
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"
	"social-platform-backend/internal/interface/dto/response"
	"social-platform-backend/package/constant"
	"sort"
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
	// Get top communities user is interested in (limited to MAX_COMMUNITIES)
	communityScores, err := s.userInterestScoreRepo.GetUserInterestScoresWithScores(userID, constant.RECOMMENDATION_MAX_COMMUNITIES)
	if err != nil {
		log.Printf("[Err] Error getting community scores in RecommendationService.GetRecommendedPosts: %v", err)
		// If no interest data, return empty list
		return []*response.PostListResponse{}, &response.Pagination{
			Total: 0,
			Page:  page,
			Limit: limit,
		}, nil
	}

	if len(communityScores) == 0 {
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

	// Fetch posts with weighted sampling based on community scores
	allPosts := []*model.Post{}
	totalFetched := 0

	for communityID, score := range communityScores {
		// Calculate dynamic post limit based on community score
		postsLimit := s.calculatePostsLimit(score)

		// Ensure we don't exceed MAX_TOTAL_POSTS
		if totalFetched+postsLimit > constant.RECOMMENDATION_MAX_TOTAL_POSTS {
			postsLimit = constant.RECOMMENDATION_MAX_TOTAL_POSTS - totalFetched
		}

		if postsLimit <= 0 {
			break
		}

		// Get recent posts from this community
		posts, _, err := s.postRepo.GetPostsByCommunityID(communityID, "new", 1, postsLimit, []string{}, &userID)
		if err != nil {
			log.Printf("[Warn] Error getting posts for community %d: %v", communityID, err)
			continue
		}
		allPosts = append(allPosts, posts...)
		totalFetched += len(posts)

		// Stop if we've reached the max total posts
		if totalFetched >= constant.RECOMMENDATION_MAX_TOTAL_POSTS {
			log.Printf("[Info] Reached max posts limit (%d), stopping fetch", constant.RECOMMENDATION_MAX_TOTAL_POSTS)
			break
		}
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
	posts, total, err := s.postRepo.GetPostsByCommunityID(communityID, "new", 1, 100, []string{}, &userID)
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

	// Sort by score descending using sort.Slice for O(n log n) performance
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	// Extract sorted posts
	result := make([]*model.Post, len(scored))
	for i, sp := range scored {
		result[i] = sp.post
	}

	return result
}

// calculatePostsLimit determines how many posts to fetch from a community based on its score
func (s *RecommendationService) calculatePostsLimit(score float64) int {
	if score >= constant.RECOMMENDATION_HIGH_SCORE_THRESHOLD {
		return constant.RECOMMENDATION_HIGH_SCORE_POSTS_LIMIT
	} else if score >= constant.RECOMMENDATION_MEDIUM_SCORE_THRESHOLD {
		return constant.RECOMMENDATION_MEDIUM_SCORE_POSTS_LIMIT
	} else if score >= constant.RECOMMENDATION_LOW_SCORE_THRESHOLD {
		return constant.RECOMMENDATION_LOW_SCORE_POSTS_LIMIT
	}
	return constant.RECOMMENDATION_MIN_SCORE_POSTS_LIMIT
}

func (s *RecommendationService) calculatePostScore(post *model.Post, preferredTags []string) float64 {
	score := 0.0

	// Factor 1: Post votes (engagement)
	score += float64(post.Vote) * constant.RECOMMENDATION_VOTE_WEIGHT

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
		score += float64(tagMatchCount) * constant.RECOMMENDATION_TAG_MATCH_WEIGHT
	}

	// Factor 3: Freshness (decay over time)
	// Posts from last 24 hours get bonus, then decay
	hoursSincePost := time.Since(post.CreatedAt).Hours()
	if hoursSincePost < constant.RECOMMENDATION_FRESHNESS_RECENT_HOURS {
		score += constant.RECOMMENDATION_FRESHNESS_RECENT_BONUS * (1.0 - hoursSincePost/constant.RECOMMENDATION_FRESHNESS_RECENT_HOURS)
	} else if hoursSincePost < constant.RECOMMENDATION_FRESHNESS_WEEK_HOURS {
		score += constant.RECOMMENDATION_FRESHNESS_WEEK_BONUS * (1.0 - (hoursSincePost-constant.RECOMMENDATION_FRESHNESS_RECENT_HOURS)/constant.RECOMMENDATION_FRESHNESS_DECAY_WINDOW)
	}

	// Factor 4: Author karma (quality indicator)
	if post.Author != nil {
		score += float64(post.Author.Karma) / constant.RECOMMENDATION_KARMA_DIVIDER
	}

	return score
}
