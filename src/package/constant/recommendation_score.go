package constant

const (
	// Vote weight - multiplier for post votes
	RECOMMENDATION_VOTE_WEIGHT = 0.5

	// Tag matching weight - bonus points per matched tag
	RECOMMENDATION_TAG_MATCH_WEIGHT = 10.0

	// Freshness bonus for posts less than 24 hours old
	RECOMMENDATION_FRESHNESS_RECENT_BONUS = 20.0

	// Freshness bonus for posts between 24 hours and 7 days old
	RECOMMENDATION_FRESHNESS_WEEK_BONUS = 10.0

	// Karma weight - divider for author karma
	RECOMMENDATION_KARMA_DIVIDER = 100.0

	// Time thresholds (in hours)
	RECOMMENDATION_FRESHNESS_RECENT_HOURS = 24.0  // 24 hours
	RECOMMENDATION_FRESHNESS_WEEK_HOURS   = 168.0 // 7 days
	RECOMMENDATION_FRESHNESS_DECAY_WINDOW = 144.0 // 6 days (168 - 24)

	// Fetching optimization constants
	RECOMMENDATION_MAX_COMMUNITIES         = 8   // Max communities to fetch from
	RECOMMENDATION_MAX_TOTAL_POSTS         = 100 // Max total posts to score
	RECOMMENDATION_MIN_POSTS_PER_COMMUNITY = 5
	RECOMMENDATION_MAX_POSTS_PER_COMMUNITY = 20

	// Score thresholds for weighted sampling
	RECOMMENDATION_HIGH_SCORE_THRESHOLD   = 80.0
	RECOMMENDATION_MEDIUM_SCORE_THRESHOLD = 50.0
	RECOMMENDATION_LOW_SCORE_THRESHOLD    = 30.0

	// Posts limits based on community score
	RECOMMENDATION_HIGH_SCORE_POSTS_LIMIT   = 20
	RECOMMENDATION_MEDIUM_SCORE_POSTS_LIMIT = 15
	RECOMMENDATION_LOW_SCORE_POSTS_LIMIT    = 10
	RECOMMENDATION_MIN_SCORE_POSTS_LIMIT    = 5
)
