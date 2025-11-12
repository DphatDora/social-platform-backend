package service

import (
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

	"github.com/redis/go-redis/v9"
)

type AuthService struct {
	userRepo                repository.UserRepository
	verificationRepo        repository.UserVerificationRepository
	passwordResetRepo       repository.PasswordResetRepository
	botTaskRepo             repository.BotTaskRepository
	communityModeratorRepo  repository.CommunityModeratorRepository
	notificationSettingRepo repository.NotificationSettingRepository
	botTaskService          *BotTaskService
	redisClient             *redis.Client
}

func NewAuthService(
	userRepo repository.UserRepository,
	verificationRepo repository.UserVerificationRepository,
	passwordResetRepo repository.PasswordResetRepository,
	botTaskRepo repository.BotTaskRepository,
	communityModeratorRepo repository.CommunityModeratorRepository,
	notificationSettingRepo repository.NotificationSettingRepository,
	botTaskService *BotTaskService,
	redisClient *redis.Client,
) *AuthService {
	return &AuthService{
		userRepo:                userRepo,
		verificationRepo:        verificationRepo,
		passwordResetRepo:       passwordResetRepo,
		botTaskRepo:             botTaskRepo,
		communityModeratorRepo:  communityModeratorRepo,
		notificationSettingRepo: notificationSettingRepo,
		botTaskService:          botTaskService,
		redisClient:             redisClient,
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
		Username:     req.Username,
		Email:        req.Email,
		Password:     &hashedPassword,
		IsActive:     false,
		Role:         constant.ROLE_USER,
		Karma:        0,
		AuthProvider: "email",
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

	go func(userID uint64, userEmail string, token string) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[Panic] Recovered in AuthService.Register background tasks: %v", r)
			}
		}()

		actions := []string{
			constant.NOTIFICATION_ACTION_GET_POST_VOTE,
			constant.NOTIFICATION_ACTION_GET_POST_NEW_COMMENT,
			constant.NOTIFICATION_ACTION_GET_COMMENT_VOTE,
			constant.NOTIFICATION_ACTION_GET_COMMENT_REPLY,
			constant.NOTIFICATION_ACTION_POST_APPROVED,
			constant.NOTIFICATION_ACTION_POST_REJECTED,
			constant.NOTIFICATION_ACTION_POST_DELETED,
			constant.NOTIFICATION_ACTION_POST_REPORTED,
			constant.NOTIFICATION_ACTION_SUBSCRIPTION_APPROVED,
			constant.NOTIFICATION_ACTION_SUBSCRIPTION_REJECTED,
		}

		now := time.Now()
		settings := make([]*model.NotificationSetting, len(actions))
		for i, action := range actions {
			settings[i] = &model.NotificationSetting{
				UserID:     userID,
				Action:     action,
				IsPush:     true,
				IsSendMail: false,
				CreatedAt:  now,
				UpdatedAt:  now,
			}
		}

		if err := s.notificationSettingRepo.CreateNotificationSettings(settings); err != nil {
			log.Printf("[Err] Error creating notification settings in AuthService.Register: %v", err)
		}

		// Send verification email
		verificationLink := fmt.Sprintf("%s/api/v1/auth/verify?token=%s", conf.Server.Url, token)
		body, err := util.RenderTemplate("package/template/email/email_verification.html", map[string]interface{}{
			"VerificationLink": template.URL(verificationLink),
			"ExpireMinutes":    conf.Auth.VerifyTokenExpirationMinutes,
		})

		if err != nil {
			log.Printf("[Err] Error rendering email template in AuthService.Register: %v", err)
			return
		}

		if s.botTaskService != nil {
			if err := s.botTaskService.CreateEmailTask(userEmail, "Verify Your Account", body); err != nil {
				log.Printf("[Err] Error creating email task in AuthService.Register: %v", err)
			}
		}
	}(user.ID, user.Email, token)

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

	// Check if user has password (email/password login)
	if user.Password == nil {
		log.Printf("[Err] User %s registered with Google, no password set", req.Email)
		return nil, fmt.Errorf("this account is registered with Google. Please use Google Sign-In")
	}

	if err := util.ComparePassword(*user.Password, req.Password); err != nil {
		log.Printf("[Err] Invalid password for user %s in AuthService.Login", req.Email)
		return nil, fmt.Errorf("invalid email or password")
	}

	// Generate JWT token
	conf := config.GetConfig()
	accessToken, err := util.GenerateJWT(
		user.ID,
		conf.Auth.AccessTokenExpirationMinutes,
		conf.Auth.JWTSecret,
		user.PasswordChangedAt,
	)
	if err != nil {
		log.Printf("[Err] Error generating JWT token in AuthService.Login: %v", err)
		return nil, fmt.Errorf("failed to generate access token")
	}

	loginResponse := &response.LoginResponse{
		AccessToken: accessToken,
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

	go func(userEmail string, token string) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[Panic] Recovered in AuthService.ForgotPassword bot task: %v", r)
			}
		}()

		resetLink := fmt.Sprintf("%s/api/v1/auth/verify-reset?token=%s", conf.Server.Url, token)
		body, err := util.RenderTemplate("package/template/email/password_reset.html", map[string]interface{}{
			"ResetLink":     template.URL(resetLink),
			"ExpireMinutes": conf.Auth.ResetTokenExpirationMinutes,
		})

		if err != nil {
			log.Printf("[Err] Error rendering email template in AuthService.ForgotPassword: %v", err)
			return
		}

		if err := s.botTaskService.CreateEmailTask(userEmail, "Reset Your Password", body); err != nil {
			log.Printf("[Err] Error creating email task in AuthService.ForgotPassword: %v", err)
		}
	}(user.Email, token)

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

	// Invalidate password cache after successful password reset
	if s.redisClient != nil {
		if err := util.InvalidatePasswordCache(s.redisClient, passwordReset.UserID); err != nil {
			log.Printf("[Warn] Error invalidating password cache for user %d: %v", passwordReset.UserID, err)
		}
	}

	if err := s.passwordResetRepo.DeletePasswordReset(passwordReset.ID); err != nil {
		log.Printf("[Err] Error deleting password reset in AuthService.ResetPassword: %v", err)
	}

	return nil
}

func (s *AuthService) ResendVerificationEmail(req *request.ResendVerificationRequest) error {
	// Check if email exists
	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		log.Printf("[Err] Error getting user by email in AuthService.ResendVerificationEmail: %v", err)
		return fmt.Errorf("email not found")
	}

	// Check if user is already active
	if user.IsActive {
		log.Printf("[Err] User %s is already active in AuthService.ResendVerificationEmail", req.Email)
		return fmt.Errorf("email is already verified")
	}

	// Delete existed verification tokens
	if err := s.verificationRepo.DeleteVerificationByUserID(user.ID); err != nil {
		log.Printf("[Err] Error deleting existing verification in AuthService.ResendVerificationEmail: %v", err)
	}

	// Generate new token
	token, err := util.GenerateToken(32)
	if err != nil {
		log.Printf("[Err] Error generating token in AuthService.ResendVerificationEmail: %v", err)
		return fmt.Errorf("failed to generate token")
	}

	conf := config.GetConfig()
	verification := &model.UserVerification{
		UserID:    user.ID,
		Token:     token,
		ExpiredAt: time.Now().UTC().Add(time.Duration(conf.Auth.VerifyTokenExpirationMinutes) * time.Minute),
		CreatedAt: time.Now().UTC(),
	}

	if err := s.verificationRepo.CreateVerification(verification); err != nil {
		log.Printf("[Err] Error creating verification in AuthService.ResendVerificationEmail: %v", err)
		return fmt.Errorf("failed to create verification")
	}

	// Create bot task for sending email in background
	go func(userEmail string, token string) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[Panic] Recovered in AuthService.ResendVerificationEmail bot task: %v", r)
			}
		}()

		verificationLink := fmt.Sprintf("%s/api/v1/auth/verify?token=%s", conf.Server.Url, token)
		body, err := util.RenderTemplate("package/template/email/email_verification.html", map[string]interface{}{
			"VerificationLink": template.URL(verificationLink),
			"ExpireMinutes":    conf.Auth.VerifyTokenExpirationMinutes,
		})

		if err != nil {
			log.Printf("[Err] Error rendering email template in AuthService.ResendVerificationEmail: %v", err)
			return
		}

		if err := s.botTaskService.CreateEmailTask(userEmail, "Verify Your Account", body); err != nil {
			log.Printf("[Err] Error creating email task in AuthService.ResendVerificationEmail: %v", err)
		}
	}(user.Email, token)

	return nil
}

func (s *AuthService) ResendResetPasswordEmail(req *request.ResendVerificationRequest) error {
	// Check if email exists
	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		log.Printf("[Err] Error getting user by email in AuthService.ResendResetPasswordEmail: %v", err)
		return fmt.Errorf("email not found")
	}

	// Check if user is active
	if !user.IsActive {
		log.Printf("[Err] User %s is not active in AuthService.ResendResetPasswordEmail", req.Email)
		return fmt.Errorf("email not verified. Please verify your email first")
	}

	// Delete any existing password reset tokens for this user
	if err := s.passwordResetRepo.DeletePasswordResetByUserID(user.ID); err != nil {
		log.Printf("[Err] Error deleting existing password reset in AuthService.ResendResetPasswordEmail: %v", err)
	}

	// Generate reset token
	token, err := util.GenerateToken(32)
	if err != nil {
		log.Printf("[Err] Error generating token in AuthService.ResendResetPasswordEmail: %v", err)
		return fmt.Errorf("failed to generate reset token")
	}

	conf := config.GetConfig()
	passwordReset := &model.PasswordReset{
		UserID:    user.ID,
		Token:     token,
		ExpiredAt: time.Now().UTC().Add(time.Duration(conf.Auth.ResetTokenExpirationMinutes) * time.Minute),
		CreatedAt: time.Now().UTC(),
	}

	if err := s.passwordResetRepo.CreatePasswordReset(passwordReset); err != nil {
		log.Printf("[Err] Error creating password reset in AuthService.ResendResetPasswordEmail: %v", err)
		return fmt.Errorf("failed to create password reset")
	}

	// Create bot task for sending email in background
	go func(userEmail string, token string) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[Panic] Recovered in AuthService.ResendResetPasswordEmail bot task: %v", r)
			}
		}()

		resetLink := fmt.Sprintf("%s/api/v1/auth/verify-reset?token=%s", conf.Server.Url, token)
		body, err := util.RenderTemplate("package/template/email/password_reset.html", map[string]interface{}{
			"ResetLink":     template.URL(resetLink),
			"ExpireMinutes": conf.Auth.ResetTokenExpirationMinutes,
		})

		if err != nil {
			log.Printf("[Err] Error rendering email template in AuthService.ResendResetPasswordEmail: %v", err)
			return
		}

		if err := s.botTaskService.CreateEmailTask(userEmail, "Reset Your Password", body); err != nil {
			log.Printf("[Err] Error creating email task in AuthService.ResendResetPasswordEmail: %v", err)
		}
	}(user.Email, token)

	return nil
}

func (s *AuthService) GoogleLogin(req *request.GoogleLoginRequest) (*response.LoginResponse, error) {
	conf := config.GetConfig()

	// Verify Google ID Token
	googleUserInfo, err := util.VerifyGoogleIDToken(req.IDToken, conf.Auth.GoogleClientID)
	if err != nil {
		log.Printf("[Err] Error verifying Google ID token in AuthService.GoogleLogin: %v", err)
		return nil, fmt.Errorf("invalid Google ID token: %w", err)
	}

	// Check if user exists by Google ID
	existingUser, err := s.userRepo.GetUserByGoogleID(googleUserInfo.GoogleID)
	if err == nil && existingUser != nil {
		// User already logged in with Google before
		log.Printf("[Info] User %s already exists with Google ID", googleUserInfo.Email)

		// Generate JWT token
		accessToken, err := util.GenerateJWT(
			existingUser.ID,
			conf.Auth.AccessTokenExpirationMinutes,
			conf.Auth.JWTSecret,
			existingUser.PasswordChangedAt,
		)
		if err != nil {
			log.Printf("[Err] Error generating JWT token in AuthService.GoogleLogin: %v", err)
			return nil, fmt.Errorf("failed to generate access token")
		}

		return &response.LoginResponse{
			AccessToken: accessToken,
		}, nil
	}

	// Check if user exists by email
	existingUser, err = s.userRepo.GetUserByEmail(googleUserInfo.Email)
	if err == nil && existingUser != nil {
		log.Printf("[Info] Linking Google account for user %s", googleUserInfo.Email)

		// set new auth provider
		newProvider := constant.ACCOUNT_TYPE_BOTH
		if existingUser.AuthProvider == constant.ACCOUNT_TYPE_GOOGLE {
			newProvider = constant.ACCOUNT_TYPE_GOOGLE
		}

		// Link Google account
		if err := s.userRepo.LinkGoogleAccount(existingUser.ID, googleUserInfo.GoogleID, newProvider); err != nil {
			log.Printf("[Err] Error linking Google account in AuthService.GoogleLogin: %v", err)
			return nil, fmt.Errorf("failed to link Google account: %w", err)
		}

		// Activate user if not already active
		if !existingUser.IsActive {
			if err := s.userRepo.ActivateUser(existingUser.ID); err != nil {
				log.Printf("[Err] Error activating user in AuthService.GoogleLogin: %v", err)
			}
		}

		accessToken, err := util.GenerateJWT(
			existingUser.ID,
			conf.Auth.AccessTokenExpirationMinutes,
			conf.Auth.JWTSecret,
			existingUser.PasswordChangedAt,
		)
		if err != nil {
			log.Printf("[Err] Error generating JWT token in AuthService.GoogleLogin: %v", err)
			return nil, fmt.Errorf("failed to generate access token")
		}

		return &response.LoginResponse{
			AccessToken: accessToken,
		}, nil
	}

	username := googleUserInfo.Name
	if username == "" {
		username = googleUserInfo.Email
	}

	var avatar *string
	if googleUserInfo.Picture != "" {
		avatar = &googleUserInfo.Picture
	}

	newUser := &model.User{
		Username:     username,
		Email:        googleUserInfo.Email,
		Password:     nil,
		GoogleID:     &googleUserInfo.GoogleID,
		AuthProvider: constant.ACCOUNT_TYPE_GOOGLE,
		IsActive:     true,
		Role:         constant.ROLE_USER,
		Karma:        0,
		Avatar:       avatar,
	}

	if err := s.userRepo.CreateUser(newUser); err != nil {
		log.Printf("[Err] Error creating user in AuthService.GoogleLogin: %v", err)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Create default notification settings for user
	go func(userID uint64) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[Err] Panic in goroutine AuthService.GoogleLogin: %v", r)
			}
		}()

		actions := []string{
			constant.NOTIFICATION_ACTION_GET_POST_VOTE,
			constant.NOTIFICATION_ACTION_GET_POST_NEW_COMMENT,
			constant.NOTIFICATION_ACTION_GET_COMMENT_VOTE,
			constant.NOTIFICATION_ACTION_GET_COMMENT_REPLY,
			constant.NOTIFICATION_ACTION_POST_APPROVED,
			constant.NOTIFICATION_ACTION_POST_REJECTED,
			constant.NOTIFICATION_ACTION_POST_DELETED,
			constant.NOTIFICATION_ACTION_POST_REPORTED,
			constant.NOTIFICATION_ACTION_SUBSCRIPTION_APPROVED,
			constant.NOTIFICATION_ACTION_SUBSCRIPTION_REJECTED,
		}

		now := time.Now()
		settings := make([]*model.NotificationSetting, len(actions))
		for i, action := range actions {
			settings[i] = &model.NotificationSetting{
				UserID:     userID,
				Action:     action,
				IsPush:     true,
				IsSendMail: false,
				CreatedAt:  now,
				UpdatedAt:  now,
			}
		}

		if err := s.notificationSettingRepo.CreateNotificationSettings(settings); err != nil {
			log.Printf("[Err] Error creating notification settings in AuthService.GoogleLogin: %v", err)
		}
	}(newUser.ID)

	accessToken, err := util.GenerateJWT(
		newUser.ID,
		conf.Auth.AccessTokenExpirationMinutes,
		conf.Auth.JWTSecret,
		newUser.PasswordChangedAt,
	)
	if err != nil {
		log.Printf("[Err] Error generating JWT token in AuthService.GoogleLogin: %v", err)
		return nil, fmt.Errorf("failed to generate access token")
	}

	return &response.LoginResponse{
		AccessToken: accessToken,
	}, nil
}
