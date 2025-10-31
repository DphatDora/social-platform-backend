package repository

import (
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/interface/dto/request"
)

type PostRepository interface {
	CreatePost(post *model.Post) error
	GetPostByID(id uint64) (*model.Post, error)
	GetPostDetailByID(id uint64) (*model.Post, error)
	UpdatePostText(id uint64, updatePost *request.UpdatePostTextRequest) error
	UpdatePostLink(id uint64, updatePost *request.UpdatePostLinkRequest) error
	UpdatePostMedia(id uint64, updatePost *request.UpdatePostMediaRequest) error
	UpdatePostPoll(id uint64, updatePost *request.UpdatePostPollRequest) error
	DeletePost(id uint64) error
	GetAllPosts(sortBy string, page, limit int) ([]*model.Post, int64, error)
	GetPostsByCommunityID(communityID uint64, sortBy string, page, limit int) ([]*model.Post, int64, error)
	SearchPostsByTitle(title, sortBy string, page, limit int) ([]*model.Post, int64, error)
}
