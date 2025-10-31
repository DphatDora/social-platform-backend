package repository

import (
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"
	"social-platform-backend/internal/interface/dto/request"
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

func (r *PostRepositoryImpl) GetPostDetailByID(id uint64) (*model.Post, error) {
	var post model.Post
	err := r.db.Table("posts").
		Select(`posts.*,
			COALESCE(SUM(CASE WHEN post_votes.vote = true THEN 1 WHEN post_votes.vote = false THEN -1 ELSE 0 END), 0) as vote`).
		Joins("LEFT JOIN post_votes ON posts.id = post_votes.post_id").
		Where("posts.id = ?", id).
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

func (r *PostRepositoryImpl) GetAllPosts(sortBy string, page, limit int) ([]*model.Post, int64, error) {
	var posts []*model.Post
	var total int64

	// Count total posts
	if err := r.db.Model(&model.Post{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := r.db.Table("posts").
		Select(`posts.*,
			COALESCE(SUM(CASE WHEN post_votes.vote = true THEN 1 WHEN post_votes.vote = false THEN -1 ELSE 0 END), 0) as vote`).
		Joins("LEFT JOIN post_votes ON posts.id = post_votes.post_id").
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

func (r *PostRepositoryImpl) GetPostsByCommunityID(communityID uint64, sortBy string, page, limit int) ([]*model.Post, int64, error) {
	var posts []*model.Post
	var total int64

	// Count total posts in community
	if err := r.db.Model(&model.Post{}).Where("community_id = ?", communityID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := r.db.Table("posts").
		Select(`posts.*,
			COALESCE(SUM(CASE WHEN post_votes.vote = true THEN 1 WHEN post_votes.vote = false THEN -1 ELSE 0 END), 0) as vote`).
		Joins("LEFT JOIN post_votes ON posts.id = post_votes.post_id").
		Where("posts.community_id = ?", communityID).
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

func (r *PostRepositoryImpl) SearchPostsByTitle(title, sortBy string, page, limit int) ([]*model.Post, int64, error) {
	var posts []*model.Post
	var total int64

	offset := (page - 1) * limit

	patterns := util.BuildSearchPattern(title)

	countQuery := r.db.Model(&model.Post{})
	for _, p := range patterns {
		countQuery = countQuery.Where("unaccent(lower(title)) LIKE unaccent(lower(?))", p)
	}
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := r.db.Table("posts").
		Select(`posts.*,
			COALESCE(SUM(CASE WHEN post_votes.vote = true THEN 1 WHEN post_votes.vote = false THEN -1 ELSE 0 END), 0) as vote`).
		Joins("LEFT JOIN post_votes ON posts.id = post_votes.post_id")

	for _, p := range patterns {
		query = query.Where("unaccent(lower(posts.title)) LIKE unaccent(lower(?))", p)
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
