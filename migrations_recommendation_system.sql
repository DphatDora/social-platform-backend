-- =================================================================
-- MIGRATION: Add Recommendation System Tables
-- Description: Creates tables for user interest scores and tag preferences
-- Date: 2025-11-08
-- =================================================================

-- Table 1: User Interest Scores
-- Stores user's interest in communities based on their interactions
CREATE TABLE IF NOT EXISTS user_interest_scores (
    user_id BIGINT NOT NULL,
    community_id BIGINT NOT NULL,
    score DOUBLE PRECISION DEFAULT 0,
    last_vote_at TIMESTAMP,
    last_join_at TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, community_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (community_id) REFERENCES communities(id) ON DELETE CASCADE
);

-- Indexes for fast lookups
CREATE INDEX IF NOT EXISTS idx_user_interest_score ON user_interest_scores(user_id, score DESC);
CREATE INDEX IF NOT EXISTS idx_community_score ON user_interest_scores(community_id, score DESC);

-- Table 2: User Tag Preferences
-- Caches user's preferred tags based on their activity
CREATE TABLE IF NOT EXISTS user_tag_preferences (
    user_id BIGINT PRIMARY KEY,
    preferred_tags TEXT[],
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- GIN index for fast array operations on tags
CREATE INDEX IF NOT EXISTS idx_user_tags ON user_tag_preferences USING GIN(preferred_tags);

-- Additional indexes for better query performance on existing tables
CREATE INDEX IF NOT EXISTS idx_post_vote_voted_at ON post_votes(user_id, voted_at DESC);
CREATE INDEX IF NOT EXISTS idx_subscription_time ON subscriptions(user_id, subscribed_at DESC);
CREATE INDEX IF NOT EXISTS idx_user_saved_post_followed ON user_saved_posts(user_id, is_followed);

-- =================================================================
-- SAMPLE DATA (Optional - For Testing)
-- =================================================================

-- Sample interest scores (run after users interact with the system)
-- INSERT INTO user_interest_scores (user_id, community_id, score, updated_at)
-- VALUES 
--     (1, 1, 15.0, NOW()),  -- User 1 likes community 1
--     (1, 2, 8.0, NOW()),   -- User 1 also likes community 2
--     (2, 1, 20.0, NOW());  -- User 2 really likes community 1

-- Sample tag preferences
-- INSERT INTO user_tag_preferences (user_id, preferred_tags, updated_at)
-- VALUES
--     (1, ARRAY['technology', 'programming', 'golang'], NOW()),
--     (2, ARRAY['gaming', 'esports', 'entertainment'], NOW());

-- =================================================================
-- ROLLBACK SCRIPT (if needed)
-- =================================================================

-- DROP INDEX IF EXISTS idx_user_saved_post_followed;
-- DROP INDEX IF EXISTS idx_subscription_time;
-- DROP INDEX IF EXISTS idx_post_vote_voted_at;
-- DROP INDEX IF EXISTS idx_user_tags;
-- DROP INDEX IF EXISTS idx_community_score;
-- DROP INDEX IF EXISTS idx_user_interest_score;
-- DROP TABLE IF EXISTS user_tag_preferences;
-- DROP TABLE IF EXISTS user_interest_scores;
