package repository

import "social-platform-backend/internal/domain/model"

type UserTagPreferenceRepository interface {
	UpsertTagPreferences(preference *model.UserTagPreference) error
	GetUserTagPreferences(userID uint64) (*model.UserTagPreference, error)
	UpdateTagsFromUserActivity(userID uint64) error
}
