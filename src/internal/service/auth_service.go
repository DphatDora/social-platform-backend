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
	"social-platform-backend/internal/interface/dto/response"
	"social-platform-backend/package/constant"
	"social-platform-backend/package/util"
	"time"
)

type AuthService struct {
	userRepo          repository.UserRepository
	verificationRepo  repository.UserVerificationRepository
	passwordResetRepo repository.PasswordResetRepository
	botTaskRepo       repository.BotTaskRepository
}

func NewAuthService(
	userRepo repository.UserRepository,
	verificationRepo repository.UserVerificationRepository,
	passwordResetRepo repository.PasswordResetRepository,
	botTaskRepo repository.BotTaskRepository,
) *AuthService {
	return &AuthService{
		userRepo:          userRepo,
		verificationRepo:  verificationRepo,
		passwordResetRepo: passwordResetRepo,
		botTaskRepo:       botTaskRepo,
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
	verificationLink := fmt.Sprintf("%s/api/v1/auth/verify?token=%s", conf.Server.Url, token)
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

func (s *AuthService) Login(req *request.LoginRequest) (*response.LoginResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		log.Printf("[Err] Error getting user by email in AuthService.Login: %v", err)
		return nil, fmt.Errorf("invalid email or password")
	}

	// Check if user is active
	if !user.IsActive {
		log.Printf("[Err] User %s is not active in AuthService.Login", req.Email)
		return nil, fmt.Errorf("email not verified. Please verify your email first")
	}

	if err := util.ComparePassword(user.Password, req.Password); err != nil {
		log.Printf("[Err] Invalid password for user %s in AuthService.Login", req.Email)
		return nil, fmt.Errorf("invalid email or password")
	}

	// Generate JWT token
	conf := config.GetConfig()
	accessToken, err := util.GenerateJWT(
		user.ID,
		user.Role,
		conf.Auth.AccessTokenExpirationMinutes,
		conf.Auth.JWTSecret,
	)
	if err != nil {
		log.Printf("[Err] Error generating JWT token in AuthService.Login: %v", err)
		return nil, fmt.Errorf("failed to generate access token")
	}

	loginResponse := &response.LoginResponse{
		Username:    user.Username,
		AccessToken: accessToken,
	}

	if user.Avatar != nil {
		loginResponse.Avatar = *user.Avatar
	}

	return loginResponse, nil
}

func (s *AuthService) ForgotPassword(req *request.ForgotPasswordRequest) error {
	// Check if email exists
	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		log.Printf("[Err] Error getting user by email in AuthService.ForgotPassword: %v", err)
		return fmt.Errorf("if your email is registered, you will receive a password reset link")
	}

	// Check if user is active
	if !user.IsActive {
		log.Printf("[Err] User %s is not active in AuthService.ForgotPassword", req.Email)
		return fmt.Errorf("email not verified. Please verify your email first")
	}

	if err := s.passwordResetRepo.DeletePasswordResetByUserID(user.ID); err != nil {
		log.Printf("[Err] Error deleting existing password reset in AuthService.ForgotPassword: %v", err)
	}

	// Generate reset token
	token, err := util.GenerateToken(32)
	if err != nil {
		log.Printf("[Err] Error generating token in AuthService.ForgotPassword: %v", err)
		return fmt.Errorf("failed to generate reset token")
	}

	conf := config.GetConfig()
	passwordReset := &model.PasswordReset{
		UserID:    user.ID,
		Token:     token,
		ExpiredAt: time.Now().Add(time.Duration(conf.Auth.ResetTokenExpirationMinutes) * time.Minute),
		CreatedAt: time.Now(),
	}

	if err := s.passwordResetRepo.CreatePasswordReset(passwordReset); err != nil {
		log.Printf("[Err] Error creating password reset in AuthService.ForgotPassword: %v", err)
		return fmt.Errorf("failed to create password reset")
	}

	// Create bot task for sending email
	resetLink := fmt.Sprintf("%s/api/v1/auth/verify-reset?token=%s", conf.Server.Url, token)
	body, err := util.RenderTemplate("package/template/email/password_reset.html", map[string]interface{}{
		"ResetLink":     template.URL(resetLink),
		"ExpireMinutes": conf.Auth.ResetTokenExpirationMinutes,
	})

	if err != nil {
		log.Printf("[Err] Error rendering email template in AuthService.ForgotPassword: %v", err)
		return fmt.Errorf("failed to render email template")
	}

	emailPayload := request.EmailPayload{
		To:      user.Email,
		Subject: "Reset Your Password",
		Body:    body,
	}

	payloadBytes, err := json.Marshal(emailPayload)
	if err != nil {
		log.Printf("[Err] Error marshaling email payload in AuthService.ForgotPassword: %v", err)
		return fmt.Errorf("failed to marshal email payload")
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
		log.Printf("[Err] Error creating bot task in AuthService.ForgotPassword: %v", err)
		return fmt.Errorf("failed to create bot task")
	}

	return nil
}

func (s *AuthService) VerifyResetToken(token string) (string, error) {
	passwordReset, err := s.passwordResetRepo.GetPasswordResetByToken(token)
	if err != nil {
		log.Printf("[Err] Error getting password reset by token in AuthService.VerifyResetToken: %v", err)
		return "", fmt.Errorf("invalid or expired token")
	}

	// Check if token is expired
	if time.Now().After(passwordReset.ExpiredAt) {
		log.Printf("[Err] Token expired in AuthService.VerifyResetToken for user %d", passwordReset.UserID)
		return "", fmt.Errorf("token has expired")
	}

	return token, nil
}

func (s *AuthService) ResetPassword(req *request.ResetPasswordRequest) error {
	// Verify token
	passwordReset, err := s.passwordResetRepo.GetPasswordResetByToken(req.Token)
	if err != nil {
		log.Printf("[Err] Error getting password reset by token in AuthService.ResetPassword: %v", err)
		return fmt.Errorf("invalid or expired token")
	}

	// Check if token is expired
	if time.Now().After(passwordReset.ExpiredAt) {
		log.Printf("[Err] Token expired in AuthService.ResetPassword for user %d", passwordReset.UserID)
		return fmt.Errorf("token has expired")
	}

	// Hash new password
	hashedPassword, err := util.HashPassword(req.NewPassword)
	if err != nil {
		log.Printf("[Err] Error hashing password in AuthService.ResetPassword: %v", err)
		return fmt.Errorf("failed to hash password")
	}

	if err := s.userRepo.UpdatePasswordAndSetChangedAt(passwordReset.UserID, hashedPassword); err != nil {
		log.Printf("[Err] Error updating password in AuthService.ResetPassword: %v", err)
		return fmt.Errorf("failed to update password")
	}

	if err := s.passwordResetRepo.DeletePasswordReset(passwordReset.ID); err != nil {
		log.Printf("[Err] Error deleting password reset in AuthService.ResetPassword: %v", err)
	}

	return nil
}
