-- DATABASE DESIGN FOR SOCIAL COMMUNICATION PLATFORM
-- v 1.0 (09/09/2025) --

-- Người dùng
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  username VARCHAR(50) NOT NULL UNIQUE,
  email VARCHAR(100) NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  karma INT NOT NULL DEFAULT 0,
  bio TEXT,
  avatar_url TEXT,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  role VARCHAR(20) NOT NULL DEFAULT 'user'
);
CREATE INDEX idx_users_created_at ON users(created_at);

-- Cộng đồng
CREATE TABLE communities (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL UNIQUE,
  description TEXT,
  created_by INT NOT NULL REFERENCES users(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  is_private BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE INDEX idx_communities_created_by ON communities(created_by);
CREATE INDEX idx_communities_created_at ON communities(created_at);

-- Quản trị cộng đồng
CREATE TABLE community_moderators (
  community_id INT NOT NULL REFERENCES communities(id),
  user_id INT NOT NULL REFERENCES users(id),
  role VARCHAR(20) NOT NULL DEFAULT 'moderator',
  joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (community_id, user_id)
);
CREATE INDEX idx_com_mod_user ON community_moderators(user_id);

-- Theo dõi cộng đồng
CREATE TABLE subscriptions (
  user_id INT NOT NULL REFERENCES users(id),
  community_id INT NOT NULL REFERENCES communities(id),
  subscribed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (user_id, community_id)
);
CREATE INDEX idx_subscriptions_community ON subscriptions(community_id);

-- Bài viết
CREATE TABLE posts (
  id SERIAL PRIMARY KEY,
  community_id INT NOT NULL REFERENCES communities(id),
  author_id INT NOT NULL REFERENCES users(id),
  title VARCHAR(255) NOT NULL,
  content TEXT,
  url TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_posts_community ON posts(community_id);
CREATE INDEX idx_posts_author ON posts(author_id);
CREATE INDEX idx_posts_created_at ON posts(created_at);

-- Bình luận
CREATE TABLE comments (
  id SERIAL PRIMARY KEY,
  post_id INT NOT NULL REFERENCES posts(id),
  author_id INT NOT NULL REFERENCES users(id),
  parent_comment_id INT REFERENCES comments(id),
  content TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_comments_post ON comments(post_id);
CREATE INDEX idx_comments_author ON comments(author_id);
CREATE INDEX idx_comments_parent ON comments(parent_comment_id);

-- Bình chọn bài viết
CREATE TABLE post_votes (
  user_id INT NOT NULL REFERENCES users(id),
  post_id INT NOT NULL REFERENCES posts(id),
  vote SMALLINT NOT NULL,
  voted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (user_id, post_id)
);
CREATE INDEX idx_post_votes_post ON post_votes(post_id);

-- Bình chọn bình luận
CREATE TABLE comment_votes (
  user_id INT NOT NULL REFERENCES users(id),
  comment_id INT NOT NULL REFERENCES comments(id),
  vote SMALLINT NOT NULL,
  voted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (user_id, comment_id)
);
CREATE INDEX idx_comment_votes_comment ON comment_votes(comment_id);

-- Huy hiệu
CREATE TABLE badges (
  id SERIAL PRIMARY KEY,
  name VARCHAR(50) NOT NULL UNIQUE,
  description TEXT,
  icon_url TEXT
);

-- Huy hiệu của người dùng
CREATE TABLE user_badges (
  user_id INT NOT NULL REFERENCES users(id),
  badge_id INT NOT NULL REFERENCES badges(id),
  awarded_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (user_id, badge_id)
);
CREATE INDEX idx_user_badges_badge ON user_badges(badge_id);

-- Tin nhắn
CREATE TABLE messages (
  id SERIAL PRIMARY KEY,
  sender_id INT NOT NULL REFERENCES users(id),
  receiver_id INT NOT NULL REFERENCES users(id),
  subject VARCHAR(255),
  content TEXT NOT NULL,
  sent_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  is_read BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE INDEX idx_messages_receiver ON messages(receiver_id);
CREATE INDEX idx_messages_sender ON messages(sender_id);

-- Thông báo
CREATE TABLE notifications (
  id SERIAL PRIMARY KEY,
  user_id INT NOT NULL REFERENCES users(id),
  type VARCHAR(50),
  reference_id INT,
  is_read BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_notifications_user ON notifications(user_id);
