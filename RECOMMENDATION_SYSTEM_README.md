# Recommendation System Implementation

## 📋 Tổng quan

Hệ thống recommendation (gợi ý) cho mạng xã hội cộng đồng, dựa trên:

- **Hành vi người dùng**: Upvote/downvote bài viết, follow bài viết, join cộng đồng
- **Scoring system**: Tính điểm sở thích của user với từng community
- **Content matching**: Matching tags và topics

## 🏗️ Kiến trúc

### 1. Models mới

- `UserInterestScore`: Lưu điểm sở thích của user với community
- `UserTagPreference`: Cache các tags yêu thích của user

### 2. Services

- `BotTaskService`: Tạo bot task cho các hành động
- `RecommendationService`: Tính toán và trả về danh sách recommend

### 3. Repositories

- `UserInterestScoreRepository`: CRUD cho interest scores
- `UserTagPreferenceRepository`: CRUD cho tag preferences

## 🔄 Luồng hoạt động

### Main API (Project này)

```
User Action (Vote/Follow/Join)
    ↓
Create BotTask → Save to Database
    ↓
Return response to user
```

### Background Worker (Project riêng)

```
Fetch unprocessed BotTasks
    ↓
Parse payload & calculate score delta
    ↓
Update UserInterestScore
    ↓
Mark task as executed
```

### Recommendation API

```
GET /api/v1/posts?sortBy=best
    ↓
Get user's top communities by score
    ↓
Fetch posts from those communities
    ↓
Score posts by: tags match, votes, freshness, author karma
    ↓
Return sorted posts
```

## 📊 Scoring Weights

| Hành động      | Điểm | Lý do              |
| -------------- | ---- | ------------------ |
| Join Community | +10  | Cam kết cao nhất   |
| Follow Post    | +3   | Quan tâm lâu dài   |
| Upvote Post    | +2   | Tương tác tích cực |
| Downvote Post  | -1   | Tín hiệu tiêu cực  |

## 🚀 Cách sử dụng

### 1. Chạy Migration

```sql
psql -U your_user -d your_database -f migrations_recommendation_system.sql
```

### 2. Update Service Dependencies

Cần update dependency injection trong `main.go` hoặc config:

```go
// Initialize repositories
userInterestScoreRepo := repository.NewUserInterestScoreRepository(db)
userTagPrefRepo := repository.NewUserTagPreferenceRepository(db)

// Initialize services
botTaskService := service.NewBotTaskService(botTaskRepo)
recommendService := service.NewRecommendationService(
    userInterestScoreRepo,
    userTagPrefRepo,
    postRepo,
    communityRepo,
)

// Update existing services with new dependencies
postService := service.NewPostService(
    postRepo,
    communityRepo,
    postVoteRepo,
    postReportRepo,
    botTaskRepo,
    userRepo,
    notificationService,
    botTaskService,      // ADD THIS
    recommendService,    // ADD THIS
)

communityService := service.NewCommunityService(
    communityRepo,
    subscriptionRepo,
    communityModeratorRepo,
    postRepo,
    postReportRepo,
    notificationService,
    botTaskService,      // ADD THIS
)

userService := service.NewUserService(
    userRepo,
    communityRepo,
    communityModeratorRepo,
    userSavedPostRepo,
    postRepo,            // ADD THIS
    botTaskService,      // ADD THIS
)
```

### 3. Setup Background Worker

Sử dụng code mẫu trong `SAMPLE_BOTTASK_WORKER.go`:

```go
// In your background worker project
processor := NewInterestScoreProcessor(db)

// Run every 10 seconds
ticker := time.NewTicker(10 * time.Second)
for range ticker.C {
    processor.ProcessBotTasks(100)
}
```

### 4. Test API

#### Get recommended posts (all communities)

```bash
GET /api/v1/posts?sortBy=best&page=1&limit=10
Authorization: Bearer {token}
```

#### Get recommended posts in specific community

```bash
GET /api/v1/communities/1/posts?sortBy=best&page=1&limit=10
Authorization: Bearer {token}
```

## 📝 API Changes

### Existing Endpoints với sortBy=best

- `GET /api/v1/posts?sortBy=best` - Recommended posts for user
- `GET /api/v1/communities/:id/posts?sortBy=best` - Recommended posts in community

### Bot Tasks Created Automatically

Bot tasks được tạo tự động khi:

- User upvote/downvote bài viết
- User follow bài viết (update is_followed = true)
- User join cộng đồng

## 🗄️ Database Schema

### user_interest_scores

```sql
user_id        BIGINT    (PK, FK → users.id)
community_id   BIGINT    (PK, FK → communities.id)
score          DOUBLE    (index)
last_vote_at   TIMESTAMP
last_join_at   TIMESTAMP
updated_at     TIMESTAMP
```

### user_tag_preferences

```sql
user_id         BIGINT    (PK, FK → users.id)
preferred_tags  TEXT[]    (GIN index)
updated_at      TIMESTAMP
```

## 🎯 Recommendation Algorithm

```
POST_SCORE =
    (Community_Interest_Score * 0.5) +
    (Tag_Match_Count * 10) +
    (Post_Vote / 10) +
    (Freshness_Hours / -24) +
    (Author_Karma / 100)
```

## ⚡ Performance Tips

1. **Batch Processing**: Worker xử lý bot tasks theo batch (100 tasks/lần)
2. **Index Usage**: Đảm bảo các indexes được tạo đúng
3. **Cache**: Có thể cache top communities của user trong Redis
4. **Limit Queries**: Chỉ lấy posts trong 30 ngày gần đây
5. **Tag Preferences**: Update định kỳ (hàng ngày/tuần), không real-time

## 🔧 Troubleshooting

### Bot tasks không được xử lý?

- Kiểm tra background worker có chạy không
- Check logs của worker service
- Xem bảng `bot_tasks` có records với `executed_at = NULL`

### Recommend trả về rỗng?

- User chưa có hành động nào → Chưa có interest score
- Kiểm tra bảng `user_interest_scores` có data không
- Fallback về sortBy mặc định (hot/new)

### Performance chậm?

- Kiểm tra indexes đã được tạo chưa
- Giảm số lượng communities trong recommendation
- Cache kết quả recommend trong Redis (TTL 5-10 phút)

## 📚 Files Changed

### Models

- `internal/domain/model/user_interest_score.go` (NEW)
- `internal/domain/model/user_tag_preference.go` (NEW)

### Repositories

- `internal/domain/repository/user_interest_score_repository.go` (NEW)
- `internal/domain/repository/user_tag_preference_repository.go` (NEW)
- `internal/infrastructure/db/repository/user_interest_score_repository_impl.go` (NEW)
- `internal/infrastructure/db/repository/user_tag_preference_repository_impl.go` (NEW)

### Services

- `internal/service/bottask_service.go` (NEW)
- `internal/service/recommendation_service.go` (NEW)
- `internal/service/post_service.go` (MODIFIED - added bot task creation + recommend)
- `internal/service/community_service.go` (MODIFIED - added bot task creation)
- `internal/service/user_service.go` (MODIFIED - added bot task creation)

### Constants & Payloads

- `package/constant/bottask_action.go` (MODIFIED)
- `package/constant/interest_action.go` (NEW)
- `package/template/payload/interest_score_payload.go` (NEW)

### Documentation

- `migrations_recommendation_system.sql` (NEW)
- `SAMPLE_BOTTASK_WORKER.go` (NEW - for background worker)

## ✅ TODO Next Steps

1. ✅ Run migration SQL
2. ⬜ Update dependency injection in main.go
3. ⬜ Setup background worker service
4. ⬜ Test bot task creation
5. ⬜ Test recommendation API
6. ⬜ Monitor performance
7. ⬜ Add Redis caching (optional)
8. ⬜ Setup periodic tag preference updates

## 📞 Support

Nếu có vấn đề, kiểm tra:

1. Database migrations đã chạy thành công
2. Background worker service đang chạy
3. Bot tasks được tạo ra (check bảng `bot_tasks`)
4. Logs của cả API và worker service
