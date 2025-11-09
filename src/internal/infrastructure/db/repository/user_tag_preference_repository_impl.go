package repository

import (
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type UserTagPreferenceRepositoryImpl struct {
	db *gorm.DB
}

func NewUserTagPreferenceRepository(db *gorm.DB) repository.UserTagPreferenceRepository {
	return &UserTagPreferenceRepositoryImpl{db: db}
}

func (r *UserTagPreferenceRepositoryImpl) UpsertTagPreferences(preference *model.UserTagPreference) error {
	return r.db.Save(preference).Error
}

func (r *UserTagPreferenceRepositoryImpl) GetUserTagPreferences(userID uint64) (*model.UserTagPreference, error) {
	var preference model.UserTagPreference
	err := r.db.Where("user_id = ?", userID).First(&preference).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &preference, nil
}

func (r *UserTagPreferenceRepositoryImpl) UpdateTagsFromUserActivity(userID uint64) error {
	// Get tags from posts that user has voted on or followed
	var tags []string

	// Get tags from upvoted posts
	query := `
		SELECT DISTINCT unnest(p.tags) as tag
		FROM posts p
		INNER JOIN post_votes pv ON p.id = pv.post_id
		WHERE pv.user_id = ? AND pv.vote = true AND p.tags IS NOT NULL
		UNION
		SELECT DISTINCT unnest(p.tags) as tag
		FROM posts p
		INNER JOIN user_saved_posts usp ON p.id = usp.post_id
		WHERE usp.user_id = ? AND usp.is_followed = true AND p.tags IS NOT NULL
		LIMIT 50
	`

	err := r.db.Raw(query, userID, userID).Pluck("tag", &tags).Error
	if err != nil {
		return err
	}

	if len(tags) == 0 {
		return nil
	}

	// Upsert tag preferences
	preference := &model.UserTagPreference{
		UserID:        userID,
		PreferredTags: pq.StringArray(tags),
	}

	return r.UpsertTagPreferences(preference)
}
