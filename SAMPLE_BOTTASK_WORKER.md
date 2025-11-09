
``` go
// =================================================================
// BACKGROUND WORKER SERVICE - INTEREST SCORE PROCESSOR
// =================================================================
// This is a SAMPLE CODE for the separate background worker service
// that processes bot tasks from the database.
//
// This code should be placed in your BACKGROUND WORKER project,
// NOT in the main API project.
//
// The worker will:
// 1. Fetch unprocessed bot tasks from database
// 2. Process tasks based on action type
// 3. Update interest scores and tag preferences
// 4. Mark tasks as executed
// =================================================================

package worker

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

// =================================================================
// MODELS (Copy these to your worker project)
// =================================================================

type BotTask struct {
	ID         uint64           `gorm:"column:id;primaryKey"`
	Action     string           `gorm:"column:action"`
	Payload    *json.RawMessage `gorm:"column:payload"`
	CreatedAt  time.Time        `gorm:"column:created_at"`
	ExecutedAt *time.Time       `gorm:"column:executed_at"`
}

func (BotTask) TableName() string {
	return "bot_tasks"
}

type UserInterestScore struct {
	UserID      uint64     `gorm:"column:user_id;primaryKey"`
	CommunityID uint64     `gorm:"column:community_id;primaryKey"`
	Score       float64    `gorm:"column:score;default:0;index"`
	LastVoteAt  *time.Time `gorm:"column:last_vote_at"`
	LastJoinAt  *time.Time `gorm:"column:last_join_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (UserInterestScore) TableName() string {
	return "user_interest_scores"
}

// =================================================================
// PAYLOAD STRUCTURES
// =================================================================

type UpdateInterestScorePayload struct {
	UserID      uint64    `json:"user_id"`
	CommunityID uint64    `json:"community_id"`
	Action      string    `json:"action"` // upvote_post, downvote_post, follow_post, join_community
	PostID      *uint64   `json:"post_id,omitempty"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// =================================================================
// CONSTANTS
// =================================================================

const (
	BOT_TASK_ACTION_UPDATE_INTEREST_SCORE = "update_interest_score"

	INTEREST_ACTION_UPVOTE_POST    = "upvote_post"
	INTEREST_ACTION_DOWNVOTE_POST  = "downvote_post"
	INTEREST_ACTION_FOLLOW_POST    = "follow_post"
	INTEREST_ACTION_JOIN_COMMUNITY = "join_community"
)

// =================================================================
// INTEREST SCORE PROCESSOR
// =================================================================

type InterestScoreProcessor struct {
	db *gorm.DB
}

func NewInterestScoreProcessor(db *gorm.DB) *InterestScoreProcessor {
	return &InterestScoreProcessor{db: db}
}

// ProcessBotTasks fetches and processes unprocessed bot tasks
func (p *InterestScoreProcessor) ProcessBotTasks(batchSize int) error {
	// Fetch unprocessed tasks
	var tasks []BotTask
	err := p.db.Where("action = ? AND executed_at IS NULL", BOT_TASK_ACTION_UPDATE_INTEREST_SCORE).
		Order("created_at ASC").
		Limit(batchSize).
		Find(&tasks).Error

	if err != nil {
		return fmt.Errorf("failed to fetch bot tasks: %w", err)
	}

	log.Printf("[Info] Found %d interest score tasks to process", len(tasks))

	// Process each task
	for _, task := range tasks {
		if err := p.ProcessSingleTask(&task); err != nil {
			log.Printf("[Err] Failed to process task %d: %v", task.ID, err)
			continue
		}

		// Mark task as executed
		now := time.Now()
		if err := p.db.Model(&BotTask{}).Where("id = ?", task.ID).Update("executed_at", now).Error; err != nil {
			log.Printf("[Err] Failed to mark task %d as executed: %v", task.ID, err)
		}
	}

	return nil
}

// ProcessSingleTask processes a single bot task
func (p *InterestScoreProcessor) ProcessSingleTask(task *BotTask) error {
	// Parse payload
	var payload UpdateInterestScorePayload
	if err := json.Unmarshal(*task.Payload, &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	// Calculate score delta based on action
	scoreDelta := getScoreDeltaForAction(payload.Action)

	log.Printf("[Info] Processing task: UserID=%d, CommunityID=%d, Action=%s, ScoreDelta=%.2f",
		payload.UserID, payload.CommunityID, payload.Action, scoreDelta)

	// Update interest score
	if err := p.updateInterestScore(payload.UserID, payload.CommunityID, scoreDelta, payload.Action); err != nil {
		return fmt.Errorf("failed to update interest score: %w", err)
	}

	return nil
}

// updateInterestScore updates or creates a user interest score record
func (p *InterestScoreProcessor) updateInterestScore(userID, communityID uint64, scoreDelta float64, action string) error {
	now := time.Now()

	// Try to find existing record
	var existingScore UserInterestScore
	err := p.db.Where("user_id = ? AND community_id = ?", userID, communityID).
		First(&existingScore).Error

	if err == gorm.ErrRecordNotFound {
		// Create new record
		newScore := &UserInterestScore{
			UserID:      userID,
			CommunityID: communityID,
			Score:       scoreDelta,
			UpdatedAt:   now,
		}

		// Set specific timestamp based on action
		switch action {
		case INTEREST_ACTION_UPVOTE_POST, INTEREST_ACTION_DOWNVOTE_POST:
			newScore.LastVoteAt = &now
		case INTEREST_ACTION_JOIN_COMMUNITY:
			newScore.LastJoinAt = &now
		}

		if err := p.db.Create(newScore).Error; err != nil {
			return fmt.Errorf("failed to create interest score: %w", err)
		}

		log.Printf("[Info] Created new interest score: UserID=%d, CommunityID=%d, Score=%.2f",
			userID, communityID, scoreDelta)
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to query interest score: %w", err)
	}

	// Update existing record
	updates := map[string]interface{}{
		"score":      gorm.Expr("score + ?", scoreDelta),
		"updated_at": now,
	}

	// Update specific timestamps based on action
	switch action {
	case INTEREST_ACTION_UPVOTE_POST, INTEREST_ACTION_DOWNVOTE_POST:
		updates["last_vote_at"] = now
	case INTEREST_ACTION_JOIN_COMMUNITY:
		updates["last_join_at"] = now
	}

	if err := p.db.Model(&UserInterestScore{}).
		Where("user_id = ? AND community_id = ?", userID, communityID).
		Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update interest score: %w", err)
	}

	newScore := existingScore.Score + scoreDelta
	log.Printf("[Info] Updated interest score: UserID=%d, CommunityID=%d, OldScore=%.2f, NewScore=%.2f",
		userID, communityID, existingScore.Score, newScore)

	return nil
}

// getScoreDeltaForAction returns the score delta for a given action
func getScoreDeltaForAction(action string) float64 {
	switch action {
	case INTEREST_ACTION_UPVOTE_POST:
		return 2.0 // Upvote adds 2 points
	case INTEREST_ACTION_DOWNVOTE_POST:
		return -1.0 // Downvote subtracts 1 point
	case INTEREST_ACTION_FOLLOW_POST:
		return 3.0 // Following a post adds 3 points
	case INTEREST_ACTION_JOIN_COMMUNITY:
		return 10.0 // Joining community adds 10 points (highest weight)
	default:
		return 0.0
	}
}

// =================================================================
// USAGE EXAMPLE - Main Worker Loop
// =================================================================

/*
func main() {
	// Initialize database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create processor
	processor := NewInterestScoreProcessor(db)

	// Run worker loop
	ticker := time.NewTicker(10 * time.Second) // Process every 10 seconds
	defer ticker.Stop()

	log.Println("[Info] Interest Score Worker started")

	for {
		select {
		case <-ticker.C:
			if err := processor.ProcessBotTasks(100); err != nil {
				log.Printf("[Err] Error processing bot tasks: %v", err)
			}
		}
	}
}
*/

// =================================================================
// TAG PREFERENCE UPDATER (Optional - Run periodically)
// =================================================================

type TagPreferenceUpdater struct {
	db *gorm.DB
}

func NewTagPreferenceUpdater(db *gorm.DB) *TagPreferenceUpdater {
	return &TagPreferenceUpdater{db: db}
}

// UpdateAllUserTagPreferences updates tag preferences for all active users
// This should be run periodically (e.g., once per day)
func (u *TagPreferenceUpdater) UpdateAllUserTagPreferences() error {
	// Get all users who have voted or followed posts in the last 30 days
	var userIDs []uint64
	query := `
		SELECT DISTINCT user_id FROM (
			SELECT user_id FROM post_votes WHERE voted_at > NOW() - INTERVAL '30 days'
			UNION
			SELECT user_id FROM user_saved_posts WHERE is_followed = true
		) as active_users
	`

	if err := u.db.Raw(query).Pluck("user_id", &userIDs).Error; err != nil {
		return fmt.Errorf("failed to fetch active users: %w", err)
	}

	log.Printf("[Info] Updating tag preferences for %d users", len(userIDs))

	for _, userID := range userIDs {
		if err := u.updateUserTagPreference(userID); err != nil {
			log.Printf("[Err] Failed to update tag preferences for user %d: %v", userID, err)
			continue
		}
	}

	return nil
}

func (u *TagPreferenceUpdater) updateUserTagPreference(userID uint64) error {
	// Get tags from posts that user has voted on or followed
	var tags []string

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

	if err := u.db.Raw(query, userID, userID).Pluck("tag", &tags).Error; err != nil {
		return fmt.Errorf("failed to fetch tags: %w", err)
	}

	if len(tags) == 0 {
		return nil // No tags to update
	}

	// Convert to PostgreSQL array format
	tagsJSON, _ := json.Marshal(tags)

	// Upsert tag preferences
	result := u.db.Exec(`
		INSERT INTO user_tag_preferences (user_id, preferred_tags, updated_at)
		VALUES (?, ?::text[], NOW())
		ON CONFLICT (user_id)
		DO UPDATE SET preferred_tags = ?::text[], updated_at = NOW()
	`, userID, string(tagsJSON), string(tagsJSON))

	if result.Error != nil {
		return fmt.Errorf("failed to upsert tag preferences: %w", result.Error)
	}

	log.Printf("[Info] Updated tag preferences for user %d: %d tags", userID, len(tags))
	return nil
}

// =================================================================
// SCORING ALGORITHM NOTES
// =================================================================

/*
SCORING WEIGHTS:
- Join Community: +10 points (highest commitment)
- Follow Post: +3 points (medium-term interest)
- Upvote Post: +2 points (positive engagement)
- Downvote Post: -1 point (negative signal)

RECOMMENDATION ALGORITHM:
1. Get top communities by score for user (score > 0)
2. Fetch recent posts from those communities
3. Score posts based on:
   - Tag matching with user preferences (+10 per matching tag)
   - Post votes/engagement (+0.5 per vote)
   - Freshness (decay over time, max +20 for new posts)
   - Author karma (+karma/100)
4. Sort by total score descending
5. Return paginated results

DATABASE INDEXES REQUIRED:
- user_interest_scores(user_id, score DESC)
- user_interest_scores(community_id, score DESC)
- user_tag_preferences using GIN index on preferred_tags
- post_votes(user_id, voted_at DESC)
- user_saved_posts(user_id, is_followed)

PERFORMANCE CONSIDERATIONS:
- Process bot tasks in batches (e.g., 100 at a time)
- Update tag preferences daily or weekly, not real-time
- Cache top communities per user (can use Redis)
- Limit recommendation queries to recent posts (last 30 days)
*/
```