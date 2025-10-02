package service

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"social-platform-backend/config"
	"social-platform-backend/internal/domain/model"
	"social-platform-backend/internal/domain/repository"
	"social-platform-backend/internal/interface/dto/request"
	"social-platform-backend/package/constant"
	"social-platform-backend/package/util"
	"time"
)

type AuthService struct {
	userRepo         repository.UserRepository
	verificationRepo repository.UserVerificationRepository
	botTaskRepo      repository.BotTaskRepository
}

func NewAuthService(
	userRepo repository.UserRepository,
	verificationRepo repository.UserVerificationRepository,
	botTaskRepo repository.BotTaskRepository,
) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		verificationRepo: verificationRepo,
		botTaskRepo:      botTaskRepo,
	}
}

func (s *AuthService) Register(req *request.RegisterRequest) error {
	// Check if email exists
	exists, err := s.userRepo.IsEmailExisted(req.Email)
	if err != nil {
		log.Printf("[Err] Error checking email existence in AuthService.Register: %v", err)
		return fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return fmt.Errorf("email already exists")
	}

	// Hash password
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		log.Printf("[Err] Error hashing password in AuthService.Register: %v", err)
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		IsActive: false,
		Role:     constant.ROLE_USER,
		Karma:    0,
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		log.Printf("[Err] Error creating user in AuthService.Register: %v", err)
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Generate verification token
	token, err := util.GenerateToken(32)
	if err != nil {
		log.Printf("[Err] Error generating token in AuthService.Register: %v", err)
		return fmt.Errorf("failed to generate token: %w", err)
	}

	conf := config.GetConfig()
	verification := &model.UserVerification{
		UserID:    user.ID,
		Token:     token,
		ExpiredAt: time.Now().Add(time.Duration(conf.Auth.VerifyTokenExpirationMinutes) * time.Minute),
		CreatedAt: time.Now(),
	}

	if err := s.verificationRepo.CreateVerification(verification); err != nil {
		log.Printf("[Err] Error creating verification in AuthService.Register: %v", err)
		return fmt.Errorf("failed to create verification: %w", err)
	}

	// Create bot task for sending email
	verificationLink := fmt.Sprintf("http://%s:%d/api/v1/auth/verify?token=%s", conf.App.Host, conf.App.Port, token)
	body, err := util.RenderTemplate("package/template/email/email_verification.html", map[string]interface{}{
		"VerificationLink": template.URL(verificationLink),
		"ExpireMinutes":    conf.Auth.VerifyTokenExpirationMinutes,
	})

	if err != nil {
		return fmt.Errorf("failed to render email template: %w", err)
	}

	emailPayload := request.EmailPayload{
		To:      user.Email,
		Subject: "Verify Your Account",
		Body:    body,
	}

	payloadBytes, err := json.Marshal(emailPayload)
	if err != nil {
		log.Printf("[Err] Error marshaling email payload in AuthService.Register: %v", err)
		return fmt.Errorf("failed to marshal email payload: %w", err)
	}

	rawPayload := json.RawMessage(payloadBytes)
	now := time.Now()
	botTask := &model.BotTask{
		Action:     "send_email",
		Payload:    &rawPayload,
		CreatedAt:  now,
		ExecutedAt: &now,
	}

	if err := s.botTaskRepo.CreateBotTask(botTask); err != nil {
		log.Printf("[Err] Error creating bot task in AuthService.Register: %v", err)
		return fmt.Errorf("failed to create bot task: %w", err)
	}

	return nil
}

func (s *AuthService) VerifyEmail(token string) error {
	verification, err := s.verificationRepo.GetVerificationByToken(token)
	if err != nil {
		log.Printf("[Err] Error getting verification by token in AuthService.VerifyEmail: %v", err)
		return fmt.Errorf("invalid or expired token")
	}

	// Check if token is expired
	if time.Now().After(verification.ExpiredAt) {
		log.Printf("[Err] Token expired in AuthService.VerifyEmail for user %d", verification.UserID)
		return fmt.Errorf("token has expired")
	}

	// Activate user
	if err := s.userRepo.ActivateUser(verification.UserID); err != nil {
		log.Printf("[Err] Error activating user in AuthService.VerifyEmail: %v", err)
		return fmt.Errorf("failed to activate user: %w", err)
	}

	// Delete verification record
	if err := s.verificationRepo.DeleteVerification(verification.ID); err != nil {
		log.Printf("[Err] Error deleting verification in AuthService.VerifyEmail: %v", err)
	}

	return nil
}
