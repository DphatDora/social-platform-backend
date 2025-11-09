package service

import (
	"encoding/json"
	"fmt"
	"log"
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"
	"social-platform-backend/package/constant"
	"social-platform-backend/package/template/payload"
	"time"
)

type BotTaskService struct {
	botTaskRepo repository.BotTaskRepository
}

func NewBotTaskService(botTaskRepo repository.BotTaskRepository) *BotTaskService {
	return &BotTaskService{
		botTaskRepo: botTaskRepo,
	}
}

func (s *BotTaskService) CreateInterestScoreTask(userID, communityID uint64, action string, postID *uint64) error {
	scorePayload := payload.UpdateInterestScorePayload{
		UserID:      userID,
		CommunityID: communityID,
		Action:      action,
		PostID:      postID,
		UpdatedAt:   time.Now(),
	}

	payloadBytes, err := json.Marshal(scorePayload)
	if err != nil {
		log.Printf("[Err] Error marshaling interest score payload: %v", err)
		return fmt.Errorf("failed to marshal payload")
	}

	rawPayload := json.RawMessage(payloadBytes)
	now := time.Now()
	botTask := &model.BotTask{
		Action:     constant.BOT_TASK_ACTION_UPDATE_INTEREST_SCORE,
		Payload:    &rawPayload,
		CreatedAt:  now,
		ExecutedAt: &now,
	}

	if err := s.botTaskRepo.CreateBotTask(botTask); err != nil {
		log.Printf("[Err] Error creating bot task for interest score: %v", err)
		return fmt.Errorf("failed to create bot task")
	}

	log.Printf("[Info] Created interest score task for user %d, community %d, action: %s", userID, communityID, action)
	return nil
}

func (s *BotTaskService) CreateKarmaTask(userID uint64, targetID *uint64, action string) error {
	karmaPayload := payload.UpdateUserKarmaPayload{
		UserId:    userID,
		TargetId:  targetID,
		Action:    action,
		UpdatedAt: time.Now(),
	}

	payloadBytes, err := json.Marshal(karmaPayload)
	if err != nil {
		log.Printf("[Err] Error marshaling karma payload: %v", err)
		return fmt.Errorf("failed to marshal karma payload")
	}

	rawPayload := json.RawMessage(payloadBytes)
	now := time.Now()
	botTask := &model.BotTask{
		Action:     constant.BOT_TASK_ACTION_UPDATE_KARMA,
		Payload:    &rawPayload,
		CreatedAt:  now,
		ExecutedAt: &now,
	}

	if err := s.botTaskRepo.CreateBotTask(botTask); err != nil {
		log.Printf("[Err] Error creating karma bot task: %v", err)
		return fmt.Errorf("failed to create karma task")
	}

	log.Printf("[Info] Created karma task for user %d, action: %s", userID, action)
	return nil
}

func (s *BotTaskService) CreateEmailTask(recipientEmail, subject, body string) error {
	emailPayload := payload.EmailPayload{
		To:      recipientEmail,
		Subject: subject,
		Body:    body,
	}

	payloadBytes, err := json.Marshal(emailPayload)
	if err != nil {
		log.Printf("[Err] Error marshaling email payload: %v", err)
		return fmt.Errorf("failed to marshal email payload")
	}

	rawPayload := json.RawMessage(payloadBytes)
	now := time.Now()
	botTask := &model.BotTask{
		Action:     constant.BOT_TASK_ACTION_SEND_EMAIL,
		Payload:    &rawPayload,
		CreatedAt:  now,
		ExecutedAt: &now,
	}

	if err := s.botTaskRepo.CreateBotTask(botTask); err != nil {
		log.Printf("[Err] Error creating email bot task: %v", err)
		return fmt.Errorf("failed to create email task")
	}

	log.Printf("[Info] Created email task for recipient: %s", recipientEmail)
	return nil
}
