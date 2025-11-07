package repository

import (
	"encoding/json"
	"fmt"
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"
	"social-platform-backend/internal/interface/dto/request"
	"social-platform-backend/package/constant"
	"social-platform-backend/package/util"
	"time"

	"gorm.io/gorm"
)

type PostRepositoryImpl struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) repository.PostRepository {
	return &PostRepositoryImpl{db: db}
}

func (r *PostRepositoryImpl) CreatePost(post *model.Post) error {
	return r.db.Create(post).Error
}

func (r *PostRepositoryImpl) GetPostByID(id uint64) (*model.Post, error) {
	var post model.Post
	err := r.db.Where("id = ?", id).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PostRepositoryImpl) GetPostDetailByID(id uint64, userID *uint64) (*model.Post, error) {
	var post model.Post

	selectFields := `posts.*,
		COALESCE(SUM(CASE WHEN post_votes.vote = true THEN 1 WHEN post_votes.vote = false THEN -1 ELSE 0 END), 0) as vote`

	// Add user_vote field if userID exists
	if userID != nil {
		selectFields += fmt.Sprintf(", MAX(CASE WHEN user_post_votes.user_id = %d THEN CAST(user_post_votes.vote AS INT) ELSE NULL END) as user_vote", *userID)
	}

	query := r.db.Table("posts").
		Select(selectFields).
		Joins("LEFT JOIN post_votes ON posts.id = post_votes.post_id")

	// Join with user's votes if userID exists
	if userID != nil {
		query = query.Joins("LEFT JOIN post_votes as user_post_votes ON posts.id = user_post_votes.post_id AND user_post_votes.user_id = ?", *userID)
	}

	err := query.Where("posts.id = ?", id).
		Group("posts.id").
		Preload("Community").
		Preload("Author").
		First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PostRepositoryImpl) UpdatePostText(id uint64, updatePost *request.UpdatePostTextRequest) error {
	updates := make(map[string]interface{})
	if updatePost.Title != nil {
		updates["title"] = *updatePost.Title
	}
	if updatePost.Content != nil {
		updates["content"] = *updatePost.Content
	}
	if updatePost.Tags != nil {
		updates["tags"] = *updatePost.Tags
	}
	if len(updates) > 0 {
		now := time.Now()
		updates["updated_at"] = now
	}
	return r.db.Model(&model.Post{}).Where("id = ?", id).Updates(updates).Error
}

func (r *PostRepositoryImpl) UpdatePostLink(id uint64, updatePost *request.UpdatePostLinkRequest) error {
	updates := make(map[string]interface{})
	if updatePost.Title != nil {
		updates["title"] = *updatePost.Title
	}
	if updatePost.Content != nil {
		updates["content"] = *updatePost.Content
	}
	if updatePost.URL != nil {
		updates["url"] = *updatePost.URL
	}
	if updatePost.Tags != nil {
		updates["tags"] = *updatePost.Tags
	}
	if len(updates) > 0 {
		now := time.Now()
		updates["updated_at"] = now
	}
	return r.db.Model(&model.Post{}).Where("id = ?", id).Updates(updates).Error
}

func (r *PostRepositoryImpl) UpdatePostMedia(id uint64, updatePost *request.UpdatePostMediaRequest) error {
	updates := make(map[string]interface{})
	if updatePost.Title != nil {
		updates["title"] = *updatePost.Title
	}
	if updatePost.Content != nil {
		updates["content"] = *updatePost.Content
	}
	if updatePost.MediaURLs != nil {
		updates["media_urls"] = *updatePost.MediaURLs
	}
	if updatePost.Tags != nil {
		updates["tags"] = *updatePost.Tags
	}
	if len(updates) > 0 {
		now := time.Now()
		updates["updated_at"] = now
	}
	return r.db.Model(&model.Post{}).Where("id = ?", id).Updates(updates).Error
}

func (r *PostRepositoryImpl) UpdatePostPoll(id uint64, updatePost *request.UpdatePostPollRequest) error {
	updates := make(map[string]interface{})
	if updatePost.Title != nil {
		updates["title"] = *updatePost.Title
	}
	if updatePost.Content != nil {
		updates["content"] = *updatePost.Content
	}
	if updatePost.PollData != nil {
		updates["poll_data"] = *updatePost.PollData
	}
	if updatePost.Tags != nil {
		updates["tags"] = *updatePost.Tags
	}
	if len(updates) > 0 {
		now := time.Now()
		updates["updated_at"] = now
	}
	return r.db.Model(&model.Post{}).Where("id = ?", id).Updates(updates).Error
}

func (r *PostRepositoryImpl) DeletePost(id uint64) error {
	return r.db.Model(&model.Post{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}

func (r *PostRepositoryImpl) UpdatePostStatus(id uint64, status string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	return r.db.Model(&model.Post{}).Where("id = ?", id).Updates(updates).Error
}

func (r *PostRepositoryImpl) UpdatePollData(postID uint64, pollData *json.RawMessage) error {
	updates := map[string]interface{}{
		"poll_data": pollData,
	}
	return r.db.Model(&model.Post{}).Where("id = ?", postID).Updates(updates).Error
}

func (r *PostRepositoryImpl) GetAllPosts(sortBy string, page, limit int, userID *uint64) ([]*model.Post, int64, error) {
	var posts []*model.Post
	var total int64

	// Count total APPROVED posts
	if err := r.db.Model(&model.Post{}).
		Where("status = ?", constant.POST_STATUS_APPROVED).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	selectFields := `posts.*,
		COALESCE(SUM(CASE WHEN post_votes.vote = true THEN 1 WHEN post_votes.vote = false THEN -1 ELSE 0 END), 0) as vote`

	// Add user_vote field if userID exists
	if userID != nil {
		selectFields += fmt.Sprintf(", MAX(CASE WHEN user_post_votes.user_id = %d THEN CAST(user_post_votes.vote AS INT) ELSE NULL END) as user_vote", *userID)
	}

	query := r.db.Table("posts").
		Select(selectFields).
		Joins("LEFT JOIN post_votes ON posts.id = post_votes.post_id")

	// Join with user's votes if userID exists
	if userID != nil {
		query = query.Joins("LEFT JOIN post_votes as user_post_votes ON posts.id = user_post_votes.post_id AND user_post_votes.user_id = ?", *userID)
	}

	query = query.Where("posts.status = ?", constant.POST_STATUS_APPROVED).
		Group("posts.id").
		Preload("Community").
		Preload("Author")

	// Sort method
	switch sortBy {
	case "hot", "top":
		query = query.Order("vote DESC")
	case "new", "best":
		fallthrough
	default:
		query = query.Order("posts.created_at DESC")
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (r *PostRepositoryImpl) GetPostsByCommunityID(communityID uint64, sortBy string, page, limit int, userID *uint64) ([]*model.Post, int64, error) {
	var posts []*model.Post
	var total int64

	// Count total APPROVED posts in community
	if err := r.db.Model(&model.Post{}).
		Where("community_id = ? AND status = ?", communityID, constant.POST_STATUS_APPROVED).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	selectFields := `posts.*,
		COALESCE(SUM(CASE WHEN post_votes.vote = true THEN 1 WHEN post_votes.vote = false THEN -1 ELSE 0 END), 0) as vote`

	// Add user_vote field if userID exists
	if userID != nil {
		selectFields += fmt.Sprintf(", MAX(CASE WHEN user_post_votes.user_id = %d THEN CAST(user_post_votes.vote AS INT) ELSE NULL END) as user_vote", *userID)
	}

	query := r.db.Table("posts").
		Select(selectFields).
		Joins("LEFT JOIN post_votes ON posts.id = post_votes.post_id")

	// Join with user's votes if userID exists
	if userID != nil {
		query = query.Joins("LEFT JOIN post_votes as user_post_votes ON posts.id = user_post_votes.post_id AND user_post_votes.user_id = ?", *userID)
	}

	query = query.Where("posts.community_id = ? AND posts.status = ?", communityID, constant.POST_STATUS_APPROVED).
		Group("posts.id").
		Preload("Community").
		Preload("Author")

	// Sort method
	switch sortBy {
	case "hot", "top":
		query = query.Order("vote DESC")
	case "new", "best":
		fallthrough
	default:
		query = query.Order("posts.created_at DESC")
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (r *PostRepositoryImpl) SearchPostsByTitle(title, sortBy string, page, limit int, userID *uint64) ([]*model.Post, int64, error) {
	var posts []*model.Post
	var total int64

	offset := (page - 1) * limit

	patterns := util.BuildSearchPattern(title)

	countQuery := r.db.Model(&model.Post{}).
		Where("status = ?", constant.POST_STATUS_APPROVED)
	for _, p := range patterns {
		countQuery = countQuery.Where("unaccent(lower(title)) LIKE unaccent(lower(?))", p)
	}
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	selectFields := `posts.*,
		COALESCE(SUM(CASE WHEN post_votes.vote = true THEN 1 WHEN post_votes.vote = false THEN -1 ELSE 0 END), 0) as vote`

	// Add user_vote field if userID exists
	if userID != nil {
		selectFields += fmt.Sprintf(", MAX(CASE WHEN user_post_votes.user_id = %d THEN CAST(user_post_votes.vote AS INT) ELSE NULL END) as user_vote", *userID)
	}

	query := r.db.Table("posts").
		Select(selectFields).
		Joins("LEFT JOIN post_votes ON posts.id = post_votes.post_id")

	// Join with user's votes if userID exists
	if userID != nil {
		query = query.Joins("LEFT JOIN post_votes as user_post_votes ON posts.id = user_post_votes.post_id AND user_post_votes.user_id = ?", *userID)
	}

	for _, p := range patterns {
		query = query.Where("unaccent(lower(posts.title)) LIKE unaccent(lower(?)) AND posts.status = ?", p, constant.POST_STATUS_APPROVED)
	}

	query = query.Group("posts.id").
		Preload("Community").
		Preload("Author")

	// Sort method
	switch sortBy {
	case "hot", "top":
		query = query.Order("vote DESC")
	case "new", "best":
		fallthrough
	default:
		query = query.Order("posts.created_at DESC")
	}

	if err := query.Offset(offset).Limit(limit).Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (r *PostRepositoryImpl) GetPostsByUserID(userID uint64, sortBy string, page, limit int) ([]*model.Post, int64, error) {
	var posts []*model.Post
	var total int64

	// Count total posts by user
	if err := r.db.Model(&model.Post{}).
		Where("author_id = ? AND deleted_at IS NULL", userID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := r.db.Table("posts").
		Select(`posts.*,
			COALESCE(SUM(CASE WHEN post_votes.vote = true THEN 1 WHEN post_votes.vote = false THEN -1 ELSE 0 END), 0) as vote,
			COUNT(DISTINCT comments.id) as comment_count`).
		Joins("LEFT JOIN post_votes ON posts.id = post_votes.post_id").
		Joins("LEFT JOIN comments ON posts.id = comments.post_id AND comments.deleted_at IS NULL").
		Where("posts.author_id = ? AND posts.status = ? AND posts.deleted_at IS NULL", userID, constant.POST_STATUS_APPROVED).
		Group("posts.id").
		Preload("Community").
		Preload("Author")

	switch sortBy {
	case constant.SORT_TOP:
		query = query.Order(`"vote" DESC, posts.created_at DESC`)
	case constant.SORT_HOT:
		query = query.Order(`comment_count DESC, "vote" DESC, posts.created_at DESC`)
	case constant.SORT_NEW:
		fallthrough
	default:
		query = query.Order("posts.created_at DESC")
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (r *PostRepositoryImpl) GetCommunityPostsForModerator(communityID uint64, status, searchTitle string, page, limit int) ([]*model.Post, int64, error) {
	var posts []*model.Post
	var total int64

	countQuery := r.db.Model(&model.Post{}).Where("community_id = ? AND deleted_at IS NULL", communityID)

	query := r.db.Table("posts").
		Select("posts.*").
		Where("posts.community_id = ? AND posts.deleted_at IS NULL", communityID)

	if status != "" {
		countQuery = countQuery.Where("status = ?", status)
		query = query.Where("posts.status = ?", status)
	}

	if searchTitle != "" {
		patterns := util.BuildSearchPattern(searchTitle)
		for _, p := range patterns {
			countQuery = countQuery.Where("unaccent(lower(title)) LIKE unaccent(lower(?))", p)
			query = query.Where("unaccent(lower(posts.title)) LIKE unaccent(lower(?))", p)
		}
	}

	// Count total
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = query.Preload("Author").Order("posts.created_at DESC")

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}
