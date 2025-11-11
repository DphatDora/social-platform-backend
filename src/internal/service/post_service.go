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

type PostService struct {
	postRepo            repository.PostRepository
	communityRepo       repository.CommunityRepository
	postVoteRepo        repository.PostVoteRepository
	postReportRepo      repository.PostReportRepository
	botTaskRepo         repository.BotTaskRepository
	userRepo            repository.UserRepository
	tagRepo             repository.TagRepository
	notificationService *NotificationService
	botTaskService      *BotTaskService
	recommendService    *RecommendationService
}

func NewPostService(
	postRepo repository.PostRepository,
	communityRepo repository.CommunityRepository,
	postVoteRepo repository.PostVoteRepository,
	postReportRepo repository.PostReportRepository,
	botTaskRepo repository.BotTaskRepository,
	userRepo repository.UserRepository,
	tagRepo repository.TagRepository,
	notificationService *NotificationService,
	botTaskService *BotTaskService,
	recommendService *RecommendationService,
) *PostService {
	return &PostService{
		postRepo:            postRepo,
		communityRepo:       communityRepo,
		postVoteRepo:        postVoteRepo,
		postReportRepo:      postReportRepo,
		botTaskRepo:         botTaskRepo,
		userRepo:            userRepo,
		tagRepo:             tagRepo,
		notificationService: notificationService,
		botTaskService:      botTaskService,
		recommendService:    recommendService,
	}
}

func (s *PostService) CreatePost(userID uint64, req *request.CreatePostRequest) error {
	// Check if community exists
	community, err := s.communityRepo.GetCommunityByID(req.CommunityID)
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

	// Determine post status based on community settings
	postStatus := constant.POST_STATUS_PENDING
	if !community.RequiresApproval {
		postStatus = constant.POST_STATUS_APPROVED
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
		Status:      postStatus,
	}

	if err := s.postRepo.CreatePost(post); err != nil {
		log.Printf("[Err] Error creating post in PostService.CreatePost: %v", err)
		return fmt.Errorf("failed to create post")
	}

	go func(userID uint64) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[Panic] Recovered in PostService.CreatePost background task: %v", r)
			}
		}()

		// Create bottask for updating karma
		if s.botTaskService != nil {
			if err := s.botTaskService.CreateKarmaTask(userID, nil, constant.KARMA_ACTION_CREATE_POST); err != nil {
				log.Printf("[Err] Error creating karma task in PostService.CreatePost: %v", err)
			}
		}
	}(userID)

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

func (s *PostService) GetAllPosts(sortBy string, page, limit int, tags []string, userID *uint64) ([]*response.PostListResponse, *response.Pagination, error) {
	// If sortBy is "best", use recommendation service
	if sortBy == "best" && userID != nil && s.recommendService != nil {
		return s.recommendService.GetRecommendedPosts(*userID, page, limit)
	}

	posts, total, err := s.postRepo.GetAllPosts(sortBy, page, limit, tags, userID)
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

func (s *PostService) GetPostsByCommunityID(communityID uint64, sortBy string, page, limit int, tags []string, userID *uint64) ([]*response.PostListResponse, *response.Pagination, error) {
	// Check if community exists
	_, err := s.communityRepo.GetCommunityByID(communityID)
	if err != nil {
		log.Printf("[Err] Community not found in PostService.GetPostsByCommunityID: %v", err)
		return nil, nil, fmt.Errorf("community not found")
	}

	// If sortBy is "best", use recommendation service
	if sortBy == "best" && userID != nil && s.recommendService != nil {
		return s.recommendService.GetRecommendedPostsByCommunity(*userID, communityID, page, limit)
	}

	posts, total, err := s.postRepo.GetPostsByCommunityID(communityID, sortBy, page, limit, tags, userID)
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

func (s *PostService) SearchPostsByTitle(title, sortBy string, page, limit int, tags []string, userID *uint64) ([]*response.PostListResponse, *response.Pagination, error) {
	posts, total, err := s.postRepo.SearchPostsByTitle(title, sortBy, page, limit, tags, userID)
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

func (s *PostService) GetPostDetailByID(postID uint64, userID *uint64) (*response.PostDetailResponse, error) {
	post, err := s.postRepo.GetPostDetailByID(postID, userID)
	if err != nil {
		log.Printf("[Err] Error getting post detail in PostService.GetPostDetailByID: %v", err)
		return nil, fmt.Errorf("post not found")
	}

	return response.NewPostDetailResponse(post), nil
}

func (s *PostService) VotePost(userID, postID uint64, vote bool) error {
	// Check if post exists
	post, err := s.postRepo.GetPostByID(postID)
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

	go func(userID uint64, post *model.Post, vote bool) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[Panic] Recovered in PostService.VotePost background task: %v", r)
			}
		}()

		// Create bot task for updating karma
		var karmaAction string
		if vote {
			karmaAction = constant.KARMA_ACTION_UPVOTE_POST
		} else {
			karmaAction = constant.KARMA_ACTION_DOWNVOTE_POST
		}

		if s.botTaskService != nil {
			if err := s.botTaskService.CreateKarmaTask(userID, &post.AuthorID, karmaAction); err != nil {
				log.Printf("[Err] Error creating karma task in PostService.VotePost: %v", err)
			}
		}

		// Create bot task for updating interest score
		if s.botTaskService != nil {
			var interestAction string
			if vote {
				interestAction = constant.INTEREST_ACTION_UPVOTE_POST
			} else {
				interestAction = constant.INTEREST_ACTION_DOWNVOTE_POST
			}

			postIDPtr := &post.ID
			if err := s.botTaskService.CreateInterestScoreTask(userID, post.CommunityID, interestAction, postIDPtr); err != nil {
				log.Printf("[Err] Error creating interest score task in PostService.VotePost: %v", err)
			}
		}

		// Send notification to post author (if not voting own post)
		if s.notificationService != nil && userID != post.AuthorID {
			voter, err := s.userRepo.GetUserByID(userID)
			if err != nil {
				log.Printf("[Err] Error getting voter in PostService.VotePost: %v", err)
				return
			}

			notifPayload := payload.PostVoteNotificationPayload{
				PostID:   post.ID,
				UserName: voter.Username,
				VoteType: vote,
			}
			s.notificationService.CreateNotification(
				post.AuthorID,
				constant.NOTIFICATION_ACTION_GET_POST_VOTE,
				notifPayload,
			)
		}
	}(userID, post, vote)

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

func (s *PostService) VotePoll(userID, postID uint64, req *request.VotePollRequest) error {
	post, err := s.postRepo.GetPostByID(postID)
	if err != nil {
		log.Printf("[Err] Post not found in PostService.VotePoll: %v", err)
		return fmt.Errorf("post not found")
	}

	if post.Type != constant.PostTypePoll {
		return fmt.Errorf("post is not a poll")
	}

	if post.PollData == nil {
		log.Printf("[Err] Poll data is nil in PostService.VotePoll")
		return fmt.Errorf("poll data not found")
	}

	var pollData payload.PollData
	if err := json.Unmarshal(*post.PollData, &pollData); err != nil {
		log.Printf("[Err] Error unmarshalling poll data in PostService.VotePoll: %v", err)
		return fmt.Errorf("invalid poll data")
	}

	// Check expiration time
	if pollData.ExpiresAt != nil && time.Now().After(*pollData.ExpiresAt) {
		return fmt.Errorf("poll has expired")
	}

	// Find the option
	optionIndex := -1
	for i, option := range pollData.Options {
		if option.ID == req.OptionID {
			optionIndex = i
			break
		}
	}

	if optionIndex == -1 {
		return fmt.Errorf("option not found")
	}

	// Check if user has already voted
	hasVoted := false
	var previousOptionIndex int
	for i, option := range pollData.Options {
		for _, voterID := range option.Voters {
			if voterID == userID {
				hasVoted = true
				previousOptionIndex = i
				break
			}
		}
		if hasVoted {
			break
		}
	}

	if hasVoted {
		// If not multiple choice and voting for same option, return error
		if !pollData.MultipleChoice && previousOptionIndex == optionIndex {
			return fmt.Errorf("already voted for this option")
		}

		// If not multiple choice and voting for different option, remove previous vote
		if !pollData.MultipleChoice && previousOptionIndex != optionIndex {
			// Remove from previous option
			newVoters := []uint64{}
			for _, voterID := range pollData.Options[previousOptionIndex].Voters {
				if voterID != userID {
					newVoters = append(newVoters, voterID)
				}
			}
			pollData.Options[previousOptionIndex].Voters = newVoters
			pollData.Options[previousOptionIndex].Votes = len(newVoters)
		}

		// If multiple choice, check if already voted for this option
		if pollData.MultipleChoice {
			for _, voterID := range pollData.Options[optionIndex].Voters {
				if voterID == userID {
					return fmt.Errorf("already voted for this option")
				}
			}
		}
	}

	// Add vote to the selected option
	pollData.Options[optionIndex].Voters = append(pollData.Options[optionIndex].Voters, userID)
	pollData.Options[optionIndex].Votes = len(pollData.Options[optionIndex].Voters)

	// Recalculate total votes (count unique voters)
	uniqueVoters := make(map[uint64]bool)
	for _, option := range pollData.Options {
		for _, voterID := range option.Voters {
			uniqueVoters[voterID] = true
		}
	}
	pollData.TotalVotes = len(uniqueVoters)

	updatedPollData, err := json.Marshal(pollData)
	if err != nil {
		log.Printf("[Err] Error marshalling poll data in PostService.VotePoll: %v", err)
		return fmt.Errorf("failed to update poll data")
	}

	rawMessage := json.RawMessage(updatedPollData)

	if err := s.postRepo.UpdatePollData(postID, &rawMessage); err != nil {
		log.Printf("[Err] Error updating poll data in PostService.VotePoll: %v", err)
		return fmt.Errorf("failed to update poll")
	}

	return nil
}

func (s *PostService) UnvotePoll(userID, postID uint64, req *request.UnvotePollRequest) error {
	post, err := s.postRepo.GetPostByID(postID)
	if err != nil {
		log.Printf("[Err] Post not found in PostService.UnvotePoll: %v", err)
		return fmt.Errorf("post not found")
	}

	if post.Type != constant.PostTypePoll {
		return fmt.Errorf("post is not a poll")
	}

	if post.PollData == nil {
		log.Printf("[Err] Poll data is nil in PostService.UnvotePoll")
		return fmt.Errorf("poll data not found")
	}

	var pollData payload.PollData
	if err := json.Unmarshal(*post.PollData, &pollData); err != nil {
		log.Printf("[Err] Error unmarshalling poll data in PostService.UnvotePoll: %v", err)
		return fmt.Errorf("invalid poll data")
	}

	// Find the option
	optionIndex := -1
	for i, option := range pollData.Options {
		if option.ID == req.OptionID {
			optionIndex = i
			break
		}
	}

	if optionIndex == -1 {
		return fmt.Errorf("option not found")
	}

	// Check if user has voted for this specific option
	hasVoted := false
	for _, voterID := range pollData.Options[optionIndex].Voters {
		if voterID == userID {
			hasVoted = true
			break
		}
	}

	if !hasVoted {
		return fmt.Errorf("you have not voted for this option")
	}

	// Remove user's vote from this specific option
	newVoters := []uint64{}
	for _, voterID := range pollData.Options[optionIndex].Voters {
		if voterID != userID {
			newVoters = append(newVoters, voterID)
		}
	}
	pollData.Options[optionIndex].Voters = newVoters
	pollData.Options[optionIndex].Votes = len(newVoters)

	// Recalculate total votes (count unique voters)
	uniqueVoters := make(map[uint64]bool)
	for _, option := range pollData.Options {
		for _, voterID := range option.Voters {
			uniqueVoters[voterID] = true
		}
	}
	pollData.TotalVotes = len(uniqueVoters)

	updatedPollData, err := json.Marshal(pollData)
	if err != nil {
		log.Printf("[Err] Error marshalling poll data in PostService.UnvotePoll: %v", err)
		return fmt.Errorf("failed to update poll data")
	}

	rawMessage := json.RawMessage(updatedPollData)

	if err := s.postRepo.UpdatePollData(postID, &rawMessage); err != nil {
		log.Printf("[Err] Error updating poll data in PostService.UnvotePoll: %v", err)
		return fmt.Errorf("failed to update poll")
	}

	return nil
}

func (s *PostService) GetPostsByUserID(userID uint64, sortBy string, page, limit int) ([]*response.PostListResponse, *response.Pagination, error) {
	// Check if user exists
	_, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		log.Printf("[Err] User not found in PostService.GetPostsByUserID: %v", err)
		return nil, nil, fmt.Errorf("user not found")
	}

	posts, total, err := s.postRepo.GetPostsByUserID(userID, sortBy, page, limit)
	if err != nil {
		log.Printf("[Err] Error getting posts by user ID in PostService.GetPostsByUserID: %v", err)
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

	return postResponses, pagination, nil
}

func (s *PostService) ReportPost(userID, postID uint64, req *request.ReportPostRequest) error {
	// Check if post exists
	post, err := s.postRepo.GetPostByID(postID)
	if err != nil {
		log.Printf("[Err] Post not found in PostService.ReportPost: %v", err)
		return fmt.Errorf("post not found")
	}

	// Check if user already reported this post
	alreadyReported, err := s.postReportRepo.IsUserReportedPost(userID, postID)
	if err != nil {
		log.Printf("[Err] Error checking if user reported post in PostService.ReportPost: %v", err)
		return fmt.Errorf("failed to check report status")
	}
	if alreadyReported {
		return fmt.Errorf("you have already reported this post")
	}

	// Create report
	report := &model.PostReport{
		PostID:     postID,
		ReporterID: userID,
		Reasons:    req.Reasons,
		Note:       req.Note,
		CreatedAt:  time.Now(),
	}

	if err := s.postReportRepo.CreatePostReport(report); err != nil {
		log.Printf("[Err] Error creating post report in PostService.ReportPost: %v", err)
		return fmt.Errorf("failed to report post")
	}

	// Send notification to post author
	go func(authorID uint64, postID uint64, reporterID uint64) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[Panic] Recovered in PostService.ReportPost notification: %v", r)
			}
		}()

		// Get reporter info
		reporter, err := s.userRepo.GetUserByID(reporterID)
		if err != nil {
			log.Printf("[Err] Error getting reporter in PostService.ReportPost: %v", err)
			return
		}

		notifPayload := payload.PostReportNotificationPayload{
			PostID:   postID,
			UserName: reporter.Username,
		}

		s.notificationService.CreateNotification(authorID, constant.NOTIFICATION_ACTION_POST_REPORTED, notifPayload)
	}(post.AuthorID, postID, userID)

	return nil
}

func (s *PostService) GetAllTags(search *string) ([]*response.TagResponse, error) {
	tags, err := s.tagRepo.GetAllTags(search)
	if err != nil {
		log.Printf("[Err] Error getting tags in PostService.GetAllTags: %v", err)
		return nil, err
	}

	tagResponses := make([]*response.TagResponse, len(tags))
	for i, tag := range tags {
		tagResponses[i] = &response.TagResponse{
			ID:   tag.ID,
			Name: tag.Name,
		}
	}

	return tagResponses, nil
}
